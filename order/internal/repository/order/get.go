package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	repoConverter "github.com/AxMdv/go-rocket-factory/order/internal/repository/converter"
)

func (r *repository) GetOrderByUUID(_ context.Context, orderUUID string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	repoOrder, ok := r.orders[orderUUID]
	if !ok {
		return &model.Order{}, model.ErrOrderNotFound
	}
	return lo.ToPtr(repoConverter.RepoOrderToModel(repoOrder)), nil
}
