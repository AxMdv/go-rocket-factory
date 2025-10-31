package converter

import (
	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
	repoModel "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/model"
)

func categoryModelToRepo(c model.Category) repoModel.Category {
	switch c {
	case model.CategoryEngine:
		return repoModel.CategoryEngine
	case model.CategoryFuel:
		return repoModel.CategoryFuel
	case model.CategoryPorthole:
		return repoModel.CategoryPorthole
	case model.CategoryWing:
		return repoModel.CategoryWing
	default:
		return repoModel.CategoryUnknown
	}
}

func categoryRepoToModel(c repoModel.Category) model.Category {
	switch c {
	case repoModel.CategoryEngine:
		return model.CategoryEngine
	case repoModel.CategoryFuel:
		return model.CategoryFuel
	case repoModel.CategoryPorthole:
		return model.CategoryPorthole
	case repoModel.CategoryWing:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func dimensionsModelToRepo(d *model.Dimensions) *repoModel.Dimensions {
	if d == nil {
		return nil
	}
	return &repoModel.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func dimensionsRepoToModel(d *repoModel.Dimensions) *model.Dimensions {
	if d == nil {
		return nil
	}
	return &model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerModelToRepo(m *model.Manufacturer) *repoModel.Manufacturer {
	if m == nil {
		return nil
	}
	return &repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func manufacturerRepoToModel(m *repoModel.Manufacturer) *model.Manufacturer {
	if m == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func PartModelToRepo(p *model.Part) *repoModel.Part {
	if p == nil {
		return nil
	}
	return &repoModel.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      categoryModelToRepo(p.Category),
		Dimensions:    dimensionsModelToRepo(p.Dimensions),
		Manufacturer:  manufacturerModelToRepo(p.Manufacturer),
		Tags:          append([]string(nil), p.Tags...),
		Metadata:      p.Metadata,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func PartRepoToModel(p *repoModel.Part) *model.Part {
	if p == nil {
		return nil
	}
	return &model.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      categoryRepoToModel(p.Category),
		Dimensions:    dimensionsRepoToModel(p.Dimensions),
		Manufacturer:  manufacturerRepoToModel(p.Manufacturer),
		Tags:          append([]string(nil), p.Tags...),
		Metadata:      p.Metadata,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func PartsModelToRepo(in []model.Part) []*repoModel.Part {
	if len(in) == 0 {
		return nil
	}
	out := make([]*repoModel.Part, 0, len(in))
	for i := range in {
		out = append(out, PartModelToRepo(&in[i]))
	}
	return out
}

func PartsRepoToModel(in []*repoModel.Part) []model.Part {
	if len(in) == 0 {
		return nil
	}
	out := make([]model.Part, 0, len(in))
	for _, p := range in {
		if mp := PartRepoToModel(p); mp != nil {
			out = append(out, *mp)
		}
	}
	return out
}

func PartsFilterModelToRepo(f *model.PartsFilter) *repoModel.PartsFilter {
	if f == nil {
		return nil
	}
	cats := make([]repoModel.Category, 0, len(f.Categories))
	for _, c := range f.Categories {
		cats = append(cats, categoryModelToRepo(c))
	}
	return &repoModel.PartsFilter{
		Uuids:                 append([]string(nil), f.Uuids...),
		Names:                 append([]string(nil), f.Names...),
		Categories:            cats,
		ManufacturerCountries: append([]string(nil), f.ManufacturerCountries...),
		Tags:                  append([]string(nil), f.Tags...),
	}
}

func PartsFilterRepoToModel(f *repoModel.PartsFilter) *model.PartsFilter {
	if f == nil {
		return nil
	}
	cats := make([]model.Category, 0, len(f.Categories))
	for _, c := range f.Categories {
		cats = append(cats, categoryRepoToModel(c))
	}
	return &model.PartsFilter{
		Uuids:                 append([]string(nil), f.Uuids...),
		Names:                 append([]string(nil), f.Names...),
		Categories:            cats,
		ManufacturerCountries: append([]string(nil), f.ManufacturerCountries...),
		Tags:                  append([]string(nil), f.Tags...),
	}
}
