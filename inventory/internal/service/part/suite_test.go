package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repoMocks "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	partRepository *repoMocks.PartRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.partRepository = repoMocks.NewPartRepository(s.T())

	s.service = NewService(s.partRepository)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
