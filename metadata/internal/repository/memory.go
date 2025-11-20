package repository

import (
	"context"
	"errors"
	"mini-marketplace/metadata/pkg/model"
)

type Repository struct {
	data map[string]*model.Metadata
}

func New() *Repository {
	return &Repository{
		data: make(map[string]*model.Metadata),
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	m, ok := r.data[id]
	if !ok {
		return nil, errors.New("metadata not found")
	}
	return m, nil
}

func (r *Repository) Put(ctx context.Context, id string, m *model.Metadata) error {
	r.data[id] = m
	return nil
}
