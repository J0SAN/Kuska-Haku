package ports

import (
	"context"
	// Ajusta "kuska-haku" al nombre real del módulo que está en tu archivo go.mod
	"kuska-haku-backend/internal/domain"
)

// PositionPublisher define el contrato para enviar la posición a otros sistemas (como Redis o WebSocket).
type PositionPublisher interface {
	Publish(ctx context.Context, pos domain.Position) error
}
