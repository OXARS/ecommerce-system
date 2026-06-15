package controllers

import (
	"ecommerce-system/models"
	"ecommerce-system/services"
	"fmt"
)

type UserController struct {
	Service services.UserService
}

// AddUser controla el registro de usuarios
func (uc *UserController) AddUser(user models.User) {
	err := uc.Service.Register(user)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Usuario registrado correctamente")
}
