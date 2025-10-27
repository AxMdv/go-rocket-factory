package main

import (
	"context"

	orderV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/openapi/order/v1"
)

type OrderHandler struct {
	orderService *orderService
}

func NewOrderHandler(os *orderService) *OrderHandler {
	return &OrderHandler{
		orderService: os,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	err := req.Validate()
	if err != nil {
		return &orderV1.BadRequestError{
			Error: "invalid request",
		}, nil
	}
	return h.orderService.CreateOrder(ctx, req)
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	err := req.Validate()
	if err != nil {
		return &orderV1.BadRequestError{
			Error: "invalid request",
		}, nil
	}
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Error: "invalid request",
		}, nil
	}
	return h.orderService.PayOrder(ctx, req, params.OrderUUID)
}

func (h *OrderHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	return h.orderService.GetOrderByUUID(ctx, params.OrderUUID)
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	return h.orderService.CancelOrderByUUID(ctx, params.OrderUUID)
}
