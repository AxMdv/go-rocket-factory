package v1

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/payment/internal/converter"
	paymentV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	orderUUID := req.GetOrderUuid()
	userUUID := req.GetUserUuid()
	if orderUUID == "" || userUUID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid or user_uuid not specified")
	}

	// Проверка  payment_method
	pm := req.GetPaymentMethod()
	if pm == paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "payment_method must be specified and not UNKNOWN")
	}
	switch pm {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported payment_method: %v", pm)
	}
	transactionUUID, err := a.paymentService.PayUserOrder(ctx, userUUID, orderUUID, converter.PaymentMethodToModel(pm))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error %s", err)
	}
	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
