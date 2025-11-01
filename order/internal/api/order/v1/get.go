package v1

import (
	"context"
	"errors"

	"github.com/AxMdv/go-rocket-factory/order/internal/converter"
	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	orderV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{Error: "order not found"}, nil
		}
		return &orderV1.BadRequestError{Error: err.Error()}, nil
	}
	resp := converter.OrderModelToDto(order)
	return &resp, nil
}
