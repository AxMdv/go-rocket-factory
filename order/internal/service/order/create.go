package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (orderUUID string, totalPriceCents int64, err error) {
	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return "", 0, fmt.Errorf("inventory error: %w", err)
	}
	if len(parts) != len(partUUIDs) {
		return "", 0, model.ErrPartsNotFound
	}

	for _, p := range parts {
		totalPriceCents += p.PriceCents
	}

	order := model.Order{
		OrderUUID:       uuid.New().String(),
		UserUUID:        userUUID,
		PartUUIDs:       partUUIDs,
		TotalPriceCents: totalPriceCents,
		Status:          model.OrderStatusPENDINGPAYMENT,
	}

	if err := s.orderRepository.CreateOrder(ctx, order); err != nil {
		return "", 0, err
	}
	return order.OrderUUID, totalPriceCents, nil
}
