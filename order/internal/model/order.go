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
	// UUID заказа.
	OrderUUID string `json:"order_uuid"`
	// UUID пользователя, сделавшего заказ.
	UserUUID string `json:"user_uuid"`
	// Список деталей, включённых в заказ.
	PartUUIDs []string `json:"part_uuids"`
	// Итоговая стоимость заказа.
	TotalPrice float64 `json:"total_price"`
	// UUID транзакции оплаты (если заказ оплачен).
	TransactionUUID *string        `json:"transaction_uuid,omitempty"`
	PaymentMethod   *PaymentMethod `json:"payment_method,omitempty"`
	Status          OrderStatus    `json:"status,omitempty"`
}

type OrderUpdateInfo struct {
	// UUID пользователя, сделавшего заказ.
	UserUUID *string `json:"user_uuid"`
	// Список деталей, включённых в заказ.
	PartUUIDs *[]string `json:"part_uuids"`
	// Итоговая стоимость заказа.
	TotalPrice *float64 `json:"total_price"`
	// UUID транзакции оплаты (если заказ оплачен).
	TransactionUUID *string        `json:"transaction_uuid,omitempty"`
	PaymentMethod   *PaymentMethod `json:"payment_method,omitempty"`
	Status          *OrderStatus   `json:"status,omitempty"`
}
