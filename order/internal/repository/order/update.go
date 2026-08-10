package order

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (r *repository) MarkPaid(ctx context.Context, orderUUID, transactionUUID string, method model.PaymentMethod) error {
	const query = `
		UPDATE orders
		SET
			status = $2,
			transaction_uuid = $3,
			payment_method = $4,
			updated_at = now()
		WHERE order_uuid = $1
			AND status = $5
	`

	res, err := r.db.Exec(
		ctx,
		query,
		orderUUID,
		model.OrderStatusPAID,
		transactionUUID,
		method,
		model.OrderStatusPENDINGPAYMENT,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return model.ErrOrderStatusConflict
	}

	return nil
}

func (r *repository) Cancel(ctx context.Context, orderUUID string) error {
	const query = `
		UPDATE orders
		SET
			status = $2,
			updated_at = now()
		WHERE order_uuid = $1
			AND status = $3
	`

	res, err := r.db.Exec(ctx, query, orderUUID, model.OrderStatusCANCELLED, model.OrderStatusPENDINGPAYMENT)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return model.ErrOrderStatusConflict
	}

	return nil
}
