package memory

import (
	"errors"
	"sync"

	"mini-marketplace/users/internal/pkg/model"
)

type UserRepository struct {
	mu   sync.RWMutex
	data map[string]model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		data: map[string]model.User{
			"u1": {ID: "u1", Name: "Alice"},
			"u2": {ID: "u2", Name: "Bob"},
		},
	}
}

func (r *UserRepository) List() []model.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.User, 0, len(r.data))
	for _, v := range r.data {
		out = append(out, v)
	}
	return out
}

func (r *UserRepository) Get(id string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.data[id]
	if !ok {
		return model.User{}, errors.New("user not found")
	}
	return u, nil
}

func (r *UserRepository) Create(u model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[u.ID]; exists {
		return errors.New("user already exists")
	}
	r.data[u.ID] = u
	return nil
}
