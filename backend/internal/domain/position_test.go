package domain

import (
	"testing"
	"time"
)

func TestPositionValidate(t *testing.T) {
	validPos := Position{
		BusID:     "bus-01",
		RouteID:   "ruta-01",
		Latitude:  -7.162,
		Longitude: -78.500,
		SpeedKmh:  25.0,
		Timestamp: time.Now(),
	}

	if err := validPos.Validate(); err != nil {
		t.Errorf("No se esperaba error en una posición válida, pero se obtuvo: %v", err)
	}

	cases := []struct {
		name string
		pos  Position
	}{
		{"sin bus_id", Position{RouteID: "ruta-01", Latitude: -7.16, Longitude: -78.50}},
		{"sin route_id", Position{BusID: "bus-01", Latitude: -7.16, Longitude: -78.50}},
		{"latitud fuera de rango", Position{BusID: "bus-01", RouteID: "ruta-01", Latitude: 150.0, Longitude: -78.50}},
		{"longitud fuera de rango", Position{BusID: "bus-01", RouteID: "ruta-01", Latitude: -7.16, Longitude: -300.0}},
		{"velocidad negativa", Position{BusID: "bus-01", RouteID: "ruta-01", Latitude: -7.16, Longitude: -78.50, SpeedKmh: -5.0}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.pos.Validate(); err == nil {
				t.Errorf("Se esperaba error de validación para: %s", c.name)
			}
		})
	}
}
