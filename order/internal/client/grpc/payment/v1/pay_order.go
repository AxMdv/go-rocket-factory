package v1

import (
	"context"

	grpcConverter "github.com/AxMdv/go-rocket-factory/order/internal/client/converter"
	"github.com/AxMdv/go-rocket-factory/order/internal/model"
)

func (c *paymentClient) PayOrder(ctx context.Context, orderUUID, userUUID string,
	method model.PaymentMethod,
) (transactionUUID string, err error) {
	payResp, err := c.pc.PayOrder(ctx, grpcConverter.CreatePayOrderRequest(orderUUID, userUUID, method))
	if err != nil {
		return "", err
	}
	return payResp.GetTransactionUuid(), nil
}
