package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestPaySuccess() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		paymentMethod   = model.PaymentMethodCARD
		transactionUUID = gofakeit.UUID()

		order = model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).
		Return(transactionUUID, nil).Once()
	s.orderRepository.On("UpdateOrder", s.ctx, orderUUID, model.OrderUpdateInfo{
		TransactionUUID: lo.ToPtr(transactionUUID),
		PaymentMethod:   lo.ToPtr(paymentMethod),
		Status:          lo.ToPtr(model.OrderStatusPAID),
	}).Return(nil).Once()

	res, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayRepoGetError() {
	var (
		orderUUID     = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCARD
		repoErr       = gofakeit.Error()
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(nil, repoErr).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestPayInvalidStatus() {
	var (
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCARD

		order = model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPAID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&order, nil).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderStatusConflict)

	order = model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusCANCELLED,
	}
	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&order, nil).Once()

	_, err = s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderStatusConflict)
}

func (s *ServiceSuite) TestPayPaymentClientError() {
	var (
		clientErr     = gofakeit.Error()
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCREDITCARD

		order = model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).
		Return("", clientErr).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, clientErr)
	s.Require().ErrorContains(err, "payment error")
}

func (s *ServiceSuite) TestPayRepoUpdateError() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		paymentMethod   = model.PaymentMethodINVESTORMONEY
		transactionUUID = gofakeit.UUID()
		repoErr         = gofakeit.Error()

		order = model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).
		Return(transactionUUID, nil).Once()
	s.orderRepository.On("UpdateOrder", s.ctx, orderUUID, model.OrderUpdateInfo{
		TransactionUUID: lo.ToPtr(transactionUUID),
		PaymentMethod:   lo.ToPtr(paymentMethod),
		Status:          lo.ToPtr(model.OrderStatusPAID),
	}).Return(repoErr).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
