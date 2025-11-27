package product

import (
	"errors"
	"mini-marketplace/products/internal/pkg/model"
)

type Repository interface {
	Get(id uint) (model.Product, error)
	DecreaseStock(id uint, qty int) error
	List() ([]model.Product, error)
	Create(p *model.Product) error
	Update(p *model.Product) error
	Delete(id uint) error
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

func (c *Controller) Get(id uint) (model.Product, error) {
	return c.repo.Get(id)
}

func (c *Controller) Create(p *model.Product) error {
	if p.Name == "" || p.Stock < 0 {
		return errors.New("invalid product")
	}
	return c.repo.Create(p)
}

func (c *Controller) Reserve(id uint, qty int) (model.Product, error) {
	if err := c.repo.DecreaseStock(id, qty); err != nil {
		return model.Product{}, err
	}
	return c.repo.Get(id)
}

func (c *Controller) Update(p *model.Product) error {
	if p.ID == 0 {
		return errors.New("id requerido")
	}
	return c.repo.Update(p)
}

func (c *Controller) Delete(id uint) error {
	return c.repo.Delete(id)
}
