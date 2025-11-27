package model

import "gorm.io/gorm"

type Product struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `json:"name"`
	Price     float64        `json:"price"`
	Stock     int            `json:"stock"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
