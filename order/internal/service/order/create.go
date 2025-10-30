package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (orderUUID string, totalPrice float64, err error) {
	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return "", 0, fmt.Errorf("inventory error: %w", err)
	}
	if len(parts) != len(partUUIDs) {
		return "", 0, model.ErrPartsNotFound
	}

	for _, p := range parts {
		totalPrice += p.Price
	}

	order := model.Order{
		OrderUUID:  uuid.New().String(),
		UserUUID:   userUUID,
		PartUUIDs:  append([]string(nil), partUUIDs...),
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	if err := s.orderRepository.CreateOrder(ctx, order); err != nil {
		return "", 0, err
	}
	return order.OrderUUID, totalPrice, nil
}
