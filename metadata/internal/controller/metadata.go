package controller

import (
	"context"

	"mini-marketplace/metadata/pkg/model"
)

type Repository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
	Put(ctx context.Context, id string, m *model.Metadata) error
}

type Controller struct {
	repo Repository
}

func New(r Repository) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	return c.repo.Get(ctx, id)
}

func (c *Controller) Put(ctx context.Context, id string, m *model.Metadata) error {
	return c.repo.Put(ctx, id, m)
}
