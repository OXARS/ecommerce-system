package services

import (
	"errors"
	"testing"

	"ecommerce-system/models"
)

// TestCreateOrder comprueba la creación de un pedido completo.
func TestCreateOrder(t *testing.T) {
	productService := NewProductService()
	userService := NewUserService()

	user, err := userService.Register(models.User{
		Name:     "Oscar Caisatoa",
		Email:    "oscar@gmail.com",
		Password: "123456",
		Role:     "admin",
	})

	if err != nil {
		t.Fatalf("no fue posible crear el usuario: %v", err)
	}

	product, err := productService.CreateProduct(models.Product{
		Name:        "Laptop",
		Description: "Laptop para programación",
		Price:       1200.50,
		Stock:       5,
	})

	if err != nil {
		t.Fatalf("no fue posible crear el producto: %v", err)
	}

	orderService := NewOrderService(
		productService,
		userService,
	)

	order, err := orderService.CreateOrder(models.Order{
		UserID: user.ID,
		Items: []models.OrderItem{
			{
				ProductID: product.ID,
				Quantity:  2,
			},
		},
	})

	if err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}

	if order.ID != 1 {
		t.Errorf(
			"se esperaba el ID 1, pero se obtuvo %d",
			order.ID,
		)
	}

	if order.Total != 2401 {
		t.Errorf(
			"se esperaba un total de 2401, pero se obtuvo %.2f",
			order.Total,
		)
	}

	if order.Status != "pendiente" {
		t.Errorf(
			"se esperaba el estado pendiente, pero se obtuvo %s",
			order.Status,
		)
	}

	updatedProduct, err := productService.GetProductByID(
		product.ID,
	)

	if err != nil {
		t.Fatalf("no fue posible consultar el producto: %v", err)
	}

	if updatedProduct.Stock != 3 {
		t.Errorf(
			"se esperaba stock 3, pero se obtuvo %d",
			updatedProduct.Stock,
		)
	}
}

// TestCreateOrderInsufficientStock verifica un pedido inválido.
func TestCreateOrderInsufficientStock(t *testing.T) {
	productService := NewProductService()
	userService := NewUserService()

	user, err := userService.Register(models.User{
		Name:     "Maria Lopez",
		Email:    "maria@gmail.com",
		Password: "abcdef",
		Role:     "cliente",
	})

	if err != nil {
		t.Fatalf("no fue posible crear el usuario: %v", err)
	}

	product, err := productService.CreateProduct(models.Product{
		Name:        "Mouse",
		Description: "Mouse inalámbrico",
		Price:       25,
		Stock:       1,
	})

	if err != nil {
		t.Fatalf("no fue posible crear el producto: %v", err)
	}

	orderService := NewOrderService(
		productService,
		userService,
	)

	_, err = orderService.CreateOrder(models.Order{
		UserID: user.ID,
		Items: []models.OrderItem{
			{
				ProductID: product.ID,
				Quantity:  5,
			},
		},
	})

	if !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf(
			"se esperaba ErrInsufficientStock, pero se obtuvo: %v",
			err,
		)
	}

	orders := orderService.GetOrders()

	if len(orders) != 0 {
		t.Errorf(
			"el pedido inválido no debía guardarse, pero existen %d pedidos",
			len(orders),
		)
	}
}
