package part

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
	repoConverter "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	mongoFilter := buildPartsFilter(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("find parts: %w", err)
	}

	var repoParts []repoModel.Part
	if err = cursor.All(ctx, &repoParts); err != nil {
		return nil, fmt.Errorf("decode parts cursor: %w", err)
	}

	return repoConverter.PartsRepoToModel(repoParts), nil
}

func buildPartsFilter(filter *model.PartsFilter) bson.M {
	mongoFilter := bson.M{}
	if filter == nil {
		return mongoFilter
	}

	if len(filter.Uuids) > 0 {
		mongoFilter["_id"] = bson.M{"$in": filter.Uuids}
	}
	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": filter.Names}
	}
	if len(filter.Categories) > 0 {
		categories := make([]repoModel.Category, 0, len(filter.Categories))
		for _, category := range filter.Categories {
			categories = append(categories, repoConverter.CategoryModelToRepo(category))
		}

		mongoFilter["category"] = bson.M{"$in": categories}
	}
	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}
	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	return mongoFilter
}
