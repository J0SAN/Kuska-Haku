package usecases

import (
	"context"
	"fmt"
	"log"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

// VehicleFeedbackInput es la calificación que envía la página pública.
type VehicleFeedbackInput struct {
	Rating   int
	Comment  string
	DeviceID string
}

// VehicleFeedbackSubmitter resuelve el token del QR y guarda la calificación
// con los datos del vehículo y del conductor asignado en ese momento.
type VehicleFeedbackSubmitter struct {
	feedbackRepo ports.VehicleFeedbackRepository
	qrRepo       ports.VehicleQRRepository
	positions    ports.PositionStore
	now          func() time.Time
}

func NewVehicleFeedbackSubmitter(
	feedbackRepo ports.VehicleFeedbackRepository,
	qrRepo ports.VehicleQRRepository,
	positions ports.PositionStore,
) *VehicleFeedbackSubmitter {
	return &VehicleFeedbackSubmitter{
		feedbackRepo: feedbackRepo,
		qrRepo:       qrRepo,
		positions:    positions,
		now:          time.Now,
	}
}

// Resolve devuelve el bus y la ruta del token. Un token inexistente, revocado
// o de un bus inactivo se traduce a domain.ErrQRNotAvailable.
func (s *VehicleFeedbackSubmitter) Resolve(ctx context.Context, token string) (*domain.VehicleQRTarget, error) {
	if err := domain.ValidateQRToken(token); err != nil {
		return nil, err
	}
	target, err := s.qrRepo.ResolveActive(ctx, token)
	if err != nil {
		return nil, err
	}
	return target, nil
}

// Submit valida la entrada, resuelve el QR y guarda la calificación. Las
// validaciones ocurren antes de cualquier escritura.
func (s *VehicleFeedbackSubmitter) Submit(ctx context.Context, token string, input VehicleFeedbackInput) (*domain.VehicleFeedback, error) {
	if err := domain.ValidateRating(input.Rating); err != nil {
		return nil, err
	}
	comment, err := domain.NormalizeFeedbackComment(input.Comment)
	if err != nil {
		return nil, err
	}
	if err := domain.ValidateDeviceID(input.DeviceID); err != nil {
		return nil, err
	}

	target, err := s.Resolve(ctx, token)
	if err != nil {
		return nil, err
	}

	busLive := s.busIsLive(ctx, target.BusCode)

	now := s.now()
	feedback := &domain.VehicleFeedback{
		FleetID:     stringPointer(target.FleetID),
		QRCodeID:    stringPointer(target.QRCodeID),
		BusCode:     target.BusCode,
		RouteName:   target.RouteName,
		DriverID:    target.DriverID,
		Rating:      input.Rating,
		Comment:     comment,
		DeviceHash:  domain.HashDeviceID(input.DeviceID),
		ServiceDate: domain.DateLima(now),
		BusLive:     busLive,
	}
	if err := s.feedbackRepo.Create(ctx, feedback); err != nil {
		return nil, fmt.Errorf("usecases: error guardando calificación: %w", err)
	}
	return feedback, nil
}

// busIsLive consulta Redis solo como señal para el admin: si falla o no hay
// posición reciente, la calificación se guarda igual con BusLive=false.
func (s *VehicleFeedbackSubmitter) busIsLive(ctx context.Context, busCode string) bool {
	if s.positions == nil {
		return false
	}
	pos, err := s.positions.GetLastPosition(ctx, busCode)
	if err != nil {
		log.Printf("aviso: no se pudo consultar la posición del bus %s: %v", busCode, err)
		return false
	}
	return pos != nil
}

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
