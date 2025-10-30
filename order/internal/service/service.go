package service

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (orderUUID string, total float64, err error)
	PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (txUUID string, err error)
	GetOrderByUUID(ctx context.Context, orderUUID string) (model.Order, error)
	CancelOrderByUUID(ctx context.Context, orderUUID string) error
}
