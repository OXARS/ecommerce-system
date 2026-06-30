package routes

import (
	"ecommerce-system/controllers"
	"ecommerce-system/services"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura y devuelve las rutas de la aplicación.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Servicios compartidos por toda la aplicación.
	productService := services.NewProductService()
	userService := services.NewUserService()

	// El pedido utiliza los mismos servicios de usuarios y productos.
	orderService := services.NewOrderService(
		productService,
		userService,
	)

	// Controladores.
	productController := controllers.NewProductController(
		productService,
	)

	userController := controllers.NewUserController(
		userService,
	)

	orderController := controllers.NewOrderController(
		orderService,
	)

	api := router.Group("/api")
	{
		api.GET("/health", controllers.Health)
		api.GET("/inventory", productController.GetInventory)

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

		orders := api.Group("/orders")
		{
			orders.POST("", orderController.CreateOrder)
			orders.GET("", orderController.GetOrders)
			orders.GET("/:id", orderController.GetOrderByID)
		}
	}

	return router
}
