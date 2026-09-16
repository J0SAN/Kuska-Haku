package domain

import (
	"errors"
	"regexp"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	Role         string
}

// NewUser crea una instancia validada, aplicando el rol 'passenger' por defecto si viene vacío.
func NewUser(email, passwordHash, fullName, role string) (*User, error) {
	if !isValidEmail(email) {
		return nil, errors.New("formato de correo inválido")
	}

	if role == "" {
		role = "passenger"
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
