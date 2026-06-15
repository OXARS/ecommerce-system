package services

import (
	"ecommerce-system/models"
	"errors"
)

type UserService struct {
	Users []models.User
}

// Register registra un nuevo usuario
func (us *UserService) Register(user models.User) error {

	if user.Name == "" {
		return errors.New("nombre vacío")
	}

	us.Users = append(us.Users, user)
	return nil
}

// GetUsers devuelve todos los usuarios
func (us *UserService) GetUsers() []models.User {
	return us.Users
}
