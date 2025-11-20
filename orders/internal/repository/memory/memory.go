package memory

import (
	"errors"
	"mini-marketplace/orders/internal/pkg/model"
	"mini-marketplace/orders/internal/repository"
	"sync"
)

type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]model.Order
}

func NewInMemoryRepository() repository.Repo {
	return &InMemoryRepository{
		data: make(map[string]model.Order),
	}
}

func (r *InMemoryRepository) Create(o model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[o.ID]; exists {
		return errors.New("order already exists")
	}
	r.data[o.ID] = o
	return nil
}

func (r *InMemoryRepository) Get(id string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.data[id]
	if !ok {
		return model.Order{}, errors.New("order not found")
	}
	return o, nil
}

func (r *InMemoryRepository) List() ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]model.Order, 0, len(r.data))
	for _, o := range r.data {
		list = append(list, o)
	}
	return list, nil
}
