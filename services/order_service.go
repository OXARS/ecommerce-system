package services

import (
	"errors"
	"sync"
	"time"

	"ecommerce-system/interfaces"
	"ecommerce-system/models"
)

// ErrOrderNotFound se devuelve cuando no existe el pedido solicitado.
var ErrOrderNotFound = errors.New("pedido no encontrado")

// OrderService administra los pedidos almacenados en memoria.
type OrderService struct {
	mu             sync.RWMutex
	orders         []models.Order
	nextID         int
	productService interfaces.ProductActions
	userService    interfaces.UserActions
}

// NewOrderService crea e inicializa el servicio de pedidos.
func NewOrderService(
	productService interfaces.ProductActions,
	userService interfaces.UserActions,
) *OrderService {
	return &OrderService{
		orders:         make([]models.Order, 0),
		nextID:         1,
		productService: productService,
		userService:    userService,
	}
}

// CreateOrder valida al usuario, reserva el stock y registra el pedido.
func (os *OrderService) CreateOrder(
	order models.Order,
) (models.Order, error) {
	if order.UserID <= 0 {
		return models.Order{}, errors.New(
			"el usuario es obligatorio",
		)
	}

	if len(order.Items) == 0 {
		return models.Order{}, errors.New(
			"el pedido debe contener al menos un producto",
		)
	}

	// Comprueba que el usuario exista.
	if _, err := os.userService.GetUserByID(order.UserID); err != nil {
		return models.Order{}, err
	}

	// Valida y descuenta el inventario.
	reservedItems, total, err := os.productService.ReserveStock(
		order.Items,
	)
	if err != nil {
		return models.Order{}, err
	}

	order.Items = reservedItems
	order.Total = total
	order.Status = "pendiente"
	order.CreatedAt = time.Now().UTC()

	os.mu.Lock()
	defer os.mu.Unlock()

	order.ID = os.nextID
	os.nextID++

	os.orders = append(os.orders, order)

	return order, nil
}

// GetOrders devuelve todos los pedidos registrados.
func (os *OrderService) GetOrders() []models.Order {
	os.mu.RLock()
	defer os.mu.RUnlock()

	ordersCopy := make([]models.Order, len(os.orders))
	copy(ordersCopy, os.orders)

	return ordersCopy
}

// GetOrderByID busca un pedido por su identificador.
func (os *OrderService) GetOrderByID(
	id int,
) (models.Order, error) {
	os.mu.RLock()
	defer os.mu.RUnlock()

	for _, order := range os.orders {
		if order.ID == id {
			return order, nil
		}
	}

	return models.Order{}, ErrOrderNotFound
}
