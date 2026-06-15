package interfaces

import "ecommerce-system/models"

type UserActions interface {
	Register(user models.User) error
	GetUsers() []models.User
}
