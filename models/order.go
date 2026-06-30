package models

import "time"

// OrderItem representa un producto y su cantidad dentro de un pedido.
type OrderItem struct {
	ProductID int     `json:"product_id" binding:"required,gt=0"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price,omitempty"`
	Subtotal  float64 `json:"subtotal,omitempty"`
}

// Order representa una compra realizada por un usuario.
type Order struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id" binding:"required,gt=0"`
	Items     []OrderItem `json:"items" binding:"required,min=1,dive"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}
