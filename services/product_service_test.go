package services

import (
	"errors"
	"sync"
	"testing"

	"ecommerce-system/models"
)

// TestCreateProduct comprueba el registro correcto de un producto.
func TestCreateProduct(t *testing.T) {
	service := NewProductService()

	product := models.Product{
		Name:        "Laptop",
		Description: "Laptop para programación",
		Price:       1200.50,
		Stock:       5,
	}

	createdProduct, err := service.CreateProduct(product)

	if err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}

	if createdProduct.ID != 1 {
		t.Errorf(
			"se esperaba el ID 1, pero se obtuvo %d",
			createdProduct.ID,
		)
	}

	if createdProduct.Name != "Laptop" {
		t.Errorf(
			"se esperaba Laptop, pero se obtuvo %s",
			createdProduct.Name,
		)
	}
}

// TestCreateProductInvalidPrice comprueba la validación del precio.
func TestCreateProductInvalidPrice(t *testing.T) {
	service := NewProductService()

	product := models.Product{
		Name:        "Producto inválido",
		Description: "Producto utilizado para una prueba",
		Price:       -20,
		Stock:       5,
	}

	_, err := service.CreateProduct(product)

	if err == nil {
		t.Fatal("se esperaba un error por precio inválido")
	}
}

// TestReserveStock comprueba el descuento correcto del inventario.
func TestReserveStock(t *testing.T) {
	service := NewProductService()

	product, err := service.CreateProduct(models.Product{
		Name:        "Teclado",
		Description: "Teclado mecánico",
		Price:       50,
		Stock:       5,
	})

	if err != nil {
		t.Fatalf("no fue posible crear el producto: %v", err)
	}

	items, total, err := service.ReserveStock(
		[]models.OrderItem{
			{
				ProductID: product.ID,
				Quantity:  2,
			},
		},
	)

	if err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}

	if total != 100 {
		t.Errorf(
			"se esperaba un total de 100, pero se obtuvo %.2f",
			total,
		)
	}

	if len(items) != 1 {
		t.Fatalf(
			"se esperaba un elemento, pero se obtuvieron %d",
			len(items),
		)
	}

	updatedProduct, err := service.GetProductByID(product.ID)

	if err != nil {
		t.Fatalf("no fue posible consultar el producto: %v", err)
	}

	if updatedProduct.Stock != 3 {
		t.Errorf(
			"se esperaba un stock de 3, pero se obtuvo %d",
			updatedProduct.Stock,
		)
	}
}

// TestConcurrentStockReservationDoesNotOversell demuestra la concurrencia.
// Dos solicitudes intentan comprar la última unidad simultáneamente.
func TestConcurrentStockReservationDoesNotOversell(
	t *testing.T,
) {
	service := NewProductService()

	product, err := service.CreateProduct(models.Product{
		Name:        "Monitor",
		Description: "Monitor con una sola unidad disponible",
		Price:       300,
		Stock:       1,
	})

	if err != nil {
		t.Fatalf("no fue posible crear el producto: %v", err)
	}

	results := make(chan error, 2)

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	// Se ejecutan dos reservas al mismo tiempo mediante goroutines.
	for i := 0; i < 2; i++ {
		go func() {
			defer waitGroup.Done()

			_, _, reservationError := service.ReserveStock(
				[]models.OrderItem{
					{
						ProductID: product.ID,
						Quantity:  1,
					},
				},
			)

			results <- reservationError
		}()
	}

	waitGroup.Wait()
	close(results)

	successfulReservations := 0
	insufficientStockErrors := 0

	for result := range results {
		switch {
		case result == nil:
			successfulReservations++

		case errors.Is(result, ErrInsufficientStock):
			insufficientStockErrors++

		default:
			t.Fatalf("se recibió un error inesperado: %v", result)
		}
	}

	if successfulReservations != 1 {
		t.Errorf(
			"se esperaba una reserva exitosa, pero se obtuvieron %d",
			successfulReservations,
		)
	}

	if insufficientStockErrors != 1 {
		t.Errorf(
			"se esperaba un error de stock, pero se obtuvieron %d",
			insufficientStockErrors,
		)
	}

	updatedProduct, err := service.GetProductByID(product.ID)

	if err != nil {
		t.Fatalf("no fue posible consultar el producto: %v", err)
	}

	if updatedProduct.Stock != 0 {
		t.Errorf(
			"se esperaba stock 0, pero se obtuvo %d",
			updatedProduct.Stock,
		)
	}
}
