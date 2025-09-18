package repository

import (
	"context"
	"errors"
	"sync"

	"mini-marketplace/metadata/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

func New() *Repository {
	return &Repository{data: map[string]*model.Metadata{}}
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (r *Repository) Put(ctx context.Context, id string, m *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = m
	return nil
}

var ErrNotFound = errors.New("metadata not found")
