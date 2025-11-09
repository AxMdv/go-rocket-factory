package model

// PaymentMethod — доменный enum слоя модели (без зависимостей от protobuf).
type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)
