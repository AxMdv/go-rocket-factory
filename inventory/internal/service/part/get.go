package part

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
)

func (s *service) GetPartByUUID(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.partRepository.GetPartByUUID(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
