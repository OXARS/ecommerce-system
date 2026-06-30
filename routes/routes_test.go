package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// performJSONRequest crea una solicitud HTTP de prueba.
func performJSONRequest(
	handler http.Handler,
	method string,
	path string,
	body interface{},
) *httptest.ResponseRecorder {
	var reader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)

		if err != nil {
			panic(err)
		}

		reader = bytes.NewBuffer(jsonBody)
	}

	request := httptest.NewRequest(method, path, reader)

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	return response
}

// TestEcommerceAPIFlow prueba la integración de varios módulos.
func TestEcommerceAPIFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := SetupRouter()

	userResponse := performJSONRequest(
		router,
		http.MethodPost,
		"/api/users",
		map[string]interface{}{
			"name":     "Oscar Caisatoa",
			"email":    "oscar@gmail.com",
			"password": "123456",
			"role":     "admin",
		},
	)

	if userResponse.Code != http.StatusCreated {
		t.Fatalf(
			"se esperaba código 201 al crear usuario, pero se obtuvo %d: %s",
			userResponse.Code,
			userResponse.Body.String(),
		)
	}

	productResponse := performJSONRequest(
		router,
		http.MethodPost,
		"/api/products",
		map[string]interface{}{
			"name":        "Laptop",
			"description": "Laptop para programación",
			"price":       1200.50,
			"stock":       5,
		},
	)

	if productResponse.Code != http.StatusCreated {
		t.Fatalf(
			"se esperaba código 201 al crear producto, pero se obtuvo %d: %s",
			productResponse.Code,
			productResponse.Body.String(),
		)
	}

	orderResponse := performJSONRequest(
		router,
		http.MethodPost,
		"/api/orders",
		map[string]interface{}{
			"user_id": 1,
			"items": []map[string]interface{}{
				{
					"product_id": 1,
					"quantity":   2,
				},
			},
		},
	)

	if orderResponse.Code != http.StatusCreated {
		t.Fatalf(
			"se esperaba código 201 al crear pedido, pero se obtuvo %d: %s",
			orderResponse.Code,
			orderResponse.Body.String(),
		)
	}

	inventoryResponse := performJSONRequest(
		router,
		http.MethodGet,
		"/api/inventory",
		nil,
	)

	if inventoryResponse.Code != http.StatusOK {
		t.Fatalf(
			"se esperaba código 200 al consultar inventario, pero se obtuvo %d",
			inventoryResponse.Code,
		)
	}

	var inventoryData struct {
		Total int `json:"total"`
		Data  []struct {
			ID    int `json:"id"`
			Stock int `json:"stock"`
		} `json:"data"`
	}

	err := json.Unmarshal(
		inventoryResponse.Body.Bytes(),
		&inventoryData,
	)

	if err != nil {
		t.Fatalf(
			"no fue posible interpretar la respuesta JSON: %v",
			err,
		)
	}

	if inventoryData.Total != 1 {
		t.Errorf(
			"se esperaba un producto, pero se obtuvo %d",
			inventoryData.Total,
		)
	}

	if len(inventoryData.Data) != 1 {
		t.Fatalf(
			"se esperaba un elemento de inventario, pero se obtuvieron %d",
			len(inventoryData.Data),
		)
	}

	if inventoryData.Data[0].Stock != 3 {
		t.Errorf(
			"se esperaba stock 3 después del pedido, pero se obtuvo %d",
			inventoryData.Data[0].Stock,
		)
	}
}

// TestCreateInvalidProduct comprueba una respuesta HTTP de error.
func TestCreateInvalidProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := SetupRouter()

	response := performJSONRequest(
		router,
		http.MethodPost,
		"/api/products",
		map[string]interface{}{
			"name":        "Producto inválido",
			"description": "Producto con precio incorrecto",
			"price":       -50,
			"stock":       2,
		},
	)

	if response.Code != http.StatusBadRequest {
		t.Errorf(
			"se esperaba código 400, pero se obtuvo %d",
			response.Code,
		)
	}
}
