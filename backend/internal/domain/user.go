package domain

import (
	"errors"
	"regexp"
	"time"
)

// Role define los permisos de un usuario autenticado dentro del sistema.
type Role string

const (
	RolePassenger Role = "passenger"
	RoleDriver    Role = "driver"
	RoleAdmin     Role = "admin"
)

func (r Role) Valid() bool {
	switch r {
	case RolePassenger, RoleDriver, RoleAdmin:
		return true
	}
	return false
}

// User es la entidad de dominio. PasswordHash nunca debe serializarse en respuestas HTTP.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	Role         string
	CreatedAt    time.Time
}

// NewUser crea una instancia validada, aplicando el rol 'passenger' por defecto si viene vacío.
func NewUser(email, passwordHash, fullName, role string) (*User, error) {
	if !isValidEmail(email) {
		return nil, errors.New("formato de correo inválido")
	}

	if role == "" {
		role = string(RolePassenger)
	}

	if !Role(role).Valid() {
		return nil, errors.New("rol inválido: debe ser passenger, driver o admin")
	}

	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Role:         role,
	}, nil
}

func isValidEmail(e string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(e)
}
