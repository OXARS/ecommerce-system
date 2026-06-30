package interfaces

import "ecommerce-system/models"

// OrderActions define las operaciones disponibles para pedidos.
type OrderActions interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrders() []models.Order
	GetOrderByID(id int) (models.Order, error)
}
