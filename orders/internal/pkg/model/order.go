package model

type Order struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	Total     float64
	CreatedAt int64
	Products  []OrderProduct `gorm:"foreignKey:OrderID"`
}

type OrderProduct struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   string
	ProductID string
	Quantity  int32
	Price     float64
}
