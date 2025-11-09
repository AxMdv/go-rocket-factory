package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		partUuid = gofakeit.UUID()

		expectedPart = model.Part{
			UUID:          partUuid,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Paragraph(1, 5, 1, " "),
			Price:         gofakeit.Price(0, 100),
			StockQuantity: int64(gofakeit.Number(0, 100)),
			Category:      model.Category(gofakeit.Number(1, 4)),
		}
	)

	s.partRepository.On("Get", s.ctx, partUuid).Return(expectedPart, nil).Once()

	part, err := s.service.Get(s.ctx, partUuid)
	s.Require().NoError(err)
	s.Require().Equal(expectedPart, part)
}

func (s *ServiceSuite) TestGetInvalidUUIDPart() {
	partUuid := gofakeit.UUID()

	s.partRepository.On("Get", s.ctx, partUuid).Return(model.Part{}, model.ErrInvalidUUID).Once()

	_, err := s.service.Get(s.ctx, partUuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrInvalidUUID)
}

func (s *ServiceSuite) TestGetNotFoundPart() {
	partUuid := gofakeit.UUID()

	s.partRepository.On("Get", s.ctx, partUuid).Return(model.Part{}, model.ErrPartNotFound).Once()

	_, err := s.service.Get(s.ctx, partUuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartNotFound)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoErr  = gofakeit.Error()
		partUuid = gofakeit.UUID()
	)

	s.partRepository.On("Get", s.ctx, partUuid).Return(model.Part{}, repoErr).Once()

	_, err := s.service.Get(s.ctx, partUuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
