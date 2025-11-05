package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		expectedOrder = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUUIDs:  []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice: gofakeit.Price(10, 100),
			Status:     model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&expectedOrder, nil).Once()

	order, err := s.service.Get(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(expectedOrder, order)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoErr   = gofakeit.Error()
		orderUUID = gofakeit.UUID()
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&model.Order{}, repoErr).Once()

	_, err := s.service.Get(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestGetNotFound() {
	orderUUID := gofakeit.UUID()

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

	_, err := s.service.Get(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}
