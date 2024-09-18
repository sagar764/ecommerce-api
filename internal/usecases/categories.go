package usecases

import (
	"context"
	"ecommerce-api/errortools"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"ecommerce-api/utils"
	"fmt"
)

func (ecommerce *EcommerceUseCases) CreateCategory(ctx context.Context, category entities.Category) (string, error) {
	validationErr := errortools.Init()
	if !utils.ValidateName(category.Name) {
		validationErr.AddValidationError(
			consts.CategoryName,
			errortools.Invalid,
			consts.CategoryName,
		)
		return "", validationErr
	}

	validChildCategories, err := ecommerce.repo.ValidateChildCategories(ctx, category.ChildCategories)
	if err != nil {
		return "", err
	}
	category.ChildCategories = validChildCategories

	validProducts, err := ecommerce.repo.ValidateProducts(ctx, category.Products)
	if err != nil {
		return "", err
	}
	category.Products = validProducts

	return ecommerce.repo.InsertCategory(ctx, category)
}

func (ecommerce *EcommerceUseCases) UpdateCategory(ctx context.Context, category entities.Category) error {
	validationErr := errortools.Init()

	if !utils.ValidateUUID(category.ID) {
		validationErr.AddValidationError(
			consts.CategoryID,
			errortools.InvalidUUID,
			consts.CategoryID,
		)
		return validationErr
	}

	if !utils.ValidateName(category.Name) {
		validationErr.AddValidationError(
			consts.CategoryName,
			errortools.Invalid,
			consts.CategoryName,
		)
		return validationErr
	}

	validChildCategories, err := ecommerce.repo.ValidateChildCategories(ctx, category.ChildCategories)
	if err != nil {
		return err
	}
	category.ChildCategories = validChildCategories

	validProducts, err := ecommerce.repo.ValidateProducts(ctx, category.Products)
	if err != nil {
		return err
	}
	category.Products = validProducts

	return ecommerce.repo.UpdateCategory(ctx, category)
}

func (ecommerce *EcommerceUseCases) DeleteCategory(ctx context.Context, categoryID string) error {
	validationErr := errortools.Init()
	// Validate category ID
	if !utils.ValidateUUID(categoryID) {
		validationErr.AddValidationError(
			consts.CategoryID,
			errortools.InvalidUUID,
			consts.CategoryID,
		)
		return validationErr
	}

	// Validate if the category exists and is active
	isActive, err := ecommerce.repo.CheckCategoryIsActive(ctx, categoryID)
	if err != nil {
		return err
	}
	if !isActive {
		validationErr.AddValidationError(
			consts.Category,
			errortools.Inactive,
			consts.Category,
		)
		return validationErr
	}

	// Check if there are any active products in this category
	hasActiveProducts, err := ecommerce.repo.HasActiveProducts(ctx, categoryID)
	if err != nil {
		return err
	}
	if hasActiveProducts {
		validationErr.AddValidationError(
			consts.Category,
			errortools.HasActiveProduct,
		)
		return validationErr
	}

	// Proceed with soft delete (set is_active to false)
	err = ecommerce.repo.DeleteCategory(ctx, categoryID)
	if err != nil {
		return err
	}

	return nil
}

func (ecommerce *EcommerceUseCases) GetCategoryByID(ctx context.Context, categoryID string) (*entities.CategoryDetails, error) {
	validationErr := errortools.Init()

	if !utils.ValidateUUID(categoryID) {
		validationErr.AddValidationError(
			consts.CategoryID,
			errortools.InvalidUUID,
			consts.CategoryID,
		)
		return nil, validationErr
	}

	category, err := ecommerce.repo.GetCategoryDetails(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch category: %w", err)
	}
	return category, nil
}

func (ecommerce *EcommerceUseCases) GetAllCategories(ctx context.Context, searchQuery string, limit, page int) ([]*entities.CategoryDetails, *entities.MetaData, error) {
	// Handle pagination and limit
	page, limit = utils.Paginate(page, limit, consts.DefaultLimit)
	offset := (page - 1) * limit

	// Fetch categories from the repository with search, pagination, and sorting
	categories, total, err := ecommerce.repo.FetchCategories(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch categories: %w", err)
	}

	// Prepare metadata for pagination
	metaData := &entities.MetaData{
		Total:       total,
		PerPage:     limit,
		CurrentPage: page,
	}

	// Populate meta data info (if required)
	metaData = utils.MetaDataInfo(metaData)

	return categories, metaData, nil
}
