package order

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (*model.Order, error) {
	const query = `
		SELECT
			order_uuid,
			user_uuid,
			part_uuids,
			total_price_cents,
			transaction_uuid,
			payment_method,
			status
		FROM orders
		WHERE order_uuid = $1
	`

	var (
		order           model.Order
		transactionUUID sql.NullString
		paymentMethod   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, orderUUID).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&order.PartUUIDs,
		&order.TotalPriceCents,
		&transactionUUID,
		&paymentMethod,
		&order.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.Order{}, model.ErrOrderNotFound
		}

		return &model.Order{}, err
	}

	if transactionUUID.Valid {
		order.TransactionUUID = &transactionUUID.String
	}
	if paymentMethod.Valid {
		method := model.PaymentMethod(paymentMethod.String)
		order.PaymentMethod = &method
	}

	return &order, nil
}
