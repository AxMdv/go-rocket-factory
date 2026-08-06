package converter

import (
	"github.com/AxMdv/go-rocket-factory/order/internal/model"
	repoModel "github.com/AxMdv/go-rocket-factory/order/internal/repository/model"
)

func RepoOrderToModel(r repoModel.Order) model.Order {
	var tx *string
	if r.TransactionUUID != nil {
		tx = r.TransactionUUID
	}

	var pm *model.PaymentMethod
	if r.PaymentMethod != nil {
		val := model.PaymentMethod(*r.PaymentMethod)
		pm = &val
	}

	return model.Order{
		OrderUUID:       r.OrderUUID,
		UserUUID:        r.UserUUID,
		PartUUIDs:       r.PartUUIDs,
		TotalPriceCents: r.TotalPriceCents,
		TransactionUUID: tx,
		PaymentMethod:   pm,
		Status:          model.OrderStatus(r.Status),
	}
}

func OrderStatusToRepo(m model.OrderStatus) repoModel.OrderStatus {
	return repoModel.OrderStatus(m)
}

func PaymentMethodToRepo(m model.PaymentMethod) repoModel.PaymentMethod {
	return repoModel.PaymentMethod(m)
}

func ModelOrderToRepo(m model.Order) repoModel.Order {
	var tx *string
	if m.TransactionUUID != nil {
		tx = m.TransactionUUID
	}

	var pm *repoModel.PaymentMethod
	if m.PaymentMethod != nil {
		val := repoModel.PaymentMethod(*m.PaymentMethod)
		pm = &val
	}

	return repoModel.Order{
		OrderUUID:       m.OrderUUID,
		UserUUID:        m.UserUUID,
		PartUUIDs:       m.PartUUIDs,
		TotalPriceCents: m.TotalPriceCents,
		TransactionUUID: tx,
		PaymentMethod:   pm,
		Status:          repoModel.OrderStatus(m.Status),
	}
}
