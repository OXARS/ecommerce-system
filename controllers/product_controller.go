package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"ecommerce-system/interfaces"
	"ecommerce-system/models"
	"ecommerce-system/services"

	"github.com/gin-gonic/gin"
)

// ProductController recibe las solicitudes web de productos.
type ProductController struct {
	Service interfaces.ProductActions
}

// NewProductController crea un controlador de productos.
func NewProductController(
	service interfaces.ProductActions,
) *ProductController {
	return &ProductController{
		Service: service,
	}
}

// CreateProduct registra un producto recibido mediante JSON.
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos del producto inválidos",
			"details": err.Error(),
		})
		return
	}

	createdProduct, err := pc.Service.CreateProduct(product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Producto creado correctamente",
		"data":    createdProduct,
	})
}

// GetProducts devuelve todos los productos registrados.
func (pc *ProductController) GetProducts(c *gin.Context) {
	products := pc.Service.GetProducts()

	c.JSON(http.StatusOK, gin.H{
		"total": len(products),
		"data":  products,
	})
}

// GetProductByID busca un producto mediante el ID de la URL.
func (pc *ProductController) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El ID debe ser un número entero positivo",
		})
		return
	}

	product, err := pc.Service.GetProductByID(id)

	if errors.Is(err, services.ErrProductNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No fue posible consultar el producto",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

// UpdateProduct actualiza un producto existente.
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El ID debe ser un número entero positivo",
		})
		return
	}

	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos del producto inválidos",
			"details": err.Error(),
		})
		return
	}

	updatedProduct, err := pc.Service.UpdateProduct(id, product)

	if errors.Is(err, services.ErrProductNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Producto actualizado correctamente",
		"data":    updatedProduct,
	})
}

// DeleteProduct elimina un producto existente.
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El ID debe ser un número entero positivo",
		})
		return
	}

	err = pc.Service.DeleteProduct(id)

	if errors.Is(err, services.ErrProductNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No fue posible eliminar el producto",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Producto eliminado correctamente",
	})
}

// GetInventory devuelve el stock actual de todos los productos.
func (pc *ProductController) GetInventory(c *gin.Context) {
	type InventoryItem struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Stock     int    `json:"stock"`
		Available bool   `json:"available"`
	}

	products := pc.Service.GetProducts()

	inventory := make(
		[]InventoryItem,
		0,
		len(products),
	)

	for _, product := range products {
		inventory = append(inventory, InventoryItem{
			ID:        product.ID,
			Name:      product.Name,
			Stock:     product.Stock,
			Available: product.Stock > 0,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(inventory),
		"data":  inventory,
	})
}
