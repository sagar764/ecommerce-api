package usecases

import (
	"context"
	"ecommerce-api/errortools"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"ecommerce-api/utils"
	"fmt"
)

func (ecommerce *EcommerceUseCases) CreateOrder(ctx context.Context, order entities.Order) (string, error) {
	// Validate order data
	if err := validateOrderData(order); err != nil {
		return "", err
	}

	// Call the repository to handle transaction and business logic
	orderID, err := ecommerce.repo.CreateOrderTransaction(ctx, order)
	if err != nil {
		return "", fmt.Errorf("could not create order: %w", err)
	}

	return orderID, nil
}

func (ecommerce *EcommerceUseCases) GetOrderByID(ctx context.Context, orderID string) (*entities.Order, error) {
	validationErr := errortools.Init()

	if !utils.ValidateUUID(orderID) {
		validationErr.AddValidationError(
			consts.OrderID,
			errortools.InvalidUUID,
			consts.OrderID,
		)
		return nil, validationErr
	}

	order, err := ecommerce.repo.FetchOrderWithDetails(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch order: %w", err)
	}
	return order, nil
}

func (ecommerce *EcommerceUseCases) GetAllOrders(ctx context.Context, searchQuery string, limit, offset int) ([]*entities.Order, *entities.MetaData, error) {
	page, limit := utils.Paginate(offset, limit, consts.DefaultLimit)
	// Fetch all orders with the specified search query, limit, and offset

	offset = (page - 1) * limit
	orders, total, err := ecommerce.repo.FetchOrdersWithDetails(ctx, searchQuery, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch orders: %w", err)
	}

	metaData := &entities.MetaData{
		Total:       total,
		PerPage:     limit,
		CurrentPage: page,
	}

	metaData = utils.MetaDataInfo(metaData)

	return orders, metaData, nil
}

func validateOrderData(order entities.Order) error {
	validationErr := errortools.Init()
	if len(order.Items) == 0 {
		validationErr.AddValidationError(
			consts.OrderItems,
			errortools.Required,
			consts.OrderItems,
		)
	}

	if order.Total <= 0 {
		validationErr.AddValidationError(
			consts.OrderTotal,
			errortools.Invalid,
			consts.OrderTotal,
		)
	}

	for _, item := range order.Items {
		if item.Quantity <= 0 {
			validationErr.AddValidationError(
				fmt.Sprintf("%s.%s", item.VariantID, consts.OrderQuantity),
				errortools.Invalid,
				consts.OrderQuantity,
			)
		}
		if item.Price <= 0 {
			validationErr.AddValidationError(
				fmt.Sprintf("%s.%s", item.VariantID, consts.OrderPrice),
				errortools.Invalid,
				consts.OrderPrice,
			)
		}
	}

	if !validationErr.Nil() {
		return validationErr
	}

	return nil
}
