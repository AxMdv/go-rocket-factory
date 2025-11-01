package order

import (
	gprcClient "github.com/AxMdv/go-rocket-factory/order/internal/client/grpc"
	"github.com/AxMdv/go-rocket-factory/order/internal/repository"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient gprcClient.InventoryClient
	paymentClient   gprcClient.PaymentClient
}

func NewOrderService(repo repository.OrderRepository, inv gprcClient.InventoryClient, pay gprcClient.PaymentClient) *service {
	return &service{
		orderRepository: repo,
		inventoryClient: inv,
		paymentClient:   pay,
	}
}
