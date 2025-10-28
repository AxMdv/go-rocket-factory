package service

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
)

type PartService interface {
	GetPartByUUID(ctx context.Context, uuid string) (model.Part, error)
	ListPartsByFilter(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
