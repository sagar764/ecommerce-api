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

// CreateProduct creates a new product in the system
// @Summary Create a new product
// @Tags Products
// @Produce json
// @Param product body entities.Product true "Product to create"
// @Success 201 {object} utils.APIResponse{data=entities.Product}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/products [post]
func (ecommerce *EcommerceController) CreateProduct(ctx *gin.Context) {
	var product entities.Product
	if err := ctx.BindJSON(&product); err != nil {
		utils.BindingError(ctx, err)
		return
	}

	productID, err := ecommerce.useCases.CreateProduct(ctx, product)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	product.ID = productID
	successResponse := utils.SuccessGenerator(product, consts.Product, consts.ProductCreateSuccess)
	ctx.JSON(http.StatusCreated, successResponse)
}

// UpdateProduct updates an existing product
// @Summary Update an existing product
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Param product body entities.Product true "Product to update"
// @Success 200 {object} utils.APIResponse{data=entities.Product}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/products/{id} [put]
func (ecommerce *EcommerceController) UpdateProduct(ctx *gin.Context) {
	var product entities.Product
	if err := ctx.BindJSON(&product); err != nil {
		utils.BindingError(ctx, err)
		return
	}

	productID := ctx.Param("id")
	product.ID = productID

	err := ecommerce.useCases.UpdateProduct(ctx, product)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	successResponse := utils.SuccessGenerator(product, consts.Product, consts.ProductUpdateSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// DeleteProduct deletes a product by ID
// @Summary Delete a product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/products/{id} [delete]
func (ecommerce *EcommerceController) DeleteProduct(ctx *gin.Context) {
	productID := ctx.Param("id")

	err := ecommerce.useCases.DeleteProduct(ctx, productID)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	successResponse := utils.SuccessGenerator(nil, consts.Product, consts.ProductDeleteSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// GetProductByID retrieves a product by its ID
// @Summary Retrieve a product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entities.ProductDetails
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/products/{id} [get]
func (ecommerce *EcommerceController) GetProductByID(ctx *gin.Context) {
	productID := ctx.Param("id")

	product, err := ecommerce.useCases.GetProductByID(ctx, productID)
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

	successResponse := utils.SuccessGenerator(product, consts.Product, consts.ProductFetchSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// GetAllProducts retrieves a list of products with pagination and search functionality
// @Summary Retrieve all products with pagination and search
// @Tags Products
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(0)
// @Param search query string false "Search query"
// @Success 200 {object} utils.APIResponse{data=map[string]interface{}}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/products [get]
func (ecommerce *EcommerceController) GetAllProducts(ctx *gin.Context) {
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

	// Call use case to get all products
	products, metadata, err := ecommerce.useCases.GetAllProducts(ctx, searchQuery, limit, page)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	// Combine products and metadata into a single response
	responseData := map[string]interface{}{
		"products": products,
		"metadata": metadata,
	}

	successResponse := utils.SuccessGenerator(responseData, consts.Product, consts.ProductFetchSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}
