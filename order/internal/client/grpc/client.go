package grpc

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (transactionUUID string, err error)
}
