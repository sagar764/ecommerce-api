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

func (ecommerce *EcommerceRepo) InsertCategory(ctx context.Context, category entities.Category) (string, error) {
	// Begin a transaction
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

	// Insert category and return the generated id
	var categoryID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO categories (category_name, is_active)
		VALUES ($1, $2)
		RETURNING id`,
		category.Name, true).Scan(&categoryID)

	if err != nil {
		return "", fmt.Errorf("could not insert category: %w", err)
	}

	// Insert child categories
	if len(category.ChildCategories) > 0 {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO category_mapping (parent_category_id, child_category_id)
			SELECT $1, unnest($2::uuid[])
			ON CONFLICT DO NOTHING`,
			categoryID, pq.Array(category.ChildCategories))

		if err != nil {
			return "", fmt.Errorf("could not insert category mappings: %w", err)
		}
	}

	// Insert product-category mappings
	if len(category.Products) > 0 {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO product_category_mapping (product_id, category_id)
			SELECT unnest($1::uuid[]), $2
			ON CONFLICT DO NOTHING`,
			pq.Array(category.Products), categoryID)

		if err != nil {
			return "", fmt.Errorf("could not insert product-category mappings: %w", err)
		}
	}

	return categoryID, nil
}

func (ecommerce *EcommerceRepo) UpdateCategory(ctx context.Context, category entities.Category) error {
	// Begin a transaction
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

	// Update category
	_, err = tx.ExecContext(ctx, `
		UPDATE categories
		SET name = $1, is_active = $2
		WHERE id = $3`,
		category.Name, true, category.ID)

	if err != nil {
		return fmt.Errorf("could not update category: %w", err)
	}

	// Update child categories
	if len(category.ChildCategories) > 0 {
		_, err = tx.ExecContext(ctx, `
			DELETE FROM category_mapping
			WHERE parent_category_id = $1
			AND child_category_id = ANY($2::uuid[])`,
			category.ID, pq.Array(category.ChildCategories))
		if err != nil {
			return fmt.Errorf("could not delete old category mappings: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO category_mapping (parent_category_id, child_category_id)
			SELECT $1, unnest($2::uuid[])`,
			category.ID, pq.Array(category.ChildCategories))
		if err != nil {
			return fmt.Errorf("could not insert new category mappings: %w", err)
		}
	}

	// Update product-category mappings
	if len(category.Products) > 0 {
		_, err = tx.ExecContext(ctx, `
			DELETE FROM product_category_mapping
			WHERE category_id = $1
			AND product_id = ANY($2::uuid[])`,
			category.ID, pq.Array(category.Products))
		if err != nil {
			return fmt.Errorf("could not delete old product-category mappings: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO product_category_mapping (product_id, category_id)
			SELECT unnest($1::uuid[]), $2`,
			pq.Array(category.Products), category.ID)
		if err != nil {
			return fmt.Errorf("could not insert new product-category mappings: %w", err)
		}
	}

	return nil
}

func (ecommerce *EcommerceRepo) DeleteCategory(ctx context.Context, categoryID string) error {
	_, err := ecommerce.db.ExecContext(ctx, `
		UPDATE categories 
		SET is_active = FALSE 
		WHERE id = $1`,
		categoryID)

	if err != nil {
		return fmt.Errorf("could not soft delete category: %w", err)
	}

	return nil
}

func (repo *EcommerceRepo) GetCategoryDetails(ctx context.Context, categoryID string) (*entities.CategoryDetails, error) {
	query := `
    WITH category_data AS (
        SELECT
            c.id AS category_id,
            c.name AS category_name,
            c.is_active AS category_is_active,
            c.created_at AS category_created_at,
            c.updated_at AS category_updated_at,
            c.created_by AS category_created_by,
            c.updated_by AS category_updated_by
        FROM categories c
        WHERE c.id = $1
    ),
    child_categories AS (
        SELECT
            c.id AS category_id,
            c.name AS category_name
        FROM categories c
        JOIN category_mapping cm ON c.id = cm.child_category_id
        WHERE cm.parent_category_id = $1
    ),
    products AS (
        SELECT
            p.id AS product_id,
            p.name AS product_name,
            p.description AS product_description,
            p.image_url AS product_image_url
        FROM products p
        JOIN product_category_mapping pcm ON p.id = pcm.product_id
        WHERE pcm.category_id = $1
    )

    SELECT
        json_build_object(
            'id', c.category_id,
            'name', c.category_name,
            'is_active', c.category_is_active,
            'created_at', c.category_created_at,
            'updated_at', c.category_updated_at,
            'created_by', c.category_created_by,
            'updated_by', c.category_updated_by
        ) AS category,
        COALESCE((SELECT json_agg(json_build_object(
            'id', c.category_id,
            'name', c.category_name
        )) FROM child_categories c), '[]') AS child_categories,
        COALESCE((SELECT json_agg(json_build_object(
            'id', p.product_id,
            'name', p.product_name,
            'description', p.product_description,
            'image_url', p.product_image_url
        )) FROM products p), '[]') AS products
    FROM category_data c
    `

	row := repo.db.QueryRowContext(ctx, query, categoryID)

	var categoryJSON, childCategoriesJSON, productsJSON string
	err := row.Scan(&categoryJSON, &childCategoriesJSON, &productsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errortools.New("category_not_found", errortools.WithDetail("Category with ID "+categoryID+" not found"))
		}
		return nil, fmt.Errorf("could not fetch category details: %w", err)
	}

	var category entities.CategoryDetails
	var childCategories []*entities.Category
	var products []*entities.Product

	// Unmarshal JSON responses
	if err := json.Unmarshal([]byte(categoryJSON), &category); err != nil {
		return nil, fmt.Errorf("could not unmarshal category details: %w", err)
	}
	if err := json.Unmarshal([]byte(childCategoriesJSON), &childCategories); err != nil {
		return nil, fmt.Errorf("could not unmarshal child categories: %w", err)
	}
	if err := json.Unmarshal([]byte(productsJSON), &products); err != nil {
		return nil, fmt.Errorf("could not unmarshal products: %w", err)
	}

	category.ChildCategories = childCategories
	category.Products = products

	return &category, nil
}

func (repo *EcommerceRepo) FetchCategories(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.CategoryDetails, int, error) {
	// Adjust the search query for case-insensitive search
	searchQuery = "%" + searchQuery + "%"

	query := `
    WITH total_count AS (
        SELECT COUNT(*) AS total FROM categories WHERE name ILIKE $1
    ),
    category_data AS (
        SELECT
            c.id AS category_id,
            c.name AS category_name,
            c.is_active AS category_is_active,
            c.created_at AS category_created_at,
            c.updated_at AS category_updated_at,
            c.created_by AS category_created_by,
            c.updated_by AS category_updated_by
        FROM categories c
        WHERE c.name ILIKE $1
        ORDER BY c.created_at DESC
        LIMIT $2 OFFSET $3
    )

    SELECT
        json_build_object(
            'id', c.category_id,
            'name', c.category_name,
            'is_active', c.category_is_active,
            'created_at', c.category_created_at,
            'updated_at', c.category_updated_at,
            'created_by', c.category_created_by,
            'updated_by', c.category_updated_by
        ) AS category,
        COALESCE((SELECT json_agg(json_build_object(
            'id', cc.id,
            'name', cc.name
        )) FROM category_mapping cm JOIN categories cc ON cm.child_category_id = cc.id WHERE cm.parent_category_id = c.category_id), '[]') AS child_categories,
        COALESCE((SELECT json_agg(json_build_object(
            'id', p.id,
            'name', p.name,
            'description', p.description,
            'image_url', p.image_url
        )) FROM products p JOIN product_category_mapping pcm ON p.id = pcm.product_id WHERE pcm.category_id = c.category_id), '[]') AS products,
        t.total
    FROM category_data c
    JOIN total_count t ON TRUE
    `

	rows, err := repo.db.QueryContext(ctx, query, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	var categories []*entities.CategoryDetails
	var total int

	for rows.Next() {
		var categoryJSON, childCategoriesJSON, productsJSON string
		if err := rows.Scan(&categoryJSON, &childCategoriesJSON, &productsJSON, &total); err != nil {
			return nil, 0, fmt.Errorf("could not scan category: %w", err)
		}

		var category entities.CategoryDetails
		var childCategories []*entities.Category
		var products []*entities.Product

		// Unmarshal JSON data into the respective structures
		if err := json.Unmarshal([]byte(categoryJSON), &category); err != nil {
			return nil, 0, fmt.Errorf("could not unmarshal category details: %w", err)
		}
		if err := json.Unmarshal([]byte(childCategoriesJSON), &childCategories); err != nil {
			return nil, 0, fmt.Errorf("could not unmarshal child categories: %w", err)
		}
		if err := json.Unmarshal([]byte(productsJSON), &products); err != nil {
			return nil, 0, fmt.Errorf("could not unmarshal products: %w", err)
		}

		// Set the fetched child categories and products in the category
		category.ChildCategories = childCategories
		category.Products = products

		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over rows: %w", err)
	}

	return categories, total, nil
}

func (ecommerce *EcommerceRepo) ValidateChildCategories(ctx context.Context, categoryIDs []string) ([]string, error) {
	if len(categoryIDs) == 0 {
		return nil, nil
	}

	// Prepare the query
	query := `SELECT id FROM categories WHERE id = ANY($1) AND is_active = TRUE`
	rows, err := ecommerce.db.Query(query, pq.Array(categoryIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Collect valid IDs
	validIDs := make([]string, 0, len(categoryIDs))
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

	if len(validIDs) != len(categoryIDs) {
		return nil, fmt.Errorf("one or more category IDs are invalid")
	}

	return validIDs, nil
}

func (ecommerce *EcommerceRepo) CheckCategoryIsActive(ctx context.Context, categoryID string) (bool, error) {
	var isActive bool
	err := ecommerce.db.QueryRowContext(ctx, `
		SELECT is_active FROM categories WHERE id = $1`,
		categoryID).Scan(&isActive)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("category not found")
		}
		return false, fmt.Errorf("could not check category status: %w", err)
	}

	return isActive, nil
}
