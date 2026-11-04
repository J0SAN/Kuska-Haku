package http

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type stubRoutes struct {
	ports.RouteRepository
	route *domain.Route
}

func (s stubRoutes) GetByID(context.Context, string) (*domain.Route, error) { return s.route, nil }

type stubFleet struct {
	ports.FleetRepository
	byRoute  []domain.Fleet
	byDriver *domain.Fleet
}

func (s stubFleet) GetByRouteID(context.Context, string) ([]domain.Fleet, error) {
	return s.byRoute, nil
}
func (s stubFleet) GetByDriverID(context.Context, string) (*domain.Fleet, error) {
	return s.byDriver, nil
}

type stubPositionStore struct {
	ports.PositionStore
	pos *domain.Position
}

func (s stubPositionStore) GetLastPosition(context.Context, string) (*domain.Position, error) {
	return s.pos, nil
}

type stubRouteProgress struct {
	stops []domain.StopDistance
	off   float64
}

func (s stubRouteProgress) Progress(context.Context, string, domain.Coordinate) (*domain.BusProgress, error) {
	return &domain.BusProgress{Stops: s.stops, OffRouteMeters: s.off}, nil
}

func activeRoute() *domain.Route {
	return &domain.Route{ID: "r1", Name: "Ruta 01"}
}

func TestETA_CalculatesSecondsAndMarksPassedStops(t *testing.T) {
	pos := &domain.Position{
		BusID:     "bus-01",
		RouteID:   "r1",
		Latitude:  -7.16,
		Longitude: -78.51,
		SpeedKmh:  36,
		Timestamp: time.Now(),
	}
	h := NewPassengerController(
		stubRoutes{route: activeRoute()},
		stubFleet{byRoute: []domain.Fleet{{BusID: "bus-01", Status: "active"}}},
		stubPositionStore{pos: pos},
		stubRouteProgress{stops: []domain.StopDistance{
			{StopID: "s1", Name: "Uno", StopOrder: 1, Passed: true},
			{StopID: "s2", Name: "Dos", StopOrder: 2, DistanceMeters: 1000},
			{StopID: "s3", Name: "Tres", StopOrder: 3, DistanceMeters: 2000},
		}},
	)

	rec := httptest.NewRecorder()
	h.RouteItem(rec, httptest.NewRequest("GET", "/api/v1/routes/r1/eta", nil))

	if rec.Code != 200 {
		t.Fatalf("esperaba 200, obtuve %d: %s", rec.Code, rec.Body.String())
	}
	var resp etaResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Buses) != 1 || len(resp.Buses[0].Stops) != 3 {
		t.Fatalf("estructura inesperada: %+v", resp)
	}
	stops := resp.Buses[0].Stops
	if !stops[0].Passed || stops[0].ETASeconds != nil {
		t.Errorf("el paradero ya pasado no debe tener ETA: %+v", stops[0])
	}
	if stops[1].ETASeconds == nil || *stops[1].ETASeconds != 100 {
		t.Errorf("ETA paradero 2 = %v, esperaba 100", stops[1].ETASeconds)
	}
	if stops[2].ETASeconds == nil || *stops[2].ETASeconds != 200 {
		t.Errorf("ETA paradero 3 = %v, esperaba 200", stops[2].ETASeconds)
	}
}

func TestETA_BusWithoutPositionIsOmitted(t *testing.T) {
	h := NewPassengerController(
		stubRoutes{route: activeRoute()},
		stubFleet{byRoute: []domain.Fleet{{BusID: "bus-01", Status: "active"}}},
		stubPositionStore{pos: nil},
		stubRouteProgress{},
	)
	rec := httptest.NewRecorder()
	h.RouteItem(rec, httptest.NewRequest("GET", "/api/v1/routes/r1/eta", nil))

	var resp etaResponse
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if rec.Code != 200 || len(resp.Buses) != 0 {
		t.Fatalf("esperaba 200 con lista vacía, obtuve %d con %d buses", rec.Code, len(resp.Buses))
	}
}

func TestETA_NonExistentRouteReturns404(t *testing.T) {
	h := NewPassengerController(
		stubRoutes{route: nil},
		stubFleet{},
		stubPositionStore{},
		stubRouteProgress{},
	)
	rec := httptest.NewRecorder()
	h.RouteItem(rec, httptest.NewRequest("GET", "/api/v1/routes/r1/eta", nil))
	if rec.Code != 404 {
		t.Fatalf("esperaba 404, obtuve %d", rec.Code)
	}
}

func TestETA_OffRouteBusDoesNotCalculateETA(t *testing.T) {
	pos := &domain.Position{
		BusID:     "bus-01",
		RouteID:   "r1",
		Latitude:  -7.16,
		Longitude: -78.51,
		SpeedKmh:  30,
		Timestamp: time.Now(),
	}
	h := NewPassengerController(
		stubRoutes{route: activeRoute()},
		stubFleet{byRoute: []domain.Fleet{{BusID: "bus-01", Status: "active"}}},
		stubPositionStore{pos: pos},
		stubRouteProgress{
			off:   450,
			stops: []domain.StopDistance{{StopID: "s1", Name: "Uno", StopOrder: 1, DistanceMeters: 800}},
		},
	)
	rec := httptest.NewRecorder()
	h.RouteItem(rec, httptest.NewRequest("GET", "/api/v1/routes/r1/eta", nil))

	var resp etaResponse
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp.Buses) != 1 {
		t.Fatalf("el bus debe aparecer aunque esté fuera de ruta: %+v", resp)
	}
	b := resp.Buses[0]
	if !b.OffRoute || b.OffRouteMeters != 450 {
		t.Errorf("esperaba offRoute con desvío 450, obtuve %+v", b)
	}
	if b.Stops[0].ETASeconds != nil {
		t.Error("un bus fuera de ruta no debe tener ETA calculado")
	}
}
