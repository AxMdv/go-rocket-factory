package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/AxMdv/go-rocket-factory/payment/internal/model"
)

func (s *service) PayUserOrder(ctx context.Context, userUUID, orderUUID string, paymentMethod model.PaymentMethod) (transactionUUID string, err error) {
	transactionUUID = uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)
	return transactionUUID, nil
}
