package ports

import (
	"context"
	// Ajusta "kuska-haku" al nombre real del módulo que está en tu archivo go.mod
	"kuska-haku-backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type AuthService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GenerateToken(user *domain.User) (string, error)
	ValidateToken(token string) (*domain.User, error)
}

// PositionPublisher define el contrato para enviar la posición a otros sistemas (como Redis o WebSocket).
type PositionPublisher interface {
	Publish(ctx context.Context, pos domain.Position) error
}
