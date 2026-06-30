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

// UserController recibe las solicitudes web de usuarios.
type UserController struct {
	Service interfaces.UserActions
}

// NewUserController crea un controlador de usuarios.
func NewUserController(
	service interfaces.UserActions,
) *UserController {
	return &UserController{
		Service: service,
	}
}

// Register registra un usuario recibido mediante JSON.
func (uc *UserController) Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos del usuario inválidos",
			"details": err.Error(),
		})
		return
	}

	createdUser, err := uc.Service.Register(user)

	if errors.Is(err, services.ErrEmailAlreadyRegistered) {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado correctamente",
		"data":    createdUser,
	})
}

// GetUsers devuelve todos los usuarios registrados.
func (uc *UserController) GetUsers(c *gin.Context) {
	users := uc.Service.GetUsers()

	c.JSON(http.StatusOK, gin.H{
		"total": len(users),
		"data":  users,
	})
}

// GetUserByID busca un usuario mediante el ID de la URL.
func (uc *UserController) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El ID debe ser un número entero positivo",
		})
		return
	}

	user, err := uc.Service.GetUserByID(id)

	if errors.Is(err, services.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No fue posible consultar el usuario",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
