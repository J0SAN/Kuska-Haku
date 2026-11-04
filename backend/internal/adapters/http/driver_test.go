package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type stubIncomes struct {
	ports.IncomeRepository
	created *domain.IncomeRecord
}

func (s *stubIncomes) Create(_ context.Context, r *domain.IncomeRecord) error {
	r.ID = "rec-1"
	cp := *r
	s.created = &cp
	return nil
}

type fakeIngestor struct {
	got []domain.Position
	res domain.IngestResult
}

func (f *fakeIngestor) Ingest(_ context.Context, b []domain.Position) (domain.IngestResult, error) {
	f.got = b
	return f.res, nil
}

func withClaims(req *http.Request, role domain.Role) *http.Request {
	ctx := context.WithValue(req.Context(), claimsContextKey, &ports.TokenClaims{
		UserID: "user-1",
		Role:   role,
	})
	return req.WithContext(ctx)
}

func TestDriver_IncomeWithoutAssignedBusReturns409(t *testing.T) {
	h := NewDriverController(stubFleet{byDriver: nil}, stubRoutes{}, &stubIncomes{})
	body, _ := json.Marshal(incomeRequest{Passengers: 10, Amount: 25})
	req := withClaims(httptest.NewRequest("POST", "/api/v1/driver/incomes", bytes.NewReader(body)), domain.RoleDriver)
	rec := httptest.NewRecorder()
	h.Incomes(rec, req)
	if rec.Code != 409 {
		t.Fatalf("esperaba 409, obtuve %d", rec.Code)
	}
}

func TestDriver_IncomeUsesServerBusAndDate(t *testing.T) {
	repo := &stubIncomes{}
	h := NewDriverController(
		stubFleet{byDriver: &domain.Fleet{ID: "fleet-uuid", BusID: "bus-01", Status: "active"}},
		stubRoutes{},
		repo,
	)

	// Aunque un cliente malintencionado envíe busCode o recordDate falsos, el servidor los ignora.
	body := []byte(`{"passengers":12,"amount":30.5,"busCode":"otro-bus","recordDate":"2001-01-01"}`)
	req := withClaims(httptest.NewRequest("POST", "/api/v1/driver/incomes", bytes.NewReader(body)), domain.RoleDriver)
	rec := httptest.NewRecorder()
	h.Incomes(rec, req)

	if rec.Code != 201 {
		t.Fatalf("esperaba 201, obtuve %d: %s", rec.Code, rec.Body.String())
	}
	if repo.created.FleetID != "fleet-uuid" || repo.created.BusCode != "bus-01" {
		t.Errorf("el bus debe salir del conductor autenticado: %+v", repo.created)
	}
	if repo.created.RecordDate != domain.TodayLima() {
		t.Errorf("la fecha debe fijarla el servidor, obtuve %q", repo.created.RecordDate)
	}
	if repo.created.DriverID != "user-1" {
		t.Errorf("driverId = %q, esperaba user-1", repo.created.DriverID)
	}
}

func TestDriver_InvalidIncomeReturns400(t *testing.T) {
	h := NewDriverController(
		stubFleet{byDriver: &domain.Fleet{ID: "fleet-uuid", BusID: "bus-01", Status: "active"}},
		stubRoutes{},
		&stubIncomes{},
	)
	body, _ := json.Marshal(incomeRequest{Passengers: -3, Amount: 10})
	req := withClaims(httptest.NewRequest("POST", "/api/v1/driver/incomes", bytes.NewReader(body)), domain.RoleDriver)
	rec := httptest.NewRecorder()
	h.Incomes(rec, req)
	if rec.Code != 400 {
		t.Fatalf("esperaba 400, obtuve %d", rec.Code)
	}
}

func buildPositionsPayload(n int) []byte {
	lat := -7.16
	lng := -78.51
	ps := make([]positionPointRequest, n)
	for i := range ps {
		ps[i] = positionPointRequest{
			Latitude:  &lat,
			Longitude: &lng,
			SpeedKmh:  20,
			Heading:   90,
			Timestamp: time.Now(),
		}
	}
	b, _ := json.Marshal(positionsBatchRequest{Positions: ps})
	return b
}

func postDriverPositions(h *DriverController, body []byte) *httptest.ResponseRecorder {
	req := withClaims(httptest.NewRequest("POST", "/api/v1/driver/positions", bytes.NewReader(body)), domain.RoleDriver)
	rec := httptest.NewRecorder()
	h.Positions(rec, req)
	return rec
}

func driverWithAssignedBus(ing *fakeIngestor) *DriverController {
	return NewDriverController(
		stubFleet{byDriver: &domain.Fleet{ID: "fleet-uuid", BusID: "bus-01", RouteID: "r1", Status: "active"}},
		stubRoutes{route: activeRoute()},
		&stubIncomes{},
		ing,
	)
}

func TestPositions_IdentityComesFromServerNotClient(t *testing.T) {
	ing := &fakeIngestor{res: domain.IngestResult{Received: 1, Saved: 1}}
	// Un cliente malicioso intenta reportar como otro bus y otra ruta.
	body := []byte(`{"positions":[{"latitude":-7.16,"longitude":-78.51,"speedKmh":20,"heading":10,"timestamp":"2026-10-01T12:00:00Z","bus_id":"bus-99","route_id":"ruta-99"}]}`)
	rec := postDriverPositions(driverWithAssignedBus(ing), body)

	if rec.Code != 200 {
		t.Fatalf("esperaba 200, obtuve %d: %s", rec.Code, rec.Body.String())
	}
	if len(ing.got) != 1 || ing.got[0].BusID != "bus-01" || ing.got[0].RouteID != "r1" {
		t.Fatalf("la identidad debe salir del bus asignado en el servidor, obtuve %+v", ing.got)
	}
}

func TestPositions_BatchLimits(t *testing.T) {
	h := driverWithAssignedBus(&fakeIngestor{})
	if rec := postDriverPositions(h, buildPositionsPayload(0)); rec.Code != 400 {
		t.Errorf("lote vacío: esperaba 400, obtuve %d", rec.Code)
	}
	if rec := postDriverPositions(h, buildPositionsPayload(101)); rec.Code != 400 {
		t.Errorf("101 posiciones: esperaba 400, obtuve %d", rec.Code)
	}
	if rec := postDriverPositions(h, buildPositionsPayload(100)); rec.Code != 200 {
		t.Errorf("100 posiciones es el máximo permitido: obtuve %d", rec.Code)
	}
	if rec := postDriverPositions(h, []byte("no es json")); rec.Code != 400 {
		t.Errorf("JSON inválido: esperaba 400, obtuve %d", rec.Code)
	}
}

func TestPositions_WithoutAssignedBusReturns409(t *testing.T) {
	h := NewDriverController(stubFleet{byDriver: nil}, stubRoutes{}, &stubIncomes{}, &fakeIngestor{})
	if rec := postDriverPositions(h, buildPositionsPayload(1)); rec.Code != 409 {
		t.Fatalf("esperaba 409, obtuve %d", rec.Code)
	}
}
