package services

import (
	"errors"
	"testing"

	"ecommerce-system/models"
)

// TestRegisterUser comprueba el registro correcto de un usuario.
func TestRegisterUser(t *testing.T) {
	service := NewUserService()

	user := models.User{
		Name:     "Oscar Caisatoa",
		Email:    "oscar@gmail.com",
		Password: "123456",
		Role:     "admin",
	}

	createdUser, err := service.Register(user)

	if err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}

	if createdUser.ID != 1 {
		t.Errorf(
			"se esperaba el ID 1, pero se obtuvo %d",
			createdUser.ID,
		)
	}

	if createdUser.Password != "" {
		t.Error("la contraseña no debe aparecer en la respuesta")
	}

	if len(service.users) != 1 {
		t.Fatalf(
			"se esperaba un usuario guardado, pero existen %d",
			len(service.users),
		)
	}

	if service.users[0].Password == "123456" {
		t.Error("la contraseña fue almacenada sin protección")
	}
}

// TestDuplicateEmail comprueba que el correo no pueda repetirse.
func TestDuplicateEmail(t *testing.T) {
	service := NewUserService()

	user := models.User{
		Name:     "Oscar Caisatoa",
		Email:    "oscar@gmail.com",
		Password: "123456",
		Role:     "admin",
	}

	_, err := service.Register(user)

	if err != nil {
		t.Fatalf("no fue posible registrar el primer usuario: %v", err)
	}

	_, err = service.Register(user)

	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Fatalf(
			"se esperaba ErrEmailAlreadyRegistered, pero se obtuvo: %v",
			err,
		)
	}
}

// TestGetUserNotFound comprueba la búsqueda de un usuario inexistente.
func TestGetUserNotFound(t *testing.T) {
	service := NewUserService()

	_, err := service.GetUserByID(99)

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf(
			"se esperaba ErrUserNotFound, pero se obtuvo: %v",
			err,
		)
	}
}
