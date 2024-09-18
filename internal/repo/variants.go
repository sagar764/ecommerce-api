package repo

import (
	"context"
	"database/sql"
	"ecommerce-api/errortools"
	"ecommerce-api/internal/entities"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (repo *EcommerceRepo) InsertVariant(ctx context.Context, variant entities.Variant) (string, error) {
	var variantID string
	err := repo.db.QueryRowContext(ctx, `
		INSERT INTO variants (name, mrp, discount_price, size, color, quantity, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		variant.Name, variant.MRP, variant.DiscountPrice, variant.Size, variant.Color, variant.Quantity, true).
		Scan(&variantID)

	if err != nil {
		return "", fmt.Errorf("could not insert variant: %w", err)
	}

	return variantID, nil
}

func (repo *EcommerceRepo) UpdateVariant(ctx context.Context, variant entities.Variant) error {
	_, err := repo.db.ExecContext(ctx, `
		UPDATE variants
		SET name = $1, mrp = $2, discount_price = $3, size = $4, color = $5, quantity = $6
		WHERE id = $7 AND is_active = TRUE`,
		variant.Name, variant.MRP, variant.DiscountPrice, variant.Size, variant.Color, variant.Quantity, variant.ID)

	if err != nil {
		return fmt.Errorf("could not update variant: %w", err)
	}

	return nil
}

func (repo *EcommerceRepo) DeleteVariant(ctx context.Context, variantID string) error {
	_, err := repo.db.ExecContext(ctx, `
		UPDATE variants
		SET is_active = FALSE
		WHERE id = $1 AND is_active = TRUE`,
		variantID)

	if err != nil {
		return fmt.Errorf("could not delete variant: %w", err)
	}

	return nil
}

func (repo *EcommerceRepo) CheckVariantIsActive(ctx context.Context, variantID string) (bool, error) {
	var isActive bool
	err := repo.db.QueryRowContext(ctx, `
		SELECT is_active
		FROM variants
		WHERE id = $1`,
		variantID).Scan(&isActive)

	if err != nil {
		return false, fmt.Errorf("could not check variant status: %w", err)
	}

	return isActive, nil
}

func (ecommerce *EcommerceRepo) ValidateVariants(ctx context.Context, variantIDs []string) ([]string, error) {
	if len(variantIDs) == 0 {
		return nil, nil
	}

	// Prepare the query
	query := `SELECT id FROM variants WHERE id = ANY($1)`
	rows, err := ecommerce.db.QueryContext(ctx, query, pq.Array(variantIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Collect valid IDs
	validIDs := make([]string, 0, len(variantIDs))
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

	if len(variantIDs) != len(validIDs) {
		return nil, fmt.Errorf("one or more variant IDs are invalid")
	}

	return validIDs, nil
}

func (ecommerce *EcommerceRepo) HasVariants(ctx context.Context, productID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM product_variant_mapping WHERE product_id = $1`
	err := ecommerce.db.QueryRowContext(ctx, query, productID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return count > 0, nil
}

func (repo *EcommerceRepo) UpdateInventoryBatch(ctx context.Context, tx *sql.Tx, items []entities.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	query := `
        UPDATE variants
        SET quantity = variants.quantity - tmp.quantity
        FROM (VALUES `

	var args []interface{}
	valueStrings := make([]string, 0, len(items))

	for i, item := range items {
		parsedID, err := uuid.Parse(item.VariantID)
		if err != nil {
			return fmt.Errorf("invalid variant ID format: %w", err)
		}
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::uuid, $%d::integer)", i*2+1, i*2+2))
		args = append(args, parsedID, item.Quantity)
	}

	// Build the batch update query
	query += strings.Join(valueStrings, ", ")
	query += `) AS tmp(variant_id, quantity)
              WHERE variants.id = tmp.variant_id AND variants.quantity >= tmp.quantity`

	// Use pgx to execute the batch update query
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("could not batch update inventory: %w", err)
	}

	return nil
}

func (repo *EcommerceRepo) CheckInventoryAvailability(ctx context.Context, items []entities.OrderItem) (bool, error) {
	// Build the query to check inventory levels
	query := `
        SELECT id, quantity
        FROM variants
        WHERE id = ANY($1::uuid[])`

	// Prepare the slice of variant IDs as strings
	varIDs := make([]string, len(items))
	for i, item := range items {
		varIDs[i] = item.VariantID
	}

	// Execute the query with pq.Array to handle the slice
	rows, err := repo.db.QueryContext(ctx, query, pq.Array(varIDs))
	if err != nil {
		return false, fmt.Errorf("could not check inventory availability: %w", err)
	}
	defer rows.Close()

	// Map to track inventory quantities
	inventory := make(map[string]int)
	for rows.Next() {
		var variantID string
		var quantity int
		if err := rows.Scan(&variantID, &quantity); err != nil {
			return false, fmt.Errorf("could not scan inventory row: %w", err)
		}
		inventory[variantID] = quantity
	}

	// Validate if the available quantity meets the order requirements
	for _, item := range items {
		if qty, ok := inventory[item.VariantID]; !ok || qty < item.Quantity {
			return false, nil // Item not found or insufficient quantity
		}
	}

	return true, nil
}

func (repo *EcommerceRepo) GetVariantByID(ctx context.Context, variantID string) (*entities.Variant, error) {
	query := `SELECT id, name, mrp, discount_price, size, color, quantity, is_active, created_by, updated_by, created_at, updated_at
	          FROM variants WHERE id = $1`

	var variant entities.Variant
	err := repo.db.QueryRowContext(ctx, query, variantID).Scan(&variant.ID, &variant.Name, &variant.MRP, &variant.DiscountPrice,
		&variant.Size, &variant.Color, &variant.Quantity, &variant.IsActive, &variant.CreatedBy, &variant.UpdatedBy, &variant.CreatedAt, &variant.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errortools.New("variant_not_found", errortools.WithDetail("Variant with ID "+variantID+" not found"))
		}
		return nil, fmt.Errorf("could not fetch variant: %w", err)
	}

	return &variant, nil
}

func (repo *EcommerceRepo) FetchVariants(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.Variant, int, error) {
	baseQuery := `
        SELECT id, name, mrp, discount_price, size, color, quantity, is_active, 
               COALESCE(created_by::text, '') AS created_by, 
               COALESCE(updated_by::text, '') AS updated_by, 
               created_at, updated_at
        FROM variants`

	var queryParams []interface{}
	if searchQuery != "" {
		baseQuery += " WHERE name ILIKE '%' || $1 || '%' "
		queryParams = append(queryParams, searchQuery)
	}

	// Ensure the LIMIT and OFFSET parameters are at the end of the query
	baseQuery += " ORDER BY mrp ASC LIMIT $" + fmt.Sprint(len(queryParams)+1) + " OFFSET $" + fmt.Sprint(len(queryParams)+2)
	queryParams = append(queryParams, limit, offset)

	// Prepare statement and execute query
	rows, err := repo.db.QueryContext(ctx, baseQuery, queryParams...)
	if err != nil {
		return nil, 0, fmt.Errorf("could not fetch variants: %w", err)
	}
	defer rows.Close()

	// Parse rows to get variants
	var variants []*entities.Variant
	for rows.Next() {
		var variant entities.Variant
		if err := rows.Scan(&variant.ID, &variant.Name, &variant.MRP, &variant.DiscountPrice,
			&variant.Size, &variant.Color, &variant.Quantity, &variant.IsActive,
			&variant.CreatedBy, &variant.UpdatedBy, &variant.CreatedAt, &variant.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("could not scan variant: %w", err)
		}
		variants = append(variants, &variant)
	}

	// Get the total count of variants matching the search query
	countQuery := `SELECT COUNT(*) FROM variants`
	if searchQuery != "" {
		countQuery += " WHERE name ILIKE '%' || $1 || '%'"
	}
	var total int
	if searchQuery != "" {
		err = repo.db.QueryRowContext(ctx, countQuery, searchQuery).Scan(&total)
	} else {
		err = repo.db.QueryRowContext(ctx, countQuery).Scan(&total)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("could not count variants: %w", err)
	}

	return variants, total, nil
}
