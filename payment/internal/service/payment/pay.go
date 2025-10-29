package payment

import (
	"context"
	"log"

	"github.com/AxMdv/go-rocket-factory/payment/internal/model"
	"github.com/google/uuid"
)

func (s *service) PayUserOrder(ctx context.Context, userUUID string, orderUUID string, paymentMethod model.PaymentMethod) (transactionUUID string, err error) {
	transactionUUID = uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)
	return transactionUUID, nil
}
