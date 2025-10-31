package order

import (
	"context"
	"fmt"

	"github.com/samber/lo"

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

	updateInfo := model.OrderUpdateInfo{
		Status:          lo.ToPtr(model.OrderStatusPAID),
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &method,
	}
	if err := s.orderRepository.UpdateOrder(ctx, order.OrderUUID, updateInfo); err != nil {
		return "", err
	}
	return transactionUUID, nil
}
