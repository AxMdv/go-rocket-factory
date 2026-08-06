package order

import (
	"context"
	"fmt"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return err
	}
	switch order.Status {
	case model.OrderStatusPAID:
		return model.ErrOrderStatusConflict
	case model.OrderStatusCANCELLED:
		return nil
	case model.OrderStatusPENDINGPAYMENT:
	default:
		return fmt.Errorf("unknown order status")
	}

	return s.orderRepository.Cancel(ctx, order.OrderUUID)
}
