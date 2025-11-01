package part

import (
	"context"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
	repoConverter "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	filteredParts := filterParts(r.parts, repoConverter.PartsFilterModelToRepo(filter))

	return repoConverter.PartsRepoToModel(filteredParts), nil
}

// filterByUUID (ИЛИ).
// Если uuids пуст — вернуть все элементы из allParts.
func filterByUUID(allParts map[string]repoModel.Part, uuids []string) []repoModel.Part {
	if len(uuids) == 0 {
		out := make([]repoModel.Part, 0, len(allParts))
		for i := range allParts {
			// берём адрес значения через временную переменную
			p := allParts[i]
			out = append(out, p)
		}
		return out
	}
	out := make([]repoModel.Part, 0, len(uuids))
	for _, id := range uuids {
		if val, ok := allParts[id]; ok {
			p := val
			out = append(out, p)
		}
	}
	return out
}

// filterByNames (ИЛИ) поверх слайса.
func filterByNames(parts []repoModel.Part, names []string) []repoModel.Part {
	if len(names) == 0 {
		return parts
	}
	nameSet := make(map[string]struct{}, len(names))
	for _, n := range names {
		nameSet[n] = struct{}{}
	}
	out := make([]repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := nameSet[p.Name]; ok {
			out = append(out, p)
		}
	}
	return out
}

// filterByCategories (ИЛИ) поверх слайса.
func filterByCategories(parts []repoModel.Part, categories []repoModel.Category) []repoModel.Part {
	if len(categories) == 0 {
		return parts
	}
	catSet := make(map[repoModel.Category]struct{}, len(categories))
	for _, c := range categories {
		catSet[c] = struct{}{}
	}
	out := make([]repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if _, ok := catSet[p.Category]; ok {
			out = append(out, p)
		}
	}
	return out
}

// filterByManufacturerCountries (ИЛИ) поверх слайса.
func filterByManufacturerCountries(parts []repoModel.Part, countries []string) []repoModel.Part {
	if len(countries) == 0 {
		return parts
	}
	countrySet := make(map[string]struct{}, len(countries))
	for _, c := range countries {
		countrySet[c] = struct{}{}
	}
	out := make([]repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if p.Manufacturer == nil {
			continue
		}
		if _, ok := countrySet[p.Manufacturer.Country]; ok {
			out = append(out, p)
		}
	}
	return out
}

// filterByTags (ИЛИ) поверх слайса.
func filterByTags(parts []repoModel.Part, tags []string) []repoModel.Part {
	if len(tags) == 0 {
		return parts
	}
	tagSet := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagSet[t] = struct{}{}
	}
	out := make([]repoModel.Part, 0, len(parts))
Outer:
	for _, p := range parts {
		for _, tag := range p.Tags {
			if _, ok := tagSet[tag]; ok {
				out = append(out, p)
				continue Outer
			}
		}
	}
	return out
}

// FilterParts — основная функция (логическое И между полями).
// Вход: map[string]repoModel.Part, т.е. значения, а не указатели.
func filterParts(allParts map[string]repoModel.Part, filter *repoModel.PartsFilter) []repoModel.Part {
	if filter == nil {
		// вернуть все
		out := make([]repoModel.Part, 0, len(allParts))
		for k := range allParts {
			p := allParts[k]
			out = append(out, p)
		}
		return out
	}
	parts := filterByUUID(allParts, filter.Uuids)
	parts = filterByNames(parts, filter.Names)
	parts = filterByCategories(parts, filter.Categories)
	parts = filterByManufacturerCountries(parts, filter.ManufacturerCountries)
	parts = filterByTags(parts, filter.Tags)
	return parts
}
