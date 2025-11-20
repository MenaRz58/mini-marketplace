package repository

import "mini-marketplace/orders/internal/pkg/model"

type Repo interface {
	Create(order model.Order) error
	Get(id string) (model.Order, error)
	List() ([]model.Order, error)
}
