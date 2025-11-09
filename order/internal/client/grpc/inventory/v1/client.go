package v1

import inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"

type inventoryClient struct {
	ic inventoryV1.InventoryServiceClient
}

func NewInventoryClient(ic inventoryV1.InventoryServiceClient) *inventoryClient {
	return &inventoryClient{
		ic: ic,
	}
}
