package model

type OrderProduct struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price,omitempty"`
}

type Order struct {
	ID        string         `json:"id"`
	UserID    string         `json:"user_id"`
	Products  []OrderProduct `json:"products"`
	Total     float64        `json:"total"`
	CreatedAt int64          `json:"created_at"`
}
