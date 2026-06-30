package interfaces

import "ecommerce-system/models"

// ProductActions define las operaciones disponibles para productos.
type ProductActions interface {
	CreateProduct(product models.Product) (models.Product, error)
	GetProducts() []models.Product
	GetProductByID(id int) (models.Product, error)
	UpdateProduct(id int, product models.Product) (models.Product, error)
	DeleteProduct(id int) error
}
