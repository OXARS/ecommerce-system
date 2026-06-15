package controllers

import (
	"ecommerce-system/models"
	"ecommerce-system/services"
	"fmt"
)

type ProductController struct {
	Service services.ProductService
}

// AddProduct controla la creación de productos
func (pc *ProductController) AddProduct(product models.Product) {
	err := pc.Service.CreateProduct(product)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Producto agregado correctamente")
}
