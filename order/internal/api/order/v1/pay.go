package v1

import (
	"context"
	"errors"

	"github.com/AxMdv/go-rocket-factory/order/internal/converter"
	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	orderV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	err := req.Validate()
	if err != nil || params.OrderUUID == "" {
		return &orderV1.BadRequestError{Error: "invalid request"}, nil
	}

	transactionUUID, err := a.orderService.PayOrder(ctx, params.OrderUUID, converter.PaymentMethodDtoToModel(req.GetPaymentMethod()))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{Error: "order not found"}, nil
		}
		if errors.Is(err, model.ErrOrderStatusConflict) {
			return &orderV1.ConflictError{Error: "order already paid or cancelled"}, nil
		}
		return &orderV1.InternalServerError{Error: err.Error()}, nil
	}
	return &orderV1.PayOrderResponse{TransactionUUID: transactionUUID}, nil
}
