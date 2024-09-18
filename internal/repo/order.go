package repo

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func (repo *EcommerceRepo) CreateOrderTransaction(ctx context.Context, order entities.Order) (string, error) {
	// Begin a transaction
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("could not begin transaction: %w", err)
	}

	// Defer rollback if the transaction fails
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Check inventory availability
	if available, err := repo.CheckInventoryAvailability(ctx, order.Items); err != nil {
		return "", err
	} else if !available {
		return "", fmt.Errorf("insufficient inventory for one or more items")
	}

	// Update inventory for each item in a batch query (optimized)
	if err := repo.UpdateInventoryBatch(ctx, tx, order.Items); err != nil {
		return "", err
	}

	// Insert the order into the database
	orderID, err := repo.InsertOrder(ctx, tx, order)
	if err != nil {
		return "", err
	}

	// Insert all order items in a single batch insert
	if err := repo.InsertOrderItemsBatch(ctx, tx, orderID, order.Items); err != nil {
		return "", err
	}

	return orderID, nil
}

func (repo *EcommerceRepo) InsertOrder(ctx context.Context, tx *sql.Tx, order entities.Order) (string, error) {
	var orderID string
	err := tx.QueryRowContext(ctx, `
        INSERT INTO orders (status, order_total)
        VALUES ($1, $2)
        RETURNING id`,
		consts.OrderStatusAccepted, order.Total).
		Scan(&orderID)

	if err != nil {
		return "", fmt.Errorf("could not insert order: %w", err)
	}

	return orderID, nil
}

func (repo *EcommerceRepo) InsertOrderItemsBatch(ctx context.Context, tx *sql.Tx, orderID string, items []entities.OrderItem) error {
	if len(items) == 0 {
		return nil // No items to insert
	}

	query := `
        INSERT INTO order_items (order_id, variant_id, quantity, price)
        VALUES `

	// Build the batch insert placeholders and arguments
	var args []interface{}
	valueStrings := make([]string, 0, len(items))

	for i, item := range items {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		args = append(args, orderID, item.VariantID, item.Quantity, item.Price)
	}

	// Final query with all the values
	query += strings.Join(valueStrings, ", ")

	// Execute the batch insert query
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("could not batch insert order items: %w", err)
	}

	return nil
}

func (repo *EcommerceRepo) FetchOrderWithDetails(ctx context.Context, orderID string) (*entities.Order, error) {
	// Fetch order details
	orderQuery := `
        SELECT id, status, order_total
        FROM orders
        WHERE id = $1`

	var order entities.Order
	err := repo.db.QueryRowContext(ctx, orderQuery, orderID).Scan(&order.ID, &order.Status, &order.Total)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("could not fetch order: %w", err)
	}

	// Fetch order items with variant and product names
	itemsQuery := `
        SELECT oi.variant_id, oi.quantity, oi.price,
               v.name AS variant_name,
               p.name AS product_name
        FROM order_items oi
        JOIN variants v ON oi.variant_id = v.id
        JOIN product_variant_mapping pvm ON v.id = pvm.variant_id
        JOIN products p ON pvm.product_id = p.id
        WHERE oi.order_id = $1`

	rows, err := repo.db.QueryContext(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch order items: %w", err)
	}
	defer rows.Close()

	var items []entities.OrderItem
	for rows.Next() {
		var item entities.OrderItem
		if err := rows.Scan(&item.VariantID, &item.Quantity, &item.Price, &item.VariantName, &item.ProductName); err != nil {
			return nil, fmt.Errorf("could not scan order item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	order.Items = items
	return &order, nil
}

func (repo *EcommerceRepo) FetchOrdersWithDetails(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.Order, int, error) {
	// Define the base query to fetch order details
	baseQuery := `
        SELECT o.id, o.status, o.order_total
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id
        LEFT JOIN product_variant_mapping pvm ON oi.variant_id = pvm.variant_id
        LEFT JOIN products p ON pvm.product_id = p.id`

	// Count query to get total records before applying limit and offset
	countQuery := `SELECT COUNT(DISTINCT o.id) FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		LEFT JOIN product_variant_mapping pvm ON oi.variant_id = pvm.variant_id
		LEFT JOIN products p ON pvm.product_id = p.id`

	// Apply search filtering if provided
	if searchQuery != "" {
		baseQuery += " WHERE p.name ILIKE '%' || $1 || '%' "
		baseQuery += " GROUP BY o.id LIMIT $2 OFFSET $3"

		countQuery += " WHERE p.name ILIKE '%' || $1 || '%'"
	} else {
		baseQuery += " GROUP BY o.id LIMIT $1 OFFSET $2"
	}

	// Get total records count
	var totalRecords int
	if searchQuery != "" {
		err := repo.db.QueryRowContext(ctx, countQuery, searchQuery).Scan(&totalRecords)
		if err != nil {
			return nil, 0, fmt.Errorf("could not fetch total record count: %w", err)
		}
	} else {
		err := repo.db.QueryRowContext(ctx, countQuery).Scan(&totalRecords)
		if err != nil {
			return nil, 0, fmt.Errorf("could not fetch total record count: %w", err)
		}
	}

	// Execute the query with limit and offset
	var orderRows *sql.Rows
	var err error
	if searchQuery != "" {
		orderRows, err = repo.db.QueryContext(ctx, baseQuery, searchQuery, limit, offset)
	} else {
		orderRows, err = repo.db.QueryContext(ctx, baseQuery, limit, offset)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("could not fetch orders: %w", err)
	}
	defer orderRows.Close()

	// Parse the rows to get order details
	var orders []*entities.Order
	orderMap := make(map[string]*entities.Order)
	for orderRows.Next() {
		var order entities.Order
		if err := orderRows.Scan(&order.ID, &order.Status, &order.Total); err != nil {
			return nil, 0, fmt.Errorf("could not scan order: %w", err)
		}
		orders = append(orders, &order)
		orderMap[order.ID] = &order
	}

	// If no orders, return early
	if len(orders) == 0 {
		return orders, totalRecords, nil
	}

	// Fetch all order items for the selected orders in a single query
	orderIDs := make([]interface{}, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}

	// Use SQL's IN clause to fetch all items for the selected orders
	itemsQuery := `
		SELECT oi.order_id, oi.variant_id, oi.quantity, oi.price,
		       v.name AS variant_name,
		       p.name AS product_name
		FROM order_items oi
		JOIN variants v ON oi.variant_id = v.id
		JOIN product_variant_mapping pvm ON v.id = pvm.variant_id
		JOIN products p ON pvm.product_id = p.id
		WHERE oi.order_id = ANY($1)`

	itemRows, err := repo.db.QueryContext(ctx, itemsQuery, pq.Array(orderIDs))
	if err != nil {
		return nil, 0, fmt.Errorf("could not fetch order items: %w", err)
	}
	defer itemRows.Close()

	// Parse the order items and assign them to corresponding orders
	for itemRows.Next() {
		var item entities.OrderItem
		var orderID string
		if err := itemRows.Scan(&orderID, &item.VariantID, &item.Quantity, &item.Price, &item.VariantName, &item.ProductName); err != nil {
			return nil, 0, fmt.Errorf("could not scan order item: %w", err)
		}

		// Assign the item to the correct order in the map
		if order, exists := orderMap[orderID]; exists {
			order.Items = append(order.Items, item)
		}
	}

	return orders, totalRecords, nil
}
