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

func TestPositionValidateForIngest(t *testing.T) {
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	base := Position{
		BusID:     "bus-01",
		RouteID:   "ruta-01",
		Latitude:  -7.16,
		Longitude: -78.51,
		SpeedKmh:  30,
		Heading:   90,
		Timestamp: now.Add(-5 * time.Second),
	}
	if err := base.ValidateForIngest(now); err != nil {
		t.Fatalf("punto válido rechazado: %v", err)
	}

	mut := func(f func(*Position)) Position { p := base; f(&p); return p }
	cases := map[string]Position{
		"sin timestamp":     mut(func(p *Position) { p.Timestamp = time.Time{} }),
		"del futuro":        mut(func(p *Position) { p.Timestamp = now.Add(5 * time.Minute) }),
		"de hace dos dias":  mut(func(p *Position) { p.Timestamp = now.Add(-48 * time.Hour) }),
		"velocidad absurda": mut(func(p *Position) { p.SpeedKmh = 400 }),
		"heading negativo":  mut(func(p *Position) { p.Heading = -1 }),
		"heading mayor 360": mut(func(p *Position) { p.Heading = 361 }),
		"latitud invalida":  mut(func(p *Position) { p.Latitude = 95 }),
	}
	for name, p := range cases {
		if err := p.ValidateForIngest(now); err == nil {
			t.Errorf("esperaba error para: %s", name)
		}
	}

	// Un reloj levemente adelantado (dentro de la tolerancia de 30s) sí se acepta.
	nearFuture := mut(func(p *Position) { p.Timestamp = now.Add(10 * time.Second) })
	if err := nearFuture.ValidateForIngest(now); err != nil {
		t.Errorf("un adelanto de 10 s debería tolerarse: %v", err)
	}
}

func TestPositionIsFresh(t *testing.T) {
	now := time.Now()
	if !(Position{Timestamp: now.Add(-30 * time.Second)}).IsFresh(now) {
		t.Error("30 s debería ser fresco")
	}
	if (Position{Timestamp: now.Add(-90 * time.Second)}).IsFresh(now) {
		t.Error("90 s no debería ser fresco")
	}
}
