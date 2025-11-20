package memory

import (
	"errors"
	"sync"

	"mini-marketplace/products/internal/pkg/model"
)

// Repo define la interfaz que el controlador va a usar
type Repo interface {
	List() []model.Product
	Get(id string) (model.Product, error)
	Create(p model.Product) error
	DecreaseStock(id string, qty int) error
}

// ProductRepository es una implementación en memoria de Repo
type ProductRepository struct {
	mu   sync.RWMutex
	data map[string]model.Product
}

// NewProductRepository crea un repositorio con datos iniciales
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		data: map[string]model.Product{
			"p1": {ID: "p1", Name: "Laptop", Price: 1500.00, Stock: 5},
			"p2": {ID: "p2", Name: "Headphones", Price: 120.00, Stock: 30},
		},
	}
}

// List devuelve todos los productos
func (r *ProductRepository) List() []model.Product {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.Product, 0, len(r.data))
	for _, v := range r.data {
		out = append(out, v)
	}
	return out
}

// Get devuelve un producto por ID
func (r *ProductRepository) Get(id string) (model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.data[id]
	if !ok {
		return model.Product{}, errors.New("product not found")
	}
	return p, nil
}

// Create agrega un nuevo producto
func (r *ProductRepository) Create(p model.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[p.ID]; exists {
		return errors.New("product already exists")
	}
	r.data[p.ID] = p
	return nil
}

// DecreaseStock reduce el stock de un producto
func (r *ProductRepository) DecreaseStock(id string, qty int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.data[id]
	if !ok {
		return errors.New("product not found")
	}
	if p.Stock < qty {
		return errors.New("insufficient stock")
	}
	p.Stock -= qty
	r.data[id] = p
	return nil
}
