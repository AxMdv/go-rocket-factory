// nolint:gosec
package part

import (
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	repoModel "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/model"
)

func createRepoParts(count int) map[string]repoModel.Part {
	err := gofakeit.Seed(time.Now().UnixNano())
	if err != nil {
		log.Fatal(err)
	}

	parts := make(map[string]repoModel.Part, count)

	categories := []repoModel.Category{
		repoModel.CategoryEngine,
		repoModel.CategoryFuel,
		repoModel.CategoryPorthole,
		repoModel.CategoryWing,
	}

	genDimensions := func() *repoModel.Dimensions {
		return &repoModel.Dimensions{
			Length: gofakeit.Float64Range(0.1, 50.0),
			Width:  gofakeit.Float64Range(0.1, 10.0),
			Height: gofakeit.Float64Range(0.1, 10.0),
			Weight: gofakeit.Float64Range(0.5, 5000.0),
		}
	}

	genManufacturer := func() *repoModel.Manufacturer {
		return &repoModel.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		}
	}

	genTags := func() []string {
		base := []string{"space", "rocket", "module", "spare", "core", "nano", "mkII", "mkIII"}
		n := gofakeit.Number(1, 4)
		out := make([]string, 0, n)
		for i := 0; i < n; i++ {
			out = append(out, base[rand.Intn(len(base))])
		}
		return out
	}

	genMetadata := func() map[string]interface{} {
		batch := gofakeit.UUID()
		lifetime := gofakeit.Float64Range(100, 100000)
		refurb := gofakeit.Bool()

		return map[string]interface{}{
			"batch":          batch,
			"lifetime_hours": lifetime,
			"refurbished":    refurb,
		}
	}

	now := time.Now()

	for i := 0; i < count; i++ {
		id := uuid.New().String()
		name := gofakeit.Company() + " Part"
		desc := gofakeit.Sentence(8)

		cat := categories[gofakeit.Number(0, len(categories)-1)]
		price := gofakeit.Price(100.0, 500000.0)
		stock := int64(gofakeit.Number(0, 500))

		part := repoModel.Part{
			UUID:          id,
			Name:          name,
			Description:   desc,
			Price:         price,
			StockQuantity: stock,
			Category:      cat,
			Dimensions:    genDimensions(),
			Manufacturer:  genManufacturer(),
			Tags:          genTags(),
			Metadata:      genMetadata(),
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		parts[id] = part
	}

	return parts
}
