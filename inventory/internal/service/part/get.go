package part

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, partUUID string) (model.Part, error) {
	err := uuid.Validate(partUUID)
	if err != nil {
		return model.Part{}, model.ErrInvalidUUID
	}
	part, err := s.partRepository.Get(ctx, partUUID)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
