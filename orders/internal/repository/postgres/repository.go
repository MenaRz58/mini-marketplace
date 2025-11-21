package postgres

import (
	"mini-marketplace/orders/internal/pkg/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Get(id string) (*model.Order, error) {
	var o model.Order
	result := r.db.Preload("Products").First(&o, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &o, nil
}

func (r *OrderRepository) List() ([]model.Order, error) {
	var orders []model.Order

	result := r.db.Preload("Products").Find(&orders)

	return orders, result.Error
}

func (r *OrderRepository) Save(o *model.Order) error {
	return r.db.Create(o).Error
}
