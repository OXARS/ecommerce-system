package interfaces

import "ecommerce-system/models"

type ProductActions interface {
	CreateProduct(product models.Product) error
	GetProducts() []models.Product
}
