package domain

import (
	"testing"
	"time"
)

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
}
