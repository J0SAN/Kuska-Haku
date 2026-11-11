package usecases

import (
	"context"
	"fmt"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

// VehicleQRManager genera, rota y revoca los QR de los vehículos.
type VehicleQRManager struct {
	repo   ports.VehicleQRRepository
	tokens ports.QRTokenGenerator
}

func NewVehicleQRManager(repo ports.VehicleQRRepository, tokens ports.QRTokenGenerator) *VehicleQRManager {
	return &VehicleQRManager{repo: repo, tokens: tokens}
}

// Generate crea el primer QR activo del vehículo. Si ya existe uno, el
// repositorio devuelve domain.ErrQRAlreadyExists y el sticker no se toca.
func (m *VehicleQRManager) Generate(ctx context.Context, fleetID, createdBy string) (*domain.VehicleQR, error) {
	token, err := m.tokens.New()
	if err != nil {
		return nil, err
	}
	qr, err := m.repo.Create(ctx, fleetID, token, createdBy)
	if err != nil {
		return nil, fmt.Errorf("usecases: error generando QR: %w", err)
	}
	return qr, nil
}

// Rotate revoca el QR vigente y crea uno nuevo, en una sola transacción.
func (m *VehicleQRManager) Rotate(ctx context.Context, fleetID, createdBy string) (*domain.VehicleQR, error) {
	token, err := m.tokens.New()
	if err != nil {
		return nil, err
	}
	qr, err := m.repo.Rotate(ctx, fleetID, token, createdBy)
	if err != nil {
		return nil, fmt.Errorf("usecases: error rotando QR: %w", err)
	}
	return qr, nil
}

// Revoke desactiva el QR vigente del vehículo.
func (m *VehicleQRManager) Revoke(ctx context.Context, fleetID string) error {
	if err := m.repo.Revoke(ctx, fleetID); err != nil {
		return fmt.Errorf("usecases: error revocando QR: %w", err)
	}
	return nil
}
