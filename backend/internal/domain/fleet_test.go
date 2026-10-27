// backend/internal/domain/fleet_test.go
package domain

import (
	"testing"
	"time"
)

func validPath() []Coordinate {
	return []Coordinate{
		{Latitude: -7.1611, Longitude: -78.5127},
		{Latitude: -7.1580, Longitude: -78.5109},
	}
}

func TestFleetStruct(t *testing.T) {
	now := time.Now()
	fleet := Fleet{
		ID:        "uuid-test",
		BusID:     "bus-123",
		RouteID:   "route-xyz",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if fleet.BusID != "bus-123" {
		t.Errorf("Error en la asignación estructural: se esperaba 'bus-123', se obtuvo '%s'", fleet.BusID)
	}
	if err := fleet.Validate(); err != nil {
		t.Errorf("Fleet válido fue rechazado: %v", err)
	}
}

func TestRoute_Validate(t *testing.T) {
	valid := Route{Name: "Ruta Centro", Path: validPath()}
	if err := valid.Validate(); err != nil {
		t.Fatalf("ruta válida rechazada: %v", err)
	}

	cases := []struct {
		name  string
		route Route
	}{
		{"sin nombre", Route{Path: validPath()}},
		{"un solo punto", Route{Name: "Ruta X", Path: validPath()[:1]}},
		{"coordenada inválida", Route{Name: "Ruta X", Path: []Coordinate{{Latitude: -7.16, Longitude: -78.51}, {Latitude: 999, Longitude: -78.51}}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.route.Validate(); err == nil {
				t.Fatalf("se esperaba error para: %s", c.name)
			}
		})
	}
}

func TestStop_Validate(t *testing.T) {
	valid := Stop{
		RouteID:   "route-1",
		Name:      "Plaza de Armas",
		Location:  Coordinate{Latitude: -7.1611, Longitude: -78.5127},
		StopOrder: 1,
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("paradero válido rechazado: %v", err)
	}

	invalidOrder := valid
	invalidOrder.StopOrder = 0
	if err := invalidOrder.Validate(); err == nil {
		t.Fatal("un paradero con orden 0 debería ser inválido")
	}
}
