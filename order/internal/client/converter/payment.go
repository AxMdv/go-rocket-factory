package converter

import (
	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	paymentV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/payment/v1"
)

func CreatePayOrderRequest(orderUUID, userUUID string, method model.PaymentMethod) *paymentV1.PayOrderRequest {
	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentMethodModelToGRPC(method),
	}
	return req
}

func paymentMethodModelToGRPC(payMethod model.PaymentMethod) paymentV1.PaymentMethod {
	var grpcPaymentMethod paymentV1.PaymentMethod
	switch payMethod {
	case model.PaymentMethodCARD:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCREDITCARD:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodINVESTORMONEY:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		grpcPaymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
	return grpcPaymentMethod
}
