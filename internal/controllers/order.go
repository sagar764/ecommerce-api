package controllers

import (
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"ecommerce-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateOrder handles the creation of a new order
// @Summary Create a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body entities.Order true "Order details"
// @Success 201 {object} utils.APIResponse{data=entities.Order}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/orders [post]
func (ecommerce *EcommerceController) CreateOrder(ctx *gin.Context) {
	var orderRequest entities.Order
	if err := ctx.BindJSON(&orderRequest); err != nil {
		utils.BindingError(ctx, err)
		return
	}

	// Create the order
	orderID, err := ecommerce.useCases.CreateOrder(ctx, orderRequest)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	orderRequest.ID = orderID
	orderRequest.Status = consts.OrderStatusAccepted
	successResponse := utils.SuccessGenerator(orderRequest, consts.Order, consts.OrderCreateSuccess)
	ctx.JSON(http.StatusCreated, successResponse)
}

// GetOrderByID retrieves an order by its ID
// @Summary Get an order by ID
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} utils.APIResponse{data=entities.Order}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/orders/{id} [get]
func (ecommerce *EcommerceController) GetOrderByID(ctx *gin.Context) {
	orderID := ctx.Param("id")

	order, err := ecommerce.useCases.GetOrderByID(ctx, orderID)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	successResponse := utils.SuccessGenerator(order, consts.Order, consts.OrderFetchSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}

// GetAllOrders retrieves a list of orders with pagination and search functionality
// @Summary Get all orders with pagination and search
// @Tags Orders
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(0)
// @Param search query string false "Search query"
// @Success 200 {object} utils.APIResponse{data=map[string]interface{}}
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/{version}/orders [get]
func (ecommerce *EcommerceController) GetAllOrders(ctx *gin.Context) {

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	searchQuery := ctx.Query("search")

	orders, metadata, err := ecommerce.useCases.GetAllOrders(ctx, searchQuery, limit, offset)
	if err != nil {
		utils.ErrorGenerator(ctx, err)
		return
	}

	orderWithMetadat := map[string]interface{}{
		"orders":   orders,
		"metadata": metadata,
	}

	successResponse := utils.SuccessGenerator(orderWithMetadat, consts.Order, consts.OrderFetchSuccess)
	ctx.JSON(http.StatusOK, successResponse)
}
