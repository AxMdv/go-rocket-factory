package converter

import (
	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	orderV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/openapi/order/v1"
)

// OrderDtoToModel конвертирует транспортный DTO в доменную модель.
func OrderDtoToModel(d orderV1.OrderDto) model.Order {
	// TransactionUUID
	var tx *string
	if v, ok := d.TransactionUUID.Get(); ok {
		tx = &v
	}

	// PaymentMethod
	var pm *model.PaymentMethod
	if v, ok := d.PaymentMethod.Get(); ok {
		mv := PaymentMethodDtoToModel(v)
		pm = &mv
	}

	return model.Order{
		OrderUUID:       d.GetOrderUUID(),
		UserUUID:        d.GetUserUUID(),
		PartUUIDs:       d.GetPartUuids(),
		TotalPrice:      d.GetTotalPrice(),
		TransactionUUID: tx,
		PaymentMethod:   pm,
		Status:          OrderStatusDtoToModel(d.GetStatus()),
	}
}

// OrderModelToDto конвертирует доменную модель в транспортный DTO.
func OrderModelToDto(m model.Order) orderV1.OrderDto {
	var tx orderV1.OptNilString
	if m.TransactionUUID != nil {
		tx = orderV1.NewOptNilString(*m.TransactionUUID)
	}

	var pm orderV1.OptPaymentMethod
	if m.PaymentMethod != nil {
		pm = orderV1.NewOptPaymentMethod(PaymentMethodModelToDto(*m.PaymentMethod))
	}

	return orderV1.OrderDto{
		OrderUUID:       m.OrderUUID,
		UserUUID:        m.UserUUID,
		PartUuids:       m.PartUUIDs,
		TotalPrice:      m.TotalPrice,
		TransactionUUID: tx,
		PaymentMethod:   pm,
		Status:          OrderStatusModelToDto(m.Status),
	}
}

// ----- Enum mapping ----- //

func OrderStatusDtoToModel(s orderV1.OrderStatus) model.OrderStatus {
	switch s {
	case orderV1.OrderStatusPAID:
		return model.OrderStatusPAID
	case orderV1.OrderStatusCANCELLED:
		return model.OrderStatusCANCELLED
	default:
		return model.OrderStatusPENDINGPAYMENT
	}
}

func OrderStatusModelToDto(s model.OrderStatus) orderV1.OrderStatus {
	switch s {
	case model.OrderStatusPAID:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCANCELLED:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func PaymentMethodDtoToModel(pm orderV1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func PaymentMethodModelToDto(pm model.PaymentMethod) orderV1.PaymentMethod {
	switch pm {
	case model.PaymentMethodCARD:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCREDITCARD:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodINVESTORMONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}
