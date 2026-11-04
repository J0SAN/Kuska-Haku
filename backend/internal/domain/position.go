package domain

import (
	"errors"
	"time"
)

const (
	// liveFreshnessWindow: solo un punto con menos de 60 s de antigüedad se difunde en vivo.
	// Un lote atrasado (cuando el celular recupera señal) se guarda en el historial pero
	// NO se transmite por Redis/WebSocket como posición actual.
	liveFreshnessWindow = 60 * time.Second

	maxHistoryPointAge   = 24 * time.Hour
	futureClockTolerance = 30 * time.Second
	maxPlausibleSpeedKmh = 150.0
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
func (p Position) Validate() error {
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

// ValidateForIngest aplica reglas adicionales de protección para la ingesta de telemetría:
// descarta puntos sin timestamp, en el futuro, de hace más de 24 h o con velocidades/rumbos imposibles.
func (p Position) ValidateForIngest(now time.Time) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if p.Timestamp.IsZero() {
		return errors.New("timestamp es requerido")
	}
	if p.Timestamp.After(now.Add(futureClockTolerance)) {
		return errors.New("timestamp en el futuro")
	}
	if now.Sub(p.Timestamp) > maxHistoryPointAge {
		return errors.New("timestamp con más de 24 h de antigüedad")
	}
	if p.SpeedKmh > maxPlausibleSpeedKmh {
		return errors.New("velocidad no plausible (más de 150 km/h)")
	}
	if p.Heading < 0 || p.Heading > 360 {
		return errors.New("heading fuera de rango válido (0 a 360)")
	}
	return nil
}

// IsFresh indica si el punto es lo bastante reciente para difundirse en vivo a los pasajeros.
func (p Position) IsFresh(now time.Time) bool {
	return now.Sub(p.Timestamp) <= liveFreshnessWindow
}

// IngestResult resume el resultado del procesamiento de un lote de posiciones GPS.
type IngestResult struct {
	Received           int    `json:"received"`
	Saved              int    `json:"saved"`
	Duplicates         int    `json:"duplicates"`
	Discarded          int    `json:"discarded"`
	FirstDiscardReason string `json:"firstDiscardReason,omitempty"`
	Live               bool   `json:"live"`
}
