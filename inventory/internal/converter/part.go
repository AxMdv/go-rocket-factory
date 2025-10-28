package converter

import (
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/AxMdv/go-rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
)

func CategoryToModel(c inventoryV1.Category) model.Category {
	switch c {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func CategoryToProto(c model.Category) inventoryV1.Category {
	switch c {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNKNOWN
	}
}

func DimensionsToModel(d *inventoryV1.Dimensions) *model.Dimensions {
	if d == nil {
		return nil
	}
	return &model.Dimensions{
		Length: d.GetLength(),
		Width:  d.GetWidth(),
		Height: d.GetHeight(),
		Weight: d.GetWeight(),
	}
}

func DimensionsToProto(d *model.Dimensions) *inventoryV1.Dimensions {
	if d == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ManufacturerToModel(m *inventoryV1.Manufacturer) *model.Manufacturer {
	if m == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    m.GetName(),
		Country: m.GetCountry(),
		Website: m.GetWebsite(),
	}
}

func ManufacturerToProto(m *model.Manufacturer) *inventoryV1.Manufacturer {
	if m == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func ValueToModel(v *inventoryV1.Value) model.Value {
	var mv model.Value
	switch k := v.GetKind().(type) {
	case *inventoryV1.Value_StringValue:
		mv.String = lo.ToPtr(k.StringValue)
	case *inventoryV1.Value_Int64Value:
		x := k.Int64Value
		mv.Int64 = &x
	case *inventoryV1.Value_DoubleValue:
		x := k.DoubleValue
		mv.Double = &x
	case *inventoryV1.Value_BoolValue:
		x := k.BoolValue
		mv.Bool = &x
	}
	return mv
}

func ValueToProto(v model.Value) *inventoryV1.Value {
	out := &inventoryV1.Value{}
	switch {
	case v.String != nil:
		out.Kind = &inventoryV1.Value_StringValue{StringValue: *v.String}
	case v.Int64 != nil:
		out.Kind = &inventoryV1.Value_Int64Value{Int64Value: *v.Int64}
	case v.Double != nil:
		out.Kind = &inventoryV1.Value_DoubleValue{DoubleValue: *v.Double}
	case v.Bool != nil:
		out.Kind = &inventoryV1.Value_BoolValue{BoolValue: *v.Bool}
	default:
		// пустое значение — оставляем Kind nil
	}
	return out
}

func timeToProto(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func timeToModel(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func PartToModel(p *inventoryV1.Part) *model.Part {
	if p == nil {
		return nil
	}

	md := make(map[string]model.Value, len(p.GetMetadata()))
	for k, v := range p.GetMetadata() {
		if v == nil {
			continue
		}
		md[k] = ValueToModel(v)
	}

	return &model.Part{
		UUID:          p.GetUuid(),
		Name:          p.GetName(),
		Description:   p.GetDescription(),
		Price:         p.GetPrice(),
		StockQuantity: p.GetStockQuantity(),
		Category:      CategoryToModel(p.GetCategory()),
		Dimensions:    DimensionsToModel(p.GetDimensions()),
		Manufacturer:  ManufacturerToModel(p.GetManufacturer()),
		Tags:          append([]string(nil), p.GetTags()...),
		Metadata:      md,
		CreatedAt:     timeToModel(p.GetCreatedAt()),
		UpdatedAt:     timeToModel(p.GetUpdatedAt()),
	}
}

func PartToProto(m *model.Part) *inventoryV1.Part {
	if m == nil {
		return nil
	}

	md := make(map[string]*inventoryV1.Value, len(m.Metadata))
	for k, v := range m.Metadata {
		md[k] = ValueToProto(v)
	}

	return &inventoryV1.Part{
		Uuid:          m.UUID,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		StockQuantity: m.StockQuantity,
		Category:      CategoryToProto(m.Category),
		Dimensions:    DimensionsToProto(m.Dimensions),
		Manufacturer:  ManufacturerToProto(m.Manufacturer),
		Tags:          append([]string(nil), m.Tags...),
		Metadata:      md,
		CreatedAt:     timeToProto(m.CreatedAt),
		UpdatedAt:     timeToProto(m.UpdatedAt),
	}
}

// PartsFilterToModel конвертирует Proto фильтр в доменную модель.
func PartsFilterToModel(f *inventoryV1.PartsFilter) model.PartsFilter {
	if f == nil {
		return model.PartsFilter{}
	}

	var cats []model.Category
	if len(f.GetCategories()) > 0 {
		cats = make([]model.Category, 0, len(f.GetCategories()))
		for _, c := range f.GetCategories() {
			cats = append(cats, CategoryToModel(c))
		}
	}

	return model.PartsFilter{
		Uuids:                 append([]string(nil), f.GetUuids()...),
		Names:                 append([]string(nil), f.GetNames()...),
		Categories:            cats,
		ManufacturerCountries: append([]string(nil), f.GetManufacturerCountries()...),
		Tags:                  append([]string(nil), f.GetTags()...),
	}
}

func PartsToProto(p []model.Part) []*inventoryV1.Part {
	pbParts := make([]*inventoryV1.Part, 0, len(p))
	for _, part := range p {
		pbParts = append(pbParts, PartToProto(&part))
	}
	return pbParts
}
