package usecases

import (
	"context"
	"ecommerce-api/errortools"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"ecommerce-api/utils"
	"fmt"
)

func (ecommerce *EcommerceUseCases) CreateProduct(ctx context.Context, product entities.Product) (string, error) {
	if err := validateProductData(product); err != nil {
		return "", err
	}

	validVariants, err := ecommerce.repo.ValidateVariants(ctx, product.ChildVariants)
	if err != nil {
		return "", err
	}
	product.ChildVariants = validVariants

	return ecommerce.repo.InsertProduct(ctx, product)
}

func (ecommerce *EcommerceUseCases) UpdateProduct(ctx context.Context, product entities.Product) error {
	validationErr := errortools.Init()
	if !utils.ValidateUUID(product.ID) {
		validationErr.AddValidationError(
			consts.ProductID,
			errortools.InvalidUUID,
			consts.ProductID,
		)
		return validationErr
	}
	if err := validateProductData(product); err != nil {
		return err
	}

	validVariants, err := ecommerce.repo.ValidateVariants(ctx, product.ChildVariants)
	if err != nil {
		return err
	}
	product.ChildVariants = validVariants

	return ecommerce.repo.UpdateProduct(ctx, product)
}

func (ecommerce *EcommerceUseCases) DeleteProduct(ctx context.Context, productID string) error {
	validationErr := errortools.Init()

	// Validate product ID
	if !utils.ValidateUUID(productID) {
		validationErr.AddValidationError(
			consts.ProductID,
			errortools.InvalidUUID,
			consts.ProductID,
		)
		return validationErr
	}

	// Check if the product is active
	isActive, err := ecommerce.repo.CheckProductIsActive(ctx, productID)
	if err != nil {
		return err
	}
	if !isActive {
		validationErr.AddValidationError(
			consts.Product,
			errortools.Inactive,
			consts.Product,
		)
		return validationErr
	}

	// Check if there are any variants associated with the product
	hasVariants, err := ecommerce.repo.HasVariants(ctx, productID)
	if err != nil {
		return err
	}
	if hasVariants {
		validationErr.AddValidationError(
			consts.Product,
			errortools.HasVariants,
		)
		return validationErr
	}

	// Proceed with soft delete (set is_active to false)
	err = ecommerce.repo.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

func (ecommerce *EcommerceUseCases) GetProductByID(ctx context.Context, productID string) (*entities.ProductDetails, error) {
	validationErr := errortools.Init()

	if !utils.ValidateUUID(productID) {
		validationErr.AddValidationError(
			consts.ProductID,
			errortools.InvalidUUID,
			consts.ProductID,
		)
		return nil, validationErr
	}

	product, err := ecommerce.repo.GetProductDetails(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch product: %w", err)
	}
	return product, nil
}

func validateProductData(product entities.Product) error {
	validationErr := errortools.Init()
	if !utils.ValidateName(product.Name) {
		validationErr.AddValidationError(
			consts.ProductName,
			errortools.Invalid,
			consts.ProductName,
		)
	}
	if err := utils.ValidateURL(product.ImageURL); err != nil {
		validationErr.AddValidationError(
			consts.ProductUrl,
			errortools.Invalid,
			consts.ProductUrl,
		)
	}

	if !validationErr.Nil() {
		return validationErr
	}

	return nil
}

func (ecommerce *EcommerceUseCases) GetAllProducts(ctx context.Context, searchQuery string, limit, page int) ([]*entities.ProductDetails, *entities.MetaData, error) {
	// Handle pagination and limit
	page, limit = utils.Paginate(page, limit, consts.DefaultLimit)
	offset := (page - 1) * limit

	// Fetch products from the repository with search, pagination, and sorting
	products, total, err := ecommerce.repo.FetchProducts(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch products: %w", err)
	}

	// Prepare metadata for pagination
	metaData := &entities.MetaData{
		Total:       total,
		PerPage:     limit,
		CurrentPage: page,
	}

	// Populate meta data info (if required)
	metaData = utils.MetaDataInfo(metaData)

	return products, metaData, nil
}
