package repository

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) error
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	MarkPaid(ctx context.Context, orderUUID, transactionUUID string, method model.PaymentMethod) error
	Cancel(ctx context.Context, orderUUID string) error
}
