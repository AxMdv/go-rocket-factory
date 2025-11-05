package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUuids = []string{partUUID1, partUUID2}

		part1 = model.Part{
			UUID:  partUUID1,
			Price: gofakeit.Price(0, 100),
		}
		part2 = model.Part{
			UUID:  partUUID2,
			Price: gofakeit.Price(0, 100),
		}
		parts = []model.Part{part1, part2}

		totalPrice = part1.Price + part2.Price

		createOrder = struct {
			UserUUID  string
			PartUUIDs []string
		}{
			UserUUID:  userUUID,
			PartUUIDs: partUuids,
		}

		expectedOrder = model.Order{
			OrderUUID:  gofakeit.UUID(),
			UserUUID:   userUUID,
			PartUUIDs:  partUuids,
			TotalPrice: totalPrice,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, partUuids).Return(parts, nil).Once()

	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		s.Require().NotEmpty(order.OrderUUID)
		s.Require().Equal(expectedOrder.UserUUID, order.UserUUID)
		s.Require().Equal(expectedOrder.PartUUIDs, order.PartUUIDs)
		s.Require().Equal(expectedOrder.TotalPrice, order.TotalPrice)
		s.Require().Equal(expectedOrder.Status, order.Status)

		return true
	})).Return(nil).Once()

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, createOrder.UserUUID, createOrder.PartUUIDs)
	s.Require().NoError(err)
	s.Require().NotEmpty(orderUUID)
	s.Require().Equal(totalPrice, expectedOrder.TotalPrice)
}

func (s *ServiceSuite) TestCreateInventoryServiceError() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUuids = []string{partUUID1, partUUID2}

		customErr = gofakeit.Error()
	)

	s.inventoryClient.On("ListParts", s.ctx, partUuids).Return([]model.Part{}, customErr).Once()

	_, _, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "inventory error")
	s.Require().ErrorContains(err, customErr.Error())
}

func (s *ServiceSuite) TestCreateSomePartsNotFound() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUuids = []string{partUUID1, partUUID2}

		part1 = model.Part{
			UUID:  partUUID1,
			Price: gofakeit.Price(0, 100),
		}

		parts = []model.Part{part1}
	)

	s.inventoryClient.On("ListParts", s.ctx, partUuids).Return(parts, nil).Once()

	_, _, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartsNotFound)
}

func (s *ServiceSuite) TestCreateRepoError() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUuids = []string{partUUID1, partUUID2}

		repoErr = gofakeit.Error()

		part1 = model.Part{
			UUID:  partUUID1,
			Price: gofakeit.Price(0, 100),
		}
		part2 = model.Part{
			UUID:  partUUID2,
			Price: gofakeit.Price(0, 100),
		}
		parts = []model.Part{part1, part2}

		totalPrice = part1.Price + part2.Price
	)

	s.inventoryClient.On("ListParts", s.ctx, partUuids).Return(parts, nil).Once()

	s.orderRepository.On("CreateOrder", s.ctx, mock.MatchedBy(func(order model.Order) bool {
		s.Require().NotEmpty(order.OrderUUID)
		s.Require().Equal(userUUID, order.UserUUID)
		s.Require().Equal(partUuids, order.PartUUIDs)
		s.Require().Equal(totalPrice, order.TotalPrice)
		s.Require().Equal(model.OrderStatusPENDINGPAYMENT, order.Status)

		return true
	})).Return(repoErr).Once()

	_, _, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
