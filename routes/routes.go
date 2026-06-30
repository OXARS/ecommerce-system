package routes

import (
	"ecommerce-system/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura y devuelve todas las rutas de la aplicación.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/health", controllers.Health)
	}

	return router
}
