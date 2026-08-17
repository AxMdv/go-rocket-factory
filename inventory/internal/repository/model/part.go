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
	UUID          string                 `bson:"_id"`
	Name          string                 `bson:"name"`
	Description   string                 `bson:"description"`
	Price         float64                `bson:"price"`
	StockQuantity int64                  `bson:"stock_quantity"`
	Category      Category               `bson:"category"`
	Dimensions    *Dimensions            `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer          `bson:"manufacturer,omitempty"`
	Tags          []string               `bson:"tags"`
	Metadata      map[string]interface{} `bson:"metadata"`
	CreatedAt     time.Time              `bson:"created_at"`
	UpdatedAt     time.Time              `bson:"updated_at"`
}
