package user

import (
	"context"
	"errors"
	"mini-marketplace/users/internal/pkg/model"
)

type Repository interface {
	List() ([]model.User, error)
	Get(id string) (model.User, error)
	Create(u model.User) error
	GetWithCredentials(id string) (*model.User, error)
}

type Controller struct {
	repo Repository
}

func NewController(r Repository) (*Controller, error) {
	if r == nil {
		return nil, errors.New("repository is nil")
	}
	return &Controller{repo: r}, nil
}

func (c *Controller) List() ([]model.User, error) {
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

func (c *Controller) Validate(id string) (bool, string) {
	u, err := c.repo.Get(id)
	if err != nil {
		return false, ""
	}
	return true, u.Name
}

func (c *Controller) Login(ctx context.Context, id string, password string) (*model.User, error) {

	user, err := c.repo.GetWithCredentials(id)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	if user.Password != password {
		return nil, errors.New("contraseña incorrecta")
	}

	return user, nil
}
