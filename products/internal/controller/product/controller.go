package product

import (
	"errors"
	"mini-marketplace/products/internal/pkg/model"
)

type Repository interface {
	Get(id string) (model.Product, error)
	Create(p model.Product) error
	DecreaseStock(id string, qty int) error
	List() ([]model.Product, error)
}

type Controller struct {
	repo Repository
}

func NewController(r Repository) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) List() ([]model.Product, error) {
	return c.repo.List()
}

func (c *Controller) Get(id string) (model.Product, error) {
	return c.repo.Get(id)
}

func (c *Controller) Create(p model.Product) error {
	if p.ID == "" || p.Name == "" || p.Stock < 0 {
		return errors.New("invalid product")
	}
	return c.repo.Create(p)
}

func (c *Controller) Reserve(id string, qty int) (model.Product, error) {
	if err := c.repo.DecreaseStock(id, qty); err != nil {
		return model.Product{}, err
	}
	return c.repo.Get(id)
}
