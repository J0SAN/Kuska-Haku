package domain

import (
	"errors"
	"regexp"
	"strings"
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
// DNI y TermsAcceptedAt son punteros porque las cuentas creadas antes de la
// migración 005 no tienen esos datos.
type User struct {
	ID              string
	Email           string
	PasswordHash    string
	FullName        string
	Role            string
	DNI             *string
	TermsAcceptedAt *time.Time
	CreatedAt       time.Time
}

var emailPattern = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,24}$`)

// NormalizeEmail recorta espacios y pasa el correo a minúsculas.
func NormalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ValidateEmail comprueba el formato de un correo ya normalizado.
func ValidateEmail(email string) error {
	if !emailPattern.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// NewUser crea una instancia validada, aplicando el rol 'passenger' por defecto si viene vacío.
func NewUser(email, passwordHash, fullName, role string) (*User, error) {
	email = NormalizeEmail(email)
	if err := ValidateEmail(email); err != nil {
		return nil, err
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
