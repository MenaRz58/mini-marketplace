package memory

import (
	"errors"
	"mini-marketplace/orders/internal/controller/order"
	"mini-marketplace/orders/internal/pkg/model"
	"sync"
)

type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]model.Order
}

func NewInMemoryRepository() order.Repo {
	return &InMemoryRepository{
		data: make(map[string]model.Order),
	}
}

// Método Create (antes Save)
func (r *InMemoryRepository) Create(o model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[o.ID]; exists {
		// Ya existe: no restamos stock ni guardamos otra vez
		return errors.New("order already exists")
	}

	r.data[o.ID] = o
	return nil
}

func (r *InMemoryRepository) Get(id string) (model.Order, error) {
	order, ok := r.data[id]
	if !ok {
		return model.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (r *InMemoryRepository) List() []model.Order {
	orders := make([]model.Order, 0, len(r.data))
	for _, o := range r.data {
		orders = append(orders, o)
	}
	return orders
}
