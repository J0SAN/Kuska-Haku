package ports

import (
	"context"
	// Ajusta "kuska-haku-backend" al nombre real del módulo que está en tu archivo go.mod
	"kuska-haku-backend/internal/domain"
)

// === FASE 2: Auth y Usuarios ===

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

// === FASE 1: Tiempo Real ===

// PositionPublisher define el contrato para enviar la posición a otros sistemas (como Redis o WebSocket).
type PositionPublisher interface {
	Publish(ctx context.Context, pos domain.Position) error
}

// === FASE 3: Flota y Trazado Espacial ===

// RouteRepository define el contrato para interactuar con las rutas en la base de datos PostGIS.
type RouteRepository interface {
	Create(ctx context.Context, route *domain.Route) error
	GetAll(ctx context.Context) ([]domain.Route, error)
}

// FleetRepository define el contrato para gestionar los vehículos vinculados a las rutas.
type FleetRepository interface {
	Create(ctx context.Context, fleet *domain.Fleet) error
	GetAll(ctx context.Context) ([]domain.Fleet, error)
}
