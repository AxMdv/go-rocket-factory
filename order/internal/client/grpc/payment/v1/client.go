package v1

import paymentV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/payment/v1"

type paymentClient struct {
	pc paymentV1.PaymentServiceClient
}

func NewPaymentClient(pc paymentV1.PaymentServiceClient) *paymentClient {
	return &paymentClient{
		pc: pc,
	}
}
