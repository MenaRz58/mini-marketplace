package user

import (
	"errors"
	"mini-marketplace/users/internal/pkg/model"
	"mini-marketplace/users/internal/repository/memory"
)

type Controller struct {
	repo memory.Repo
	// Eliminamos productsClient y ordersClient. ¡Users es independiente!
}

func NewController(r memory.Repo) (*Controller, error) {
	if r == nil {
		return nil, errors.New("repository is nil")
	}
	return &Controller{repo: r}, nil
}

func (c *Controller) List() []model.User {
	return c.repo.List()
}

func (c *Controller) Get(id string) (model.User, error) {
	return c.repo.Get(id)
}

func (c *Controller) Create(id, name string) (model.User, error) {
	if id == "" || name == "" {
		return model.User{}, errors.New("invalid user")
	}
	u := model.User{ID: id, Name: name}
	if err := c.repo.Create(u); err != nil {
		return model.User{}, err
	}
	return u, nil
}

// ✅ Lógica para ValidateUser
func (c *Controller) Validate(id string) (bool, string) {
	u, err := c.repo.Get(id)
	if err != nil {
		return false, ""
	}
	return true, u.Name
}
