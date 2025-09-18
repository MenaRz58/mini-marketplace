package user

import "mini-marketplace/users/internal/pkg/model"

type Repo interface {
	List() []model.User
	Get(id string) (model.User, error)
	Create(u model.User) error
}

type Controller struct {
	repo Repo
}

func NewController(r Repo) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) List() []model.User                { return c.repo.List() }
func (c *Controller) Get(id string) (model.User, error) { return c.repo.Get(id) }
func (c *Controller) Create(u model.User) error         { return c.repo.Create(u) }
