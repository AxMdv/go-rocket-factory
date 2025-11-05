package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/AxMdv/go-rocket-factory/order/internal/client/grpc/mocks"
	repoMocks "github.com/AxMdv/go-rocket-factory/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	orderRepository *repoMocks.OrderRepository

	inventoryClient *grpcMocks.InventoryClient
	paymentClient   *grpcMocks.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = repoMocks.NewOrderRepository(s.T())

	s.inventoryClient = grpcMocks.NewInventoryClient(s.T())
	s.paymentClient = grpcMocks.NewPaymentClient(s.T())

	s.service = NewOrderService(s.orderRepository, s.inventoryClient, s.paymentClient)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
