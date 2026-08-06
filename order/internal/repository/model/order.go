package model

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type Order struct {
	OrderUUID       string         `json:"order_uuid"`
	UserUUID        string         `json:"user_uuid"`
	PartUUIDs       []string       `json:"part_uuids"`
	TotalPriceCents int64          `json:"total_price_cents"`
	TransactionUUID *string        `json:"transaction_uuid,omitempty"`
	PaymentMethod   *PaymentMethod `json:"payment_method,omitempty"`
	Status          OrderStatus    `json:"status,omitempty"`
}
