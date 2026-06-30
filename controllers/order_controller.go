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

// OrderController recibe las solicitudes web de pedidos.
type OrderController struct {
	Service interfaces.OrderActions
}

// NewOrderController crea un controlador de pedidos.
func NewOrderController(
	service interfaces.OrderActions,
) *OrderController {
	return &OrderController{
		Service: service,
	}
}

// CreateOrder registra un pedido recibido mediante JSON.
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos del pedido inválidos",
			"details": err.Error(),
		})
		return
	}

	createdOrder, err := oc.Service.CreateOrder(order)

	switch {
	case errors.Is(err, services.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return

	case errors.Is(err, services.ErrProductNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return

	case errors.Is(err, services.ErrInsufficientStock):
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return

	case err != nil:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pedido creado correctamente",
		"data":    createdOrder,
	})
}

// GetOrders devuelve todos los pedidos registrados.
func (oc *OrderController) GetOrders(c *gin.Context) {
	orders := oc.Service.GetOrders()

	c.JSON(http.StatusOK, gin.H{
		"total": len(orders),
		"data":  orders,
	})
}

// GetOrderByID busca un pedido por el ID recibido en la URL.
func (oc *OrderController) GetOrderByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El ID debe ser un número entero positivo",
		})
		return
	}

	order, err := oc.Service.GetOrderByID(id)

	if errors.Is(err, services.ErrOrderNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No fue posible consultar el pedido",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}
