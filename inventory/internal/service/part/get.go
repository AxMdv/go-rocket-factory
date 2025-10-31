package part

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.partRepository.Get(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
