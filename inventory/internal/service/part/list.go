package part

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
)

func (s *service) ListPartsByFilter(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	part, err := s.partRepository.ListPartsByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return part, nil
}
