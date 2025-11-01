package converter

import (
	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartsGRPCToModel(grpcParts []*inventoryV1.Part) []model.Part {
	modelParts := make([]model.Part, 0, len(grpcParts))
	for _, part := range grpcParts {
		modelPart := model.Part{
			UUID:  part.GetUuid(),
			Price: part.GetPrice(),
		}
		modelParts = append(modelParts, modelPart)
	}
	return modelParts
}
