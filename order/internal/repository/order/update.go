package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	repoConverter "github.com/AxMdv/go-rocket-factory/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, orderUUID string, updateInfo model.OrderUpdateInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[orderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	// Обновляем поля, только если они были установлены в запросе
	if updateInfo.UserUUID != nil {
		order.UserUUID = *updateInfo.UserUUID
	}
	if updateInfo.PartUUIDs != nil {
		order.PartUUIDs = *updateInfo.PartUUIDs
	}

	if updateInfo.TotalPrice != nil {
		order.TotalPrice = *updateInfo.TotalPrice
	}
	if updateInfo.TransactionUUID != nil {
		order.TransactionUUID = updateInfo.TransactionUUID
	}
	if updateInfo.PaymentMethod != nil {
		order.PaymentMethod = lo.ToPtr(repoConverter.PaymentMethodToRepo(*updateInfo.PaymentMethod))
	}
	if updateInfo.Status != nil {
		order.Status = repoConverter.OrderStatusToRepo(*updateInfo.Status)
	}

	r.orders[orderUUID] = order

	return nil
}
