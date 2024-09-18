package repo

import (
	"context"
	"database/sql"
	"ecommerce-api/errortools"
	"ecommerce-api/internal/entities"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

func (ecommerce *EcommerceRepo) InsertProduct(ctx context.Context, product entities.Product) (string, error) {
	tx, err := ecommerce.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var productID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO products (name, description, image_url, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		product.Name, product.Description, product.ImageURL, true).Scan(&productID)

	if err != nil {
		return "", fmt.Errorf("could not insert product: %w", err)
	}

	if len(product.ChildVariants) > 0 {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO product_variant_mapping (product_id, variant_id)
			SELECT $1, unnest($2::uuid[])`,
			productID, pq.Array(product.ChildVariants))

		if err != nil {
			return "", fmt.Errorf("could not insert product-variant mappings: %w", err)
		}
	}

	return productID, nil
}

func (ecommerce *EcommerceRepo) UpdateProduct(ctx context.Context, product entities.Product) error {
	tx, err := ecommerce.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, `
		UPDATE products 
		SET name = $1, description = $2, image_url = $3, updated_at = NOW()
		WHERE id = $4`,
		product.Name, product.Description, product.ImageURL, product.ID)

	if err != nil {
		return fmt.Errorf("could not update product: %w", err)
	}

	if len(product.ChildVariants) > 0 {
		_, err = tx.ExecContext(ctx, `
		DELETE FROM product_variant_mapping WHERE product_id = $1`, product.ID)

		if err != nil {
			return fmt.Errorf("could not delete old product-variant mappings: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO product_variant_mapping (product_id, variant_id)
			SELECT $1, unnest($2::uuid[])`,
			product.ID, pq.Array(product.ChildVariants))

		if err != nil {
			return fmt.Errorf("could not insert product-variant mappings: %w", err)
		}
	}

	return nil
}

func (ecommerce *EcommerceRepo) DeleteProduct(ctx context.Context, productID string) error {
	_, err := ecommerce.db.ExecContext(ctx, `
		UPDATE products 
		SET is_active = FALSE 
		WHERE id = $1`, productID)

	if err != nil {
		return fmt.Errorf("could not delete product: %w", err)
	}

	return nil
}

func (ecommerce *EcommerceRepo) ValidateProducts(ctx context.Context, productIDs []string) ([]string, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}

	// Prepare the query
	query := `SELECT id FROM products WHERE id = ANY($1) AND is_active = TRUE`
	rows, err := ecommerce.db.Query(query, pq.Array(productIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Collect valid IDs
	validIDs := make([]string, 0, len(productIDs))
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		validIDs = append(validIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(validIDs) != len(productIDs) {
		return nil, fmt.Errorf("one or more product IDs are invalid")
	}

	return validIDs, nil
}

func (ecommerce *EcommerceRepo) HasActiveProducts(ctx context.Context, categoryID string) (bool, error) {
	var count int
	err := ecommerce.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM products p
		JOIN product_category_mapping pcm ON p.id = pcm.product_id
		WHERE pcm.category_id = $1 AND p.is_active = TRUE`,
		categoryID).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("could not check for active products: %w", err)
	}

	return count > 0, nil
}

func (ecommerce *EcommerceRepo) CheckProductIsActive(ctx context.Context, productID string) (bool, error) {
	var isActive bool
	query := `SELECT is_active FROM products WHERE id = $1`
	err := ecommerce.db.QueryRowContext(ctx, query, productID).Scan(&isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("product not found")
		}
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return isActive, nil
}

func (repo *EcommerceRepo) GetProductDetails(ctx context.Context, productID string) (*entities.ProductDetails, error) {
	query := `
    WITH product_data AS (
        SELECT
            p.id AS product_id,
            p.name AS product_name,
            p.description AS product_description,
            p.image_url AS product_image_url,
            p.is_active AS product_is_active,
            p.created_at AS product_created_at,
            p.updated_at AS product_updated_at,
            p.created_by AS product_created_by,
            p.updated_by AS product_updated_by
        FROM products p
        WHERE p.id = $1
    ),
    categories_data AS (
        SELECT
            c.id AS category_id,
            c.name AS category_name
        FROM categories c
        JOIN product_category_mapping pcm ON c.id = pcm.category_id
        WHERE pcm.product_id = $1
    ),
    variants_data AS (
        SELECT
            v.id AS variant_id,
            v.name AS variant_name,
            v.mrp AS variant_mrp,
            v.discount_price AS variant_discount_price,
            v.size AS variant_size,
            v.color AS variant_color,
            v.quantity AS variant_quantity,
            v.is_active AS variant_is_active,
            COALESCE(v.created_at, NOW()) AS variant_created_at,
            COALESCE(v.updated_at, NOW()) AS variant_updated_at
        FROM variants v
        JOIN product_variant_mapping pvm ON v.id = pvm.variant_id
        WHERE pvm.product_id = $1
    )

    SELECT
        json_build_object(
            'id', p.product_id,
            'name', p.product_name,
            'description', p.product_description,
            'image_url', p.product_image_url,
            'is_active', p.product_is_active,
            'created_at', p.product_created_at,
            'updated_at', p.product_updated_at,
            'created_by', p.product_created_by,
            'updated_by', p.product_updated_by
        ) AS product,
        COALESCE((SELECT json_agg(json_build_object(
            'id', c.category_id,
            'name', c.category_name
        )) FROM categories_data c), '[]') AS categories,
        COALESCE((SELECT json_agg(json_build_object(
            'id', v.variant_id,
            'name', v.variant_name,
            'mrp', v.variant_mrp,
            'discount_price', v.variant_discount_price,
            'size', v.variant_size,
            'color', v.variant_color,
            'quantity', v.variant_quantity,
            'is_active', v.variant_is_active,
            'created_at', v.variant_created_at,
            'updated_at', v.variant_updated_at
        )) FROM variants_data v), '[]') AS variants
    FROM product_data p
    `

	row := repo.db.QueryRowContext(ctx, query, productID)

	var productJSON, categoriesJSON, variantsJSON string
	err := row.Scan(&productJSON, &categoriesJSON, &variantsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errortools.New("product_not_found", errortools.WithDetail("Variant with ID "+productID+" not found"))
		}
		return nil, fmt.Errorf("could not fetch product details: %w", err)
	}

	var product entities.ProductDetails
	var categories []*entities.Category
	var variants []*entities.Variant

	// Unmarshal JSON responses
	if err := json.Unmarshal([]byte(productJSON), &product); err != nil {
		return nil, fmt.Errorf("could not unmarshal product details: %w", err)
	}
	if err := json.Unmarshal([]byte(categoriesJSON), &categories); err != nil {
		return nil, fmt.Errorf("could not unmarshal categories: %w", err)
	}
	if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
		return nil, fmt.Errorf("could not unmarshal variants: %w", err)
	}

	product.Categories = categories
	product.Variants = variants

	return &product, nil
}

func (repo *EcommerceRepo) FetchProducts(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.ProductDetails, int, error) {
	var products []*entities.ProductDetails
	var total int

	query := `
    WITH total_count AS (
        SELECT COUNT(*) AS total FROM products WHERE name ILIKE '%' || $1 || '%'
    ),
    product_data AS (
        SELECT
            p.id AS product_id,
            p.name AS product_name,
            p.description AS product_description,
            p.image_url AS product_image_url,
            p.is_active AS product_is_active,
            p.created_at AS product_created_at,
            p.updated_at AS product_updated_at,
            p.created_by AS product_created_by,
            p.updated_by AS product_updated_by
        FROM products p
        WHERE p.name ILIKE '%' || $1 || '%'
        ORDER BY p.created_at DESC
        LIMIT $2 OFFSET $3
    ),
    categories_data AS (
        SELECT
            c.id AS category_id,
            c.name AS category_name,
            pcm.product_id AS product_id
        FROM categories c
        JOIN product_category_mapping pcm ON c.id = pcm.category_id
    ),
    variants_data AS (
        SELECT
            v.id AS variant_id,
            v.name AS variant_name,
            v.mrp AS variant_mrp,
            v.discount_price AS variant_discount_price,
            v.size AS variant_size,
            v.color AS variant_color,
            v.quantity AS variant_quantity,
            v.is_active AS variant_is_active,
            COALESCE(v.created_at, NOW()) AS variant_created_at,
            COALESCE(v.updated_at, NOW()) AS variant_updated_at,
            pvm.product_id AS product_id
        FROM variants v
        JOIN product_variant_mapping pvm ON v.id = pvm.variant_id
    )

    SELECT
        json_build_object(
            'id', p.product_id,
            'name', p.product_name,
            'description', p.product_description,
            'image_url', p.product_image_url,
            'is_active', p.product_is_active,
            'created_at', p.product_created_at,
            'updated_at', p.product_updated_at,
            'created_by', p.product_created_by,
            'updated_by', p.product_updated_by
        ) AS product,
        COALESCE((
            SELECT json_agg(json_build_object(
                'id', c.category_id,
                'name', c.category_name
            )) 
            FROM categories_data c WHERE c.product_id = p.product_id
        ), '[]') AS categories,
        COALESCE((
            SELECT json_agg(json_build_object(
                'id', v.variant_id,
                'name', v.variant_name,
                'mrp', v.variant_mrp,
                'discount_price', v.variant_discount_price,
                'size', v.variant_size,
                'color', v.variant_color,
                'quantity', v.variant_quantity,
                'is_active', v.variant_is_active,
                'created_at', v.variant_created_at,
                'updated_at', v.variant_updated_at
            )) 
            FROM variants_data v WHERE v.product_id = p.product_id
        ), '[]') AS variants,
        (SELECT total FROM total_count) AS total
    FROM product_data p;
    `

	rows, err := repo.db.QueryContext(ctx, query, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var productJSON, categoriesJSON, variantsJSON string
		if err := rows.Scan(&productJSON, &categoriesJSON, &variantsJSON, &total); err != nil {
			return nil, 0, fmt.Errorf("could not scan product: %w", err)
		}

		var product entities.ProductDetails
		var categories []*entities.Category
		var variants []*entities.Variant

		// Unmarshal JSON data into the respective structures
		if err := json.Unmarshal([]byte(productJSON), &product); err != nil {
			return nil, 0, fmt.Errorf("could not unmarshal product details: %w", err)
		}
		if err := json.Unmarshal([]byte(categoriesJSON), &categories); err != nil {
			return nil, 0, fmt.Errorf("could not unmarshal categories: %w", err)
		}
		if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
			return nil, 0, fmt.Errorf("could not unmarshal variants: %w", err)
		}

		// Set the fetched categories and variants in the product
		product.Categories = categories
		product.Variants = variants

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, total, nil
}
