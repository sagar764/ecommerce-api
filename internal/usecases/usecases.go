package usecases

import (
	"context"
	"ecommerce-api/internal/entities"
	"ecommerce-api/internal/repo"
)

type EcommerceUseCases struct {
	repo repo.EcommerceRepoImply
}

type EcommerceUseCaseImply interface {
	CreateCategory(ctx context.Context, category entities.Category) (string, error)
	UpdateCategory(ctx context.Context, category entities.Category) error
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategoryByID(ctx context.Context, categoryID string) (*entities.CategoryDetails, error)
	GetAllCategories(ctx context.Context, searchQuery string, limit, page int) ([]*entities.CategoryDetails, *entities.MetaData, error)

	CreateProduct(ctx context.Context, product entities.Product) (string, error)
	UpdateProduct(ctx context.Context, product entities.Product) error
	DeleteProduct(ctx context.Context, productID string) error
	GetProductByID(ctx context.Context, productID string) (*entities.ProductDetails, error)
	GetAllProducts(ctx context.Context, searchQuery string, limit, page int) ([]*entities.ProductDetails, *entities.MetaData, error)

	CreateVariant(ctx context.Context, variant entities.Variant) (string, error)
	UpdateVariant(ctx context.Context, variant entities.Variant) error
	DeleteVariant(ctx context.Context, variantID string) error
	GetVariantByID(ctx context.Context, variantID string) (*entities.Variant, error)
	GetAllVariants(ctx context.Context, searchQuery string, limit, page int) ([]*entities.Variant, *entities.MetaData, error)

	CreateOrder(ctx context.Context, order entities.Order) (string, error)
	GetOrderByID(ctx context.Context, orderID string) (*entities.Order, error)
	GetAllOrders(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.Order, *entities.MetaData, error)
}

// NewUserUseCases
func NewEcommerceUseCases(ecommerceRepo repo.EcommerceRepoImply) EcommerceUseCaseImply {
	return &EcommerceUseCases{
		repo: ecommerceRepo,
	}
}
