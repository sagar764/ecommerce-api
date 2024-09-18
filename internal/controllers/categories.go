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

// CreateCategory creates a new category
// @Summary Create a new category
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body entities.Category true "Category data"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/category [post]
func (ecommerce *EcommerceController) CreateCategory(ctx *gin.Context) {
	var category entities.Category
	if err := ctx.BindJSON(&category); err != nil {
		utils.BindingError(ctx, err)
		return
	}

	categoryID, err := ecommerce.useCases.CreateCategory(ctx, category)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	category.ID = categoryID
	successResponse := utils.SuccessGenerator(category, consts.Category, consts.CategoryCreateSuccess)
	ctx.JSON(http.StatusCreated, successResponse)
}

// UpdateCategory updates an existing category
// @Summary Update an existing category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body entities.Category true "Category data"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/category/{id} [put]
func (ecommerce *EcommerceController) UpdateCategory(ctx *gin.Context) {
	var category entities.Category
	if err := ctx.BindJSON(&category); err != nil {
		utils.BindingError(ctx, err)
		return
	}

	categoryID := ctx.Param("id")

	category.ID = categoryID
	err := ecommerce.useCases.UpdateCategory(ctx, category)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	successResponse := utils.SuccessGenerator(category, consts.Category, consts.CategoryUpdateSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// DeleteCategory deletes an existing category
// @Summary Delete an existing category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/category/{id} [delete]
func (ecommerce *EcommerceController) DeleteCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")

	err := ecommerce.useCases.DeleteCategory(ctx, categoryID)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	successResponse := utils.SuccessGenerator(nil, consts.Category, consts.CategoryDeleteSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// GetCategoryByID retrieves details of a specific category by its ID
// @Summary Get details of a specific category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} utils.APIResponse{data=entities.CategoryDetails}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/category/{id} [get]
func (ecommerce *EcommerceController) GetCategoryByID(ctx *gin.Context) {
	categoryID := ctx.Param("id")

	category, err := ecommerce.useCases.GetCategoryByID(ctx, categoryID)
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

	successResponse := utils.SuccessGenerator(category, consts.Category, consts.CategoryFetchSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// GetAllCategories retrieves a paginated list of categories
// @Summary Get all categories with optional pagination and search
// @Tags Categories
// @Accept json
// @Produce json
// @Param limit query int false "Number of categories to retrieve" default(10)
// @Param page query int false "Page number for pagination" default(0)
// @Param search query string false "Search term to filter categories"
// @Success 200 {object} utils.APIResponse{data=map[string]interface{}}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/categories [get]
func (ecommerce *EcommerceController) GetAllCategories(ctx *gin.Context) {
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

	// Call use case to get all categories
	categories, metadata, err := ecommerce.useCases.GetAllCategories(ctx, searchQuery, limit, page)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	// Combine categories and metadata into a single response
	responseData := map[string]interface{}{
		"categories": categories,
		"metadata":   metadata,
	}

	successResponse := utils.SuccessGenerator(responseData, consts.Category, consts.CategoryFetchSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}
