package services

import (
	"errors"
	"net/mail"
	"strings"
	"sync"

	"ecommerce-system/models"

	"golang.org/x/crypto/bcrypt"
)

// Errores propios del módulo de usuarios.
var (
	ErrUserNotFound           = errors.New("usuario no encontrado")
	ErrEmailAlreadyRegistered = errors.New("el correo ya está registrado")
)

// UserService administra los usuarios almacenados en memoria.
type UserService struct {
	mu     sync.RWMutex
	users  []models.User
	nextID int
}

// NewUserService crea e inicializa el servicio de usuarios.
func NewUserService() *UserService {
	return &UserService{
		users:  make([]models.User, 0),
		nextID: 1,
	}
}

// validateUser comprueba que la información sea válida.
func validateUser(user models.User) error {
	if len(user.Name) < 3 {
		return errors.New("el nombre debe tener al menos 3 caracteres")
	}

	address, err := mail.ParseAddress(user.Email)
	if err != nil || address.Address != user.Email {
		return errors.New("el correo electrónico no es válido")
	}

	if len(user.Password) < 6 {
		return errors.New("la contraseña debe tener al menos 6 caracteres")
	}

	if user.Role != "admin" && user.Role != "cliente" {
		return errors.New("el rol debe ser admin o cliente")
	}

	return nil
}

// sanitizeUser crea una copia sin mostrar la contraseña.
func sanitizeUser(user models.User) models.User {
	user.Password = ""
	return user
}

// Register registra un usuario y protege su contraseña.
func (us *UserService) Register(
	user models.User,
) (models.User, error) {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Role = strings.ToLower(strings.TrimSpace(user.Role))

	if err := validateUser(user); err != nil {
		return models.User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return models.User{}, errors.New(
			"no fue posible proteger la contraseña",
		)
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	// Comprueba que el correo no exista antes de registrar.
	for _, existingUser := range us.users {
		if existingUser.Email == user.Email {
			return models.User{}, ErrEmailAlreadyRegistered
		}
	}

	user.ID = us.nextID
	us.nextID++
	user.Password = string(hashedPassword)

	us.users = append(us.users, user)

	return sanitizeUser(user), nil
}

// GetUsers devuelve una copia segura de todos los usuarios.
func (us *UserService) GetUsers() []models.User {
	us.mu.RLock()
	defer us.mu.RUnlock()

	usersCopy := make([]models.User, len(us.users))

	for index, user := range us.users {
		usersCopy[index] = sanitizeUser(user)
	}

	return usersCopy
}

// GetUserByID busca un usuario por su identificador.
func (us *UserService) GetUserByID(
	id int,
) (models.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, user := range us.users {
		if user.ID == id {
			return sanitizeUser(user), nil
		}
	}

	return models.User{}, ErrUserNotFound
}
