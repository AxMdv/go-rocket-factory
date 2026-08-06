package order

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order) error {
	const query = `
		INSERT INTO orders (
			order_uuid,
			user_uuid,
			part_uuids,
			total_price_cents,
			status
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		order.OrderUUID,
		order.UserUUID,
		order.PartUUIDs,
		order.TotalPriceCents,
		order.Status,
	)
	if err != nil {
		return err
	}

	return nil
}
