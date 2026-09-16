package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	// Recuerda ajustar el nombre de tu módulo aquí también
	"kuska-haku-backend/internal/domain"
)

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secretKey: []byte(secret)}
}

type CustomClaims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *JWTService) GenerateToken(user *domain.User) (string, error) {
	claims := CustomClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Expiración en 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenStr string) (*domain.User, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return &domain.User{
			ID:   claims.UserID,
			Role: claims.Role,
		}, nil
	}

	return nil, errors.New("token inválido")
}
