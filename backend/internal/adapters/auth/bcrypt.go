package auth

import "golang.org/x/crypto/bcrypt"

type BcryptService struct{}

func NewBcryptService() *BcryptService {
	return &BcryptService{}
}

func (s *BcryptService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Costo 14 asegura una encriptación robusta
	return string(bytes), err
}

func (s *BcryptService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
