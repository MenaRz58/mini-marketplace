package model

import "time"

type Order struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Total     float64
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	Products  []OrderProduct `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderProduct struct {
	ID        uint   `gorm:"primaryKey"`
	OrderID   string `gorm:"index"`
	ProductID int32
	Quantity  int32
	Price     float64
}
