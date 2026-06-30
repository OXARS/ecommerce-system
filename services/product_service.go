package services

import (
	"errors"
	"strings"
	"sync"

	"ecommerce-system/models"
)

// ErrProductNotFound se devuelve cuando no existe el producto solicitado.
var ErrProductNotFound = errors.New("producto no encontrado")

// ProductService administra los productos guardados en memoria.
type ProductService struct {
	mu       sync.RWMutex
	products []models.Product
	nextID   int
}

// NewProductService crea e inicializa el servicio de productos.
func NewProductService() *ProductService {
	return &ProductService{
		products: make([]models.Product, 0),
		nextID:   1,
	}
}

// validateProduct comprueba que los datos del producto sean válidos.
func validateProduct(product models.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("el nombre del producto es obligatorio")
	}

	if strings.TrimSpace(product.Description) == "" {
		return errors.New("la descripción del producto es obligatoria")
	}

	if product.Price <= 0 {
		return errors.New("el precio debe ser mayor que cero")
	}

	if product.Stock < 0 {
		return errors.New("el stock no puede ser negativo")
	}

	return nil
}

// CreateProduct registra un nuevo producto.
func (ps *ProductService) CreateProduct(
	product models.Product,
) (models.Product, error) {
	if err := validateProduct(product); err != nil {
		return models.Product{}, err
	}

	// Bloquea temporalmente la escritura para evitar conflictos
	// cuando varias solicitudes crean productos simultáneamente.
	ps.mu.Lock()
	defer ps.mu.Unlock()

	product.ID = ps.nextID
	ps.nextID++

	ps.products = append(ps.products, product)

	return product, nil
}

// GetProducts devuelve una copia de todos los productos.
func (ps *ProductService) GetProducts() []models.Product {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	productsCopy := make([]models.Product, len(ps.products))
	copy(productsCopy, ps.products)

	return productsCopy
}

// GetProductByID busca un producto por su identificador.
func (ps *ProductService) GetProductByID(
	id int,
) (models.Product, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, product := range ps.products {
		if product.ID == id {
			return product, nil
		}
	}

	return models.Product{}, ErrProductNotFound
}

// UpdateProduct modifica los datos de un producto existente.
func (ps *ProductService) UpdateProduct(
	id int,
	updatedProduct models.Product,
) (models.Product, error) {
	if err := validateProduct(updatedProduct); err != nil {
		return models.Product{}, err
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	for index, product := range ps.products {
		if product.ID == id {
			updatedProduct.ID = id
			ps.products[index] = updatedProduct

			return updatedProduct, nil
		}
	}

	return models.Product{}, ErrProductNotFound
}

// DeleteProduct elimina un producto por su identificador.
func (ps *ProductService) DeleteProduct(id int) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for index, product := range ps.products {
		if product.ID == id {
			ps.products = append(
				ps.products[:index],
				ps.products[index+1:]...,
			)

			return nil
		}
	}

	return ErrProductNotFound
}
