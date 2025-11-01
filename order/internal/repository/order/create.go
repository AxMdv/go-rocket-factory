package order

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	repoConverter "github.com/AxMdv/go-rocket-factory/order/internal/repository/converter"
)

func (r *repository) CreateOrder(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.OrderUUID] = repoConverter.ModelOrderToRepo(order)
	return nil
}
