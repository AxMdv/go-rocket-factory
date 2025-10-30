package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	orderV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	err := req.Validate()
	if err != nil {
		return &orderV1.BadRequestError{
			Error: "invalid request",
		}, nil
	}
	orderUUID, totalPrice, err := a.orderService.CreateOrder(ctx, req.GetUserUUID(), req.GetPartUuids())
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) {
			return &orderV1.NotFoundError{Error: "some parts not found"}, nil
		}
		return &orderV1.InternalServerError{Error: fmt.Sprintf("internal err %v", err)}, nil
	}
	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
