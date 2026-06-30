package interfaces

import "ecommerce-system/models"

// UserActions define las operaciones disponibles para usuarios.
type UserActions interface {
	Register(user models.User) (models.User, error)
	GetUsers() []models.User
	GetUserByID(id int) (models.User, error)
}
