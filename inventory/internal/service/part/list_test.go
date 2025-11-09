package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestListSuccess() {
	var (
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUUID3 = gofakeit.UUID()
		name      = gofakeit.Name()
		category  = model.Category(gofakeit.Number(1, 4))

		filter = model.PartsFilter{
			Uuids:      []string{partUUID1, partUUID2, partUUID3},
			Names:      []string{name},
			Categories: []model.Category{category},
		}

		expectedPartList = []model.Part{
			{
				UUID:          partUUID1,
				Name:          name,
				Description:   gofakeit.Paragraph(1, 5, 1, " "),
				Price:         gofakeit.Price(0, 100),
				StockQuantity: int64(gofakeit.Number(0, 100)),
				Category:      category,
			},
			{
				UUID:          partUUID2,
				Name:          name,
				Description:   gofakeit.Paragraph(1, 5, 1, " "),
				Price:         gofakeit.Price(0, 100),
				StockQuantity: int64(gofakeit.Number(0, 100)),
				Category:      category,
			},
		}
	)

	s.partRepository.On("List", s.ctx, &filter).Return(expectedPartList, nil).Once()

	partList, err := s.service.List(s.ctx, &filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedPartList, partList)
}

func (s *ServiceSuite) TestListRepoError() {
	var (
		repoErr   = gofakeit.Error()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()
		partUUID3 = gofakeit.UUID()
		name      = gofakeit.Name()
		category  = model.Category(gofakeit.Number(1, 4))

		filter = model.PartsFilter{
			Uuids:      []string{partUUID1, partUUID2, partUUID3},
			Names:      []string{name},
			Categories: []model.Category{category},
		}
	)

	s.partRepository.On("List", s.ctx, &filter).Return(nil, repoErr).Once()

	_, err := s.service.List(s.ctx, &filter)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
