package controllers

import (
	"ecommerce-api/errortools"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"ecommerce-api/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateVariant handles the creation of a new product variant
// @Summary Create a new product variant
// @Tags Variants
// @Produce json
// @Param variant body entities.Variant true "Variant details"
// @Success 201 {object} utils.APIResponse{data=entities.Variant}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/variants [post]
func (ecommerce *EcommerceController) CreateVariant(ctx *gin.Context) {
	var variant entities.Variant
	if err := ctx.BindJSON(&variant); err != nil {
		utils.BindingError(ctx, err)
		return
	}

	variantID, err := ecommerce.useCases.CreateVariant(ctx, variant)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	variant.ID = variantID
	successResponse := utils.SuccessGenerator(variant, consts.Variant, consts.VariantCreateSuccess)
	ctx.JSON(http.StatusCreated, successResponse)
}

// UpdateVariant updates an existing product variant
// @Summary Update an existing product variant
// @Tags Variants
// @Produce json
// @Param id path string true "Variant ID"
// @Param variant body entities.Variant true "Variant details"
// @Success 200 {object} utils.APIResponse{data=entities.Variant}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/variants/{id} [put]
func (ecommerce *EcommerceController) UpdateVariant(ctx *gin.Context) {
	var variant entities.Variant
	if err := ctx.BindJSON(&variant); err != nil {
		utils.BindingError(ctx, err)
		return
	}

	variantID := ctx.Param("id")
	variant.ID = variantID

	err := ecommerce.useCases.UpdateVariant(ctx, variant)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	successResponse := utils.SuccessGenerator(variant, consts.Variant, consts.VariantUpdateSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// DeleteVariant deletes a product variant by ID
// @Summary Delete a product variant by ID
// @Tags Variants
// @Produce json
// @Param id path string true "Variant ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/variants/{id} [delete]
func (ecommerce *EcommerceController) DeleteVariant(ctx *gin.Context) {
	variantID := ctx.Param("id")

	err := ecommerce.useCases.DeleteVariant(ctx, variantID)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	successResponse := utils.SuccessGenerator(nil, consts.Variant, consts.VariantDeleteSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// GetVariantByID retrieves a product variant by ID
// @Summary Retrieve a product variant by ID
// @Tags Variants
// @Produce json
// @Param id path string true "Variant ID"
// @Success 200 {object} utils.APIResponse{data=entities.Variant}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/variants/{id} [get]
func (ecommerce *EcommerceController) GetVariantByID(ctx *gin.Context) {
	variantID := ctx.Param("id")

	variant, err := ecommerce.useCases.GetVariantByID(ctx, variantID)
	if err != nil {
		var e *errortools.Error
		if errors.As(err, &e) {
			utils.ErrorGenerator(ctx, e)
		} else {
			internalErr := errortools.New(errortools.InternalServerErrCode, errortools.WithDetail(err.Error()))
			utils.ErrorGenerator(ctx, internalErr)
		}
		return
	}

	// Return successful response
	successResponse := utils.SuccessGenerator(variant, consts.Variant, consts.VariantFetchSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// GetAllVariants retrieves a list of variants with pagination and search functionality
// @Summary Get all variants with pagination and search
// @Tags Variants
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(0)
// @Param search query string false "Search query"
// @Success 200 {object} utils.APIResponse{data=map[string]interface{}}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/variants [get]
func (ecommerce *EcommerceController) GetAllVariants(ctx *gin.Context) {
	// Parse query parameters
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	searchQuery := ctx.Query("search")

	// Call use case to get all variants
	variants, metadata, err := ecommerce.useCases.GetAllVariants(ctx, searchQuery, limit, page)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	// Combine variants and metadata into a single response
	responseData := map[string]interface{}{
		"variants": variants,
		"metadata": metadata,
	}

	successResponse := utils.SuccessGenerator(responseData, consts.Variant, consts.VariantFetchSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}
