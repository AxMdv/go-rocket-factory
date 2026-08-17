package part

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
	repoModel "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/model"
)

func TestBuildPartsFilter(t *testing.T) {
	tests := []struct {
		name     string
		filter   *model.PartsFilter
		expected bson.M
	}{
		{
			name:     "nil filter",
			expected: bson.M{},
		},
		{
			name:     "empty filter",
			filter:   &model.PartsFilter{},
			expected: bson.M{},
		},
		{
			name: "all fields",
			filter: &model.PartsFilter{
				Uuids:                 []string{"uuid-1", "uuid-2"},
				Names:                 []string{"Engine", "Fuel Tank"},
				Categories:            []model.Category{model.CategoryEngine, model.CategoryFuel},
				ManufacturerCountries: []string{"Russia", "USA"},
				Tags:                  []string{"rocket", "core"},
			},
			expected: bson.M{
				"_id":                  bson.M{"$in": []string{"uuid-1", "uuid-2"}},
				"name":                 bson.M{"$in": []string{"Engine", "Fuel Tank"}},
				"category":             bson.M{"$in": []repoModel.Category{repoModel.CategoryEngine, repoModel.CategoryFuel}},
				"manufacturer.country": bson.M{"$in": []string{"Russia", "USA"}},
				"tags":                 bson.M{"$in": []string{"rocket", "core"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, buildPartsFilter(tt.filter))
		})
	}
}
