package models

// User representa un usuario registrado en el e-commerce.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin cliente"`
}
