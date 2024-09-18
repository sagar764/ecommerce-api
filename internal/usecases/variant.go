package usecases

import (
	"context"
	"ecommerce-api/errortools"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"ecommerce-api/utils"
	"fmt"
)

func (ecommerce *EcommerceUseCases) CreateVariant(ctx context.Context, variant entities.Variant) (string, error) {

	// Validate the variant data
	if err := validateVariantData(variant); err != nil {
		return "", err
	}

	return ecommerce.repo.InsertVariant(ctx, variant)
}

func (ecommerce *EcommerceUseCases) UpdateVariant(ctx context.Context, variant entities.Variant) error {
	validationErr := errortools.Init()

	// Validate the variant ID
	if !utils.ValidateUUID(variant.ID) {
		validationErr.AddValidationError(
			consts.VariantID,
			errortools.InvalidUUID,
			consts.VariantID,
		)
		return validationErr
	}

	// Validate the variant data
	if err := validateVariantData(variant); err != nil {
		return err
	}

	return ecommerce.repo.UpdateVariant(ctx, variant)
}

func (ecommerce *EcommerceUseCases) DeleteVariant(ctx context.Context, variantID string) error {
	validationErr := errortools.Init()

	// Validate variant ID
	if !utils.ValidateUUID(variantID) {
		validationErr.AddValidationError(
			consts.VariantID,
			errortools.InvalidUUID,
			consts.VariantID,
		)
		return validationErr
	}

	// Check if the variant is active
	isActive, err := ecommerce.repo.CheckVariantIsActive(ctx, variantID)
	if err != nil {
		return err
	}
	if !isActive {
		validationErr.AddValidationError(
			consts.Variant,
			errortools.Inactive,
			consts.Variant,
		)
		return validationErr
	}

	// Proceed with soft delete (set is_active to false)
	return ecommerce.repo.DeleteVariant(ctx, variantID)
}

func (ecommerce *EcommerceUseCases) GetVariantByID(ctx context.Context, variantID string) (*entities.Variant, error) {
	validationErr := errortools.Init()

	// Validate the variant ID
	if !utils.ValidateUUID(variantID) {
		validationErr.AddValidationError(
			consts.VariantID,
			errortools.InvalidUUID,
			consts.VariantID,
		)
		return nil, validationErr
	}
	return ecommerce.repo.GetVariantByID(ctx, variantID)
}

func (ecommerce *EcommerceUseCases) GetAllVariants(ctx context.Context, searchQuery string, limit, page int) ([]*entities.Variant, *entities.MetaData, error) {
	page, limit = utils.Paginate(page, limit, consts.DefaultLimit)
	offset := (page - 1) * limit

	// Fetch variants with search, limit, offset, and sorting
	variants, total, err := ecommerce.repo.FetchVariants(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch variants: %w", err)
	}

	metaData := &entities.MetaData{
		Total:       total,
		PerPage:     limit,
		CurrentPage: page,
	}

	metaData = utils.MetaDataInfo(metaData)

	return variants, metaData, nil
}

func validateVariantData(variant entities.Variant) error {
	validationErr := errortools.Init()

	if !utils.ValidateName(variant.Name) {
		validationErr.AddValidationError(
			consts.VariantName,
			errortools.Invalid,
			consts.VariantName,
		)
	}
	if variant.MRP <= 0 {
		validationErr.AddValidationError(
			consts.VariantMrp,
			errortools.Invalid,
			consts.VariantMrp,
		)
	}
	if variant.DiscountPrice != nil {
		if *variant.DiscountPrice < 0.0 {
			validationErr.AddValidationError(
				consts.VariantDiscountPrice,
				errortools.Invalid,
				consts.VariantDiscountPrice,
			)
		}
	}
	if variant.Quantity <= 0 {
		validationErr.AddValidationError(
			consts.VariantQuantity,
			errortools.Invalid,
			consts.VariantQuantity,
		)
	}

	if !validationErr.Nil() {
		return validationErr
	}

	return nil
}
