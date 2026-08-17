package part

import (
	"context"
	"fmt"

	repoModel "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/model"
)

func (r *repository) insertParts(ctx context.Context, parts []repoModel.Part) error {
	if len(parts) == 0 {
		return nil
	}

	documents := make([]interface{}, len(parts))
	for i := range parts {
		documents[i] = parts[i]
	}

	if _, err := r.collection.InsertMany(ctx, documents); err != nil {
		return fmt.Errorf(
			"insert %d parts into collection %q: %w",
			len(parts),
			r.collection.Name(),
			err,
		)
	}

	return nil
}
