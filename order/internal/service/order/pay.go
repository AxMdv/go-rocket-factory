package order

import (
	"context"
	"fmt"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return "", err
	}
	if order.Status != model.OrderStatusPENDINGPAYMENT {
		return "", model.ErrOrderStatusConflict
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.OrderUUID, order.UserUUID, method)
	if err != nil {
		return "", fmt.Errorf("payment error: %w", err)
	}

	if err := s.orderRepository.MarkPaid(ctx, order.OrderUUID, transactionUUID, method); err != nil {
		return "", err
	}
	return transactionUUID, nil
}
