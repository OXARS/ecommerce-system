package services

import (
	"errors"
	"strings"
	"sync"

	"ecommerce-system/models"
)

// Errores propios del módulo de productos.
var (
	// ErrProductNotFound se devuelve cuando no existe el producto solicitado.
	ErrProductNotFound = errors.New("producto no encontrado")

	// ErrInsufficientStock se devuelve cuando no existen unidades suficientes.
	ErrInsufficientStock = errors.New("stock insuficiente")
)

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
		return errors.New(
			"la descripción del producto es obligatoria",
		)
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

	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)

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
			updatedProduct.Name = strings.TrimSpace(
				updatedProduct.Name,
			)
			updatedProduct.Description = strings.TrimSpace(
				updatedProduct.Description,
			)

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

// ReserveStock valida y descuenta de manera segura el stock de un pedido.
func (ps *ProductService) ReserveStock(
	items []models.OrderItem,
) ([]models.OrderItem, float64, error) {
	if len(items) == 0 {
		return nil, 0, errors.New(
			"el pedido debe contener al menos un producto",
		)
	}

	// Solo una solicitud puede modificar el inventario a la vez.
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Agrupa las cantidades cuando un producto se repite.
	quantities := make(map[int]int)

	for _, item := range items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return nil, 0, errors.New(
				"el producto y la cantidad deben ser mayores que cero",
			)
		}

		quantities[item.ProductID] += item.Quantity
	}

	// Guarda la posición de cada producto dentro del inventario.
	productIndexes := make(map[int]int)

	for index, product := range ps.products {
		productIndexes[product.ID] = index
	}

	// Primero comprueba que todos los productos existan
	// y tengan suficiente stock.
	for productID, quantity := range quantities {
		index, exists := productIndexes[productID]

		if !exists {
			return nil, 0, ErrProductNotFound
		}

		if ps.products[index].Stock < quantity {
			return nil, 0, ErrInsufficientStock
		}
	}

	// Como todas las validaciones fueron correctas,
	// descuenta el stock de cada producto.
	for productID, quantity := range quantities {
		index := productIndexes[productID]
		ps.products[index].Stock -= quantity
	}

	reservedItems := make(
		[]models.OrderItem,
		0,
		len(items),
	)

	total := 0.0

	// Calcula el precio unitario y subtotal de cada elemento.
	for _, item := range items {
		index := productIndexes[item.ProductID]
		product := ps.products[index]

		item.UnitPrice = product.Price
		item.Subtotal = product.Price * float64(item.Quantity)

		total += item.Subtotal

		reservedItems = append(reservedItems, item)
	}

	return reservedItems, total, nil
}
