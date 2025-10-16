package main

import (
	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
)

// Фильтр по UUID (ИЛИ)
func filterByUUID(parts map[string]*inventoryV1.Part, uuids []string) []*inventoryV1.Part {
	if len(uuids) == 0 {
		resultParts := make([]*inventoryV1.Part, 0, len(parts))
		for _, part := range parts {
			resultParts = append(resultParts, part)
		}
		return resultParts
	}
	filteredParts := make([]*inventoryV1.Part, 0)
	for _, uuid := range uuids {
		part, ok := parts[uuid]
		if ok {
			filteredParts = append(filteredParts, part)
		}
	}
	return filteredParts
}

// Фильтр по именам (ИЛИ)
func filterByNames(parts []*inventoryV1.Part, names []string) []*inventoryV1.Part {
	if len(names) == 0 {
		return parts
	}
	nameSet := make(map[string]struct{}, len(names)) // map[nameFromFilter]struct{}{}
	for _, n := range names {
		nameSet[n] = struct{}{}
	}
	var filtered []*inventoryV1.Part
	for _, p := range parts {
		if _, ok := nameSet[p.Name]; ok {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

// Фильтр по категориям (ИЛИ)
func filterByCategories(parts []*inventoryV1.Part, categories []inventoryV1.Category) []*inventoryV1.Part {
	if len(categories) == 0 {
		return parts
	}
	catSet := make(map[inventoryV1.Category]struct{}, len(categories))
	for _, c := range categories {
		catSet[c] = struct{}{}
	}
	var filtered []*inventoryV1.Part
	for _, p := range parts {
		if _, ok := catSet[p.Category]; ok {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

// Фильтр по странам (ИЛИ)
func filterByManufacturerCountries(parts []*inventoryV1.Part, countries []string) []*inventoryV1.Part {
	if len(countries) == 0 {
		return parts
	}
	countrySet := make(map[string]struct{}, len(countries))
	for _, c := range countries {
		countrySet[c] = struct{}{}
	}
	var filtered []*inventoryV1.Part
	for _, p := range parts {
		if _, ok := countrySet[p.Manufacturer.Country]; ok {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

// Фильтр по тегам (ИЛИ)
func filterByTags(parts []*inventoryV1.Part, tags []string) []*inventoryV1.Part {
	if len(tags) == 0 {
		return parts
	}
	tagSet := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagSet[t] = struct{}{}
	}
	var filtered []*inventoryV1.Part
Outer:
	for _, p := range parts {
		for _, tag := range p.Tags {
			if _, ok := tagSet[tag]; ok {
				filtered = append(filtered, p)
				continue Outer
			}
		}
	}
	return filtered
}

// Основная функция фильтрации (логическое И между полями)
func filterParts(allParts map[string]*inventoryV1.Part, filter *inventoryV1.PartsFilter) []*inventoryV1.Part {

	parts := filterByUUID(allParts, filter.Uuids)
	parts = filterByNames(parts, filter.Names)
	parts = filterByCategories(parts, filter.Categories)
	parts = filterByManufacturerCountries(parts, filter.ManufacturerCountries)
	parts = filterByTags(parts, filter.Tags)
	return parts
}
