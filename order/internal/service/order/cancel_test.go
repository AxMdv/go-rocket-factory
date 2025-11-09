package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		order = model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&order, nil).Once()
	s.orderRepository.On("UpdateOrder", s.ctx, orderUUID, model.OrderUpdateInfo{
		Status: lo.ToPtr(model.OrderStatusCANCELLED),
	}).Return(nil).Once()

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelSuccessWhenCancelled() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		order = model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusCANCELLED,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&order, nil).Once()

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelRepoGetError() {
	var (
		repoErr   = gofakeit.Error()
		orderUUID = gofakeit.UUID()
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(nil, repoErr)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCancelInvalidStatus() {
	var (
		orderUUID = gofakeit.UUID()

		order = model.Order{
			OrderUUID: orderUUID,
			UserUUID:  gofakeit.UUID(),
			Status:    model.OrderStatusPAID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&order, nil)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderStatusConflict)
}

func (s *ServiceSuite) TestCancelRepoUpdateError() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		repoErr   = gofakeit.Error()

		order = model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&order, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, orderUUID, model.OrderUpdateInfo{
		Status: lo.ToPtr(model.OrderStatusCANCELLED),
	}).Return(repoErr)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
