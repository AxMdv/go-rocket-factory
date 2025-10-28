package v1

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/converter"
	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.ListPartsByFilter(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error %s", err)
	}
	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
