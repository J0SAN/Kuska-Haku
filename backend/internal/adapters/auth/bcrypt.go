package auth

import "golang.org/x/crypto/bcrypt"

type BcryptService struct{}

func NewBcryptService() *BcryptService {
	return &BcryptService{}
}

// Implementación de ports.PasswordHasher
func (b *BcryptService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b *BcryptService) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Alias de retrocompatibilidad para tests existentes de Fase 2/3
func (b *BcryptService) HashPassword(password string) (string, error) {
	return b.Hash(password)
}

func (b *BcryptService) CheckPasswordHash(password, hash string) bool {
	return b.Compare(hash, password) == nil
}
