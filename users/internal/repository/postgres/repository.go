package postgres

import (
	"mini-marketplace/users/internal/pkg/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) List() ([]model.User, error) {
	var users []model.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *UserRepository) Get(id string) (model.User, error) {
	var u model.User
	if result := r.db.First(&u, "id = ?", id); result.Error != nil {
		return model.User{}, result.Error
	}
	return u, nil
}

func (r *UserRepository) Create(u model.User) error {
	return r.db.Create(&u).Error
}

func (r *UserRepository) GetWithCredentials(id string) (*model.User, error) {
	var u model.User

	result := r.db.Where("id = ?", id).First(&u)

	if result.Error != nil {
		return nil, result.Error
	}

	return &u, nil
}
