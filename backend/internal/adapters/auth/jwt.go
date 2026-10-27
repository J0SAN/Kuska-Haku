package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secretKey: []byte(secret)}
}

type customClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Generate implementa ports.TokenService
func (j *JWTService) Generate(user domain.User) (string, error) {
	now := time.Now()
	claims := customClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// Parse implementa ports.TokenService devolviendo *ports.TokenClaims
func (j *JWTService) Parse(tokenString string) (*ports.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return j.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token inválido o expirado")
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return nil, errors.New("claims de token inválidos")
	}

	userID := claims.UserID
	if userID == "" {
		userID = claims.Subject
	}

	return &ports.TokenClaims{
		UserID: userID,
		Email:  claims.Email,
		Role:   domain.Role(claims.Role),
	}, nil
}

// Alias de retrocompatibilidad para tests de Fase 2/3
func (j *JWTService) GenerateToken(user *domain.User) (string, error) {
	return j.Generate(*user)
}

func (j *JWTService) ValidateToken(tokenString string) (*domain.User, error) {
	claims, err := j.Parse(tokenString)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  string(claims.Role),
	}, nil
}
