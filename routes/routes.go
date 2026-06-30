package routes

import (
	"ecommerce-system/controllers"
	"ecommerce-system/services"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura y devuelve las rutas de la aplicación.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	productService := services.NewProductService()
	productController := controllers.NewProductController(
		productService,
	)

	userService := services.NewUserService()
	userController := controllers.NewUserController(
		userService,
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

		users := api.Group("/users")
		{
			users.POST("", userController.Register)
			users.GET("", userController.GetUsers)
			users.GET("/:id", userController.GetUserByID)
		}
	}

	return router
}
