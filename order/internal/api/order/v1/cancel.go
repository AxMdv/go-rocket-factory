package v1

import (
	"context"
	"errors"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	orderV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{Error: "order not found"}, nil
		}
		if errors.Is(err, model.ErrConflict) {
			return &orderV1.ConflictError{Error: "cannot cancel paid order"}, nil
		}
		return &orderV1.InternalServerError{Error: err.Error()}, nil
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
