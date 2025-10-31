package order

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return err
	}
	switch order.Status {
	case model.OrderStatusPAID:
		return model.ErrConflict
	case model.OrderStatusCANCELLED:
		return nil
	case model.OrderStatusPENDINGPAYMENT:
	default:
		return fmt.Errorf("unknown order status")
	}
	updateInfo := model.OrderUpdateInfo{
		Status: lo.ToPtr(model.OrderStatusCANCELLED),
	}
	return s.orderRepository.UpdateOrder(ctx, order.OrderUUID, updateInfo)
}
