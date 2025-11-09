package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/converter"
	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
)

// GetPart возвращает информацию о детали по её идентификатору.
func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}
		if errors.Is(err, model.ErrInvalidUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid part UUID %s", req.GetUuid())
		}
		return nil, status.Errorf(codes.Internal, "internal error %s", err)

	}
	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
