package services

import (
	"ecommerce-system/models"
	"errors"
)

type ProductService struct {
	Products []models.Product
}

// CreateProduct agrega un producto al sistema
func (ps *ProductService) CreateProduct(product models.Product) error {

	// Validación de stock
	if product.Stock <= 0 {
		return errors.New("stock inválido")
	}

	ps.Products = append(ps.Products, product)
	return nil
}

// GetProducts devuelve todos los productos
func (ps *ProductService) GetProducts() []models.Product {
	return ps.Products
}
