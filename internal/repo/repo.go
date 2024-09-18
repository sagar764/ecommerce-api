package repo

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/entities"
)

type EcommerceRepo struct {
	db *sql.DB
}

type EcommerceRepoImply interface {
	InsertCategory(ctx context.Context, category entities.Category) (string, error)
	UpdateCategory(ctx context.Context, category entities.Category) error
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategoryDetails(ctx context.Context, categoryID string) (*entities.CategoryDetails, error)
	FetchCategories(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.CategoryDetails, int, error)
	ValidateChildCategories(ctx context.Context, categoryIDs []string) ([]string, error)
	CheckCategoryIsActive(ctx context.Context, categoryID string) (bool, error)

	InsertProduct(ctx context.Context, product entities.Product) (string, error)
	UpdateProduct(ctx context.Context, product entities.Product) error
	DeleteProduct(ctx context.Context, productID string) error
	GetProductDetails(ctx context.Context, productID string) (*entities.ProductDetails, error)
	FetchProducts(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.ProductDetails, int, error)
	ValidateProducts(ctx context.Context, productIDs []string) ([]string, error)
	HasActiveProducts(ctx context.Context, categoryID string) (bool, error)
	CheckProductIsActive(ctx context.Context, productID string) (bool, error)

	InsertVariant(ctx context.Context, variant entities.Variant) (string, error)
	UpdateVariant(ctx context.Context, variant entities.Variant) error
	DeleteVariant(ctx context.Context, variantID string) error
	GetVariantByID(ctx context.Context, variantID string) (*entities.Variant, error)
	FetchVariants(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.Variant, int, error)
	ValidateVariants(ctx context.Context, variantIDs []string) ([]string, error)
	HasVariants(ctx context.Context, productID string) (bool, error)
	CheckVariantIsActive(ctx context.Context, variantID string) (bool, error)
	UpdateInventoryBatch(ctx context.Context, tx *sql.Tx, items []entities.OrderItem) error
	CheckInventoryAvailability(ctx context.Context, items []entities.OrderItem) (bool, error)

	CreateOrderTransaction(ctx context.Context, order entities.Order) (string, error)
	InsertOrder(ctx context.Context, tx *sql.Tx, order entities.Order) (string, error)
	InsertOrderItemsBatch(ctx context.Context, tx *sql.Tx, orderID string, items []entities.OrderItem) error
	FetchOrderWithDetails(ctx context.Context, orderID string) (*entities.Order, error)
	FetchOrdersWithDetails(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.Order, int, error)
}

func NewEcommerceRepo(db *sql.DB) EcommerceRepoImply {
	return &EcommerceRepo{db: db}
}
