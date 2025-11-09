package payment

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/AxMdv/go-rocket-factory/payment/internal/model"
)

func (s *ServiceSuite) TestPayUserOrderSuccess() {
	userUUID := gofakeit.UUID()
	orderUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethod(gofakeit.Number(1, 4))
	transactionId, err := s.service.PayUserOrder(s.ctx, userUUID, orderUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().NotEmpty(transactionId)
}
