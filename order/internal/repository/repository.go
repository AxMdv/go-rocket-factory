package repository

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) error
	GetOrderByUUID(ctx context.Context, orderUUID string) (*model.Order, error)
	UpdateOrder(ctx context.Context, orderUUID string, order model.OrderUpdateInfo) error
}
