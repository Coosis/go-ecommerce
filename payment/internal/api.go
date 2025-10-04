package internal

import (
	"context"
)

type GetOrderRequest struct {
	OrderID string
}

type OrderItems struct {
	ProductID string
	ProductName string
	SkuID string
	SkuName string

	Qty int32
	UnitPriceCents int32
}

type GetOrderResponse struct {
	OrderID string
	Items []OrderItems
}

type Api interface {
	GetOrder(context.Context, *GetOrderRequest) (*GetOrderResponse, error);
}
