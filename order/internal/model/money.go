package model

import "math"

const centsInUnit = 100

func PriceToCents(price float64) int64 {
	return int64(math.Round(price * centsInUnit))
}

func CentsToPrice(cents int64) float64 {
	return float64(cents) / centsInUnit
}
