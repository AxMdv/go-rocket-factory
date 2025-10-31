package part

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
	repoConverter "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return *repoConverter.PartRepoToModel(&repoPart), nil
}
