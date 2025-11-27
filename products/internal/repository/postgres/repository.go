package postgres

import (
	"errors"

	"gorm.io/gorm/clause"

	"mini-marketplace/products/internal/pkg/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Get(id uint) (model.Product, error) {
	var p model.Product
	// Busca por ID
	result := r.db.First(&p, "id = ?", id)
	return p, result.Error
}

func (r *ProductRepository) List() ([]model.Product, error) {
	var products []model.Product
	result := r.db.Find(&products)
	return products, result.Error
}

func (r *ProductRepository) Create(p *model.Product) error {
	result := r.db.Create(p)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProductRepository) Update(p *model.Product) error {
	return r.db.Save(p).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) DecreaseStock(id uint, qty int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var p model.Product

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("product not found")
			}
			return err
		}

		if p.Stock < qty {
			return errors.New("insufficient stock")
		}

		p.Stock -= qty

		return tx.Save(&p).Error
	})
}
