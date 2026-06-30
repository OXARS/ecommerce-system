package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health verifica que el servicio web se encuentre disponible.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API del sistema e-commerce funcionando correctamente",
	})
}
