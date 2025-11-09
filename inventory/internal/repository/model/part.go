package model

import "time"

// Category — доменный enum
type Category int

const (
	CategoryUnknown Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimensions struct {
	Length float64 // см
	Width  float64 // см
	Height float64 // см
	Weight float64 // кг
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// // Value — oneof-представление для metadata.
// type Value struct {
// 	String *string
// 	Int64  *int64
// 	Double *float64
// 	Bool   *bool
// }

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]interface{} // can be one of (string, int64, float64, bool)
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
