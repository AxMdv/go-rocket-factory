package order

import (
	"sync"

	repoModel "github.com/AxMdv/go-rocket-factory/order/internal/repository/model"
)

type repository struct {
	mu     sync.RWMutex
	orders map[string]repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]repoModel.Order),
	}
}
