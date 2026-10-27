package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
