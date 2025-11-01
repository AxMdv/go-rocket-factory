package part

import (
	"sync"

	repoModel "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/model"
)

type repository struct {
	mu    sync.RWMutex
	parts map[string]repoModel.Part
}

func NewRepository() *repository {
	repo := &repository{
		parts: make(map[string]repoModel.Part, 20),
	}
	repo.parts = createRepoParts(20)
	return repo
}
