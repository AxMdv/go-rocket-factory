package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	orderV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/payment/v1"
)

type orderService struct {
	mu              sync.RWMutex
	orders          map[string]*orderV1.OrderDto
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderService(ic inventoryV1.InventoryServiceClient, pc paymentV1.PaymentServiceClient) *orderService {
	return &orderService{
		orders:          make(map[string]*orderV1.OrderDto),
		inventoryClient: ic,
		paymentClient:   pc,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	resp, err := s.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{Uuids: req.GetPartUuids()},
	})
	if err != nil {
		return &orderV1.BadRequestError{Error: fmt.Sprintf("inventory error: %v", err)}, nil
	}
	if len(resp.GetParts()) != len(req.GetPartUuids()) {
		return &orderV1.NotFoundError{Error: "some parts not found"}, nil
	}

	var total float64
	for _, p := range resp.GetParts() {
		total += p.GetPrice()
	}

	order := &orderV1.OrderDto{
		OrderUUID:  uuid.New().String(),
		UserUUID:   req.GetUserUUID(),
		PartUuids:  req.GetPartUuids(),
		TotalPrice: total,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	s.mu.Lock()
	s.orders[order.OrderUUID] = order
	s.mu.Unlock()

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: total,
	}, nil
}

func (s *orderService) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, orderUUID string) (orderV1.PayOrderRes, error) {
	s.mu.RLock()
	order, ok := s.orders[orderUUID]
	s.mu.RUnlock()
	if !ok {
		return &orderV1.NotFoundError{Error: "order not found"}, nil
	}
	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{Error: "already paid"}, nil
	}
	if order.Status == orderV1.OrderStatusCANCELLED {
		return &orderV1.ConflictError{Error: "cannot pay cancelled order"}, nil
	}
	payMethod := req.GetPaymentMethod()

	var grpcPaymentMethod paymentV1.PaymentMethod
	switch payMethod {
	case orderV1.PaymentMethodCARD:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
	paymentReq := &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: grpcPaymentMethod,
	}
	payResp, err := s.paymentClient.PayOrder(ctx, paymentReq)
	if err != nil {
		return &orderV1.BadRequestError{Error: fmt.Sprintf("payment error: %v", err)}, nil
	}

	s.mu.Lock()
	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewOptNilString(payResp.GetTransactionUuid())
	order.PaymentMethod = orderV1.NewOptPaymentMethod(payMethod)
	s.mu.Unlock()

	return &orderV1.PayOrderResponse{TransactionUUID: payResp.GetTransactionUuid()}, nil
}

func (s *orderService) GetOrderByUUID(_ context.Context, uuid string) (orderV1.GetOrderByUUIDRes, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[uuid]
	if !ok {
		return &orderV1.NotFoundError{Error: "order not found"}, nil
	}
	return order, nil
}

func (s *orderService) CancelOrderByUUID(_ context.Context, uuid string) (orderV1.CancelOrderRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[uuid]
	if !ok {
		return &orderV1.NotFoundError{Error: "order not found"}, nil
	}

	switch order.Status {
	case orderV1.OrderStatusPAID:
		return &orderV1.ConflictError{Error: "cannot cancel paid order"}, nil
	case orderV1.OrderStatusCANCELLED:
		return &orderV1.CancelOrderNoContent{}, nil
	case orderV1.OrderStatusPENDINGPAYMENT:
		order.Status = orderV1.OrderStatusCANCELLED
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
