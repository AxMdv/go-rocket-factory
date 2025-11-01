package v1

import (
	"context"

	grpcConverter "github.com/AxMdv/go-rocket-factory/order/internal/client/converter"
	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *inventoryClient) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := c.ic.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{Uuids: uuids},
	})
	if err != nil {
		return []model.Part{}, err
	}
	return grpcConverter.PartsGRPCToModel(resp.GetParts()), nil
}
