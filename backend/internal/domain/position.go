package domain

import (
	"errors"
	"time"
)

// Position representa la ubicación GPS de un vehículo de transporte público.
type Position struct {
	BusID     string    `json:"bus_id"`
	RouteID   string    `json:"route_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	SpeedKmh  float64   `json:"speed_kmh,omitempty"`
	Heading   float64   `json:"heading,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Validate asegura que los datos básicos de la posición tengan sentido físico y lógico.
func (p *Position) Validate() error {
	if p.BusID == "" {
		return errors.New("bus_id es requerido")
	}
	if p.RouteID == "" {
		return errors.New("route_id es requerido")
	}
	if p.Latitude < -90 || p.Latitude > 90 {
		return errors.New("latitud inválida: debe estar entre -90 y 90")
	}
	if p.Longitude < -180 || p.Longitude > 180 {
		return errors.New("longitud inválida: debe estar entre -180 y 180")
	}
	if p.SpeedKmh < 0 {
		return errors.New("speed_kmh no puede ser negativo")
	}
	return nil
}
