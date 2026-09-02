package domain

import (
	"testing"
	"time"
)

func TestPosition_Validate(t *testing.T) {
	// Prueba 1: Un escenario correcto (ejemplo: coordenadas por Cajamarca)
	validPos := Position{
		BusID:     "bus-01",
		RouteID:   "ruta-01",
		Latitude:  -7.162,
		Longitude: -78.500,
		Timestamp: time.Now(),
	}

	if err := validPos.Validate(); err != nil {
		t.Errorf("No se esperaba error en una posición válida, pero se obtuvo: %v", err)
	}

	// Prueba 2: Un escenario incorrecto
	invalidPos := Position{
		BusID:    "bus-01",
		RouteID:  "ruta-01",
		Latitude: 150.0, // Falla porque excede 90
	}

	if err := invalidPos.Validate(); err == nil {
		t.Errorf("Se esperaba un error por latitud inválida, pero pasó la validación")
	}
}
