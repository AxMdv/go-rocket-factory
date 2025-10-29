package service

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/payment/internal/model"
)

type PaymentService interface {
	PayUserOrder(ctx context.Context, userUUID string, orderUUID string, paymentMethod model.PaymentMethod) (transactionUUID string, err error)
}
