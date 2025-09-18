package product

import "mini-marketplace/products/internal/pkg/model"

type Repo interface {
	List() []model.Product
	Get(id string) (model.Product, error)
	Create(p model.Product) error
	DecreaseStock(id string, qty int) error
}

type Controller struct {
	repo Repo
}

func NewController(r Repo) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) List() []model.Product {
	return c.repo.List()
}

func (c *Controller) Get(id string) (model.Product, error) {
	return c.repo.Get(id)
}

func (c *Controller) Create(p model.Product) error {
	return c.repo.Create(p)
}

func (c *Controller) Reserve(id string, qty int) error {
	return c.repo.DecreaseStock(id, qty)
}
