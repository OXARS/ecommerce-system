package routes

import (
	"ecommerce-system/controllers"
	"ecommerce-system/services"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura y devuelve las rutas de la aplicación.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Se crea una sola instancia compartida del servicio.
	productService := services.NewProductService()

	// El controlador recibe el servicio mediante su interfaz.
	productController := controllers.NewProductController(
		productService,
	)

	api := router.Group("/api")
	{
		api.GET("/health", controllers.Health)

		products := api.Group("/products")
		{
			products.POST("", productController.CreateProduct)
			products.GET("", productController.GetProducts)
			products.GET("/:id", productController.GetProductByID)
			products.PUT("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}
	}

	return router
}
