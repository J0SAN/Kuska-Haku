package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
	"kuska-haku-backend/internal/usecases"
)

// --- Fakes de los puertos ---

type stubVehicleQRRepo struct {
	ports.VehicleQRRepository
	summaries  []domain.VehicleQRSummary
	target     *domain.VehicleQRTarget
	resolveErr error
	createErr  error
	rotateErr  error
	revokeErr  error
	revoked    string
}

func (s *stubVehicleQRRepo) ListWithFleet(_ context.Context) ([]domain.VehicleQRSummary, error) {
	return s.summaries, nil
}

func (s *stubVehicleQRRepo) Create(_ context.Context, fleetID, token, _ string) (*domain.VehicleQR, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}
	return &domain.VehicleQR{ID: "qr-1", FleetID: fleetID, Token: token}, nil
}

func (s *stubVehicleQRRepo) Rotate(_ context.Context, fleetID, token, _ string) (*domain.VehicleQR, error) {
	if s.rotateErr != nil {
		return nil, s.rotateErr
	}
	return &domain.VehicleQR{ID: "qr-2", FleetID: fleetID, Token: token}, nil
}

func (s *stubVehicleQRRepo) Revoke(_ context.Context, fleetID string) error {
	s.revoked = fleetID
	return s.revokeErr
}

func (s *stubVehicleQRRepo) ResolveActive(_ context.Context, _ string) (*domain.VehicleQRTarget, error) {
	return s.target, s.resolveErr
}

type stubTokenGenerator struct{ token string }

func (s *stubTokenGenerator) New() (string, error) { return s.token, nil }

type stubFeedbackRepo struct {
	ports.VehicleFeedbackRepository
	list      []domain.VehicleFeedback
	summaries []domain.VehicleRatingSummary
	createErr error
	statusErr error
	status    string
	updatedID string
	fromDate  string
}

func (s *stubFeedbackRepo) Create(_ context.Context, feedback *domain.VehicleFeedback) error {
	if s.createErr != nil {
		return s.createErr
	}
	feedback.ID = "fb-1"
	feedback.Status = domain.FeedbackStatusNew
	return nil
}

func (s *stubFeedbackRepo) List(_ context.Context, _ domain.VehicleFeedbackFilter) ([]domain.VehicleFeedback, error) {
	return s.list, nil
}

func (s *stubFeedbackRepo) Summaries(_ context.Context, fromDate string) ([]domain.VehicleRatingSummary, error) {
	s.fromDate = fromDate
	return s.summaries, nil
}

func (s *stubFeedbackRepo) UpdateStatus(_ context.Context, id, status string) error {
	if s.statusErr != nil {
		return s.statusErr
	}
	s.updatedID, s.status = id, status
	return nil
}

type fbStubPositionStore struct {
	ports.PositionStore
	position *domain.Position
}

func (s *fbStubPositionStore) GetLastPosition(_ context.Context, _ string) (*domain.Position, error) {
	return s.position, nil
}

type stubTokenService struct {
	ports.TokenService
}

func (s *stubTokenService) Parse(token string) (*ports.TokenClaims, error) {
	switch token {
	case "admin-token":
		return &ports.TokenClaims{UserID: "admin-1", Role: domain.RoleAdmin}, nil
	case "driver-token":
		return &ports.TokenClaims{UserID: "driver-1", Role: domain.RoleDriver}, nil
	}
	return nil, errors.New("token inválido")
}

// --- Helpers de prueba ---

func sampleTarget() *domain.VehicleQRTarget {
	driver := "driver-1"
	return &domain.VehicleQRTarget{
		QRCodeID:  "qr-1",
		FleetID:   "fleet-1",
		BusCode:   "bus-01",
		RouteName: "R-01 Centro",
		DriverID:  &driver,
	}
}

func newPublicMux(qrRepo *stubVehicleQRRepo, feedbackRepo *stubFeedbackRepo, positions ports.PositionStore) *http.ServeMux {
	submitter := usecases.NewVehicleFeedbackSubmitter(feedbackRepo, qrRepo, positions)
	controller := NewPublicFeedbackController(submitter)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/public/vehicle-feedback/{token}", controller.Resolve)
	mux.HandleFunc("POST /api/v1/public/vehicle-feedback/{token}", controller.Submit)
	return mux
}

func doGET(mux *http.ServeMux, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func doPOST(mux *http.ServeMux, path, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, strings.NewReader(body)))
	return rec
}

const validBody = `{"rating":4,"comment":"Buen servicio","deviceId":"0123456789abcdef0123456789abcdef"}`

// --- Pruebas públicas ---

func TestPublicResolve_DevuelveSoloBusYRuta(t *testing.T) {
	mux := newPublicMux(&stubVehicleQRRepo{target: sampleTarget()}, &stubFeedbackRepo{}, nil)

	rec := doGET(mux, "/api/v1/public/vehicle-feedback/Abcdefghijklmnopqrstuv")
	if rec.Code != http.StatusOK {
		t.Fatalf("esperaba 200, obtuve %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `"busCode":"bus-01"`) || !strings.Contains(body, `"routeName":"R-01 Centro"`) {
		t.Fatalf("respuesta inesperada: %s", body)
	}
	for _, forbidden := range []string{"driver", "driver-1", "driverName"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("la respuesta pública no debe incluir %q: %s", forbidden, body)
		}
	}
}

func TestPublicResolve_TokenNoDisponibleDevuelve404(t *testing.T) {
	mux := newPublicMux(&stubVehicleQRRepo{resolveErr: domain.ErrQRNotAvailable}, &stubFeedbackRepo{}, nil)

	rec := doGET(mux, "/api/v1/public/vehicle-feedback/Abcdefghijklmnopqrstuv")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("esperaba 404, obtuve %d", rec.Code)
	}
}

func TestPublicResolve_TokenConFormatoInvalidoDevuelve404(t *testing.T) {
	mux := newPublicMux(&stubVehicleQRRepo{target: sampleTarget()}, &stubFeedbackRepo{}, nil)

	rec := doGET(mux, "/api/v1/public/vehicle-feedback/corto")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("esperaba 404, obtuve %d", rec.Code)
	}
}

func TestPublicSubmit_GuardaYDevuelve201(t *testing.T) {
	mux := newPublicMux(&stubVehicleQRRepo{target: sampleTarget()}, &stubFeedbackRepo{}, nil)

	rec := doPOST(mux, "/api/v1/public/vehicle-feedback/Abcdefghijklmnopqrstuv", validBody)
	if rec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuve %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"status":"received"`) {
		t.Fatalf("respuesta inesperada: %s", rec.Body.String())
	}
}

func TestPublicSubmit_MapeaErrores(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		qrRepo     *stubVehicleQRRepo
		feedback   *stubFeedbackRepo
		wantStatus int
	}{
		{
			name:       "rating inválido",
			body:       `{"rating":0,"deviceId":"0123456789abcdef0123456789abcdef"}`,
			qrRepo:     &stubVehicleQRRepo{target: sampleTarget()},
			feedback:   &stubFeedbackRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "comentario largo",
			body:       `{"rating":4,"comment":"` + strings.Repeat("a", 501) + `","deviceId":"0123456789abcdef0123456789abcdef"}`,
			qrRepo:     &stubVehicleQRRepo{target: sampleTarget()},
			feedback:   &stubFeedbackRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "deviceId inválido",
			body:       `{"rating":4,"deviceId":"abc"}`,
			qrRepo:     &stubVehicleQRRepo{target: sampleTarget()},
			feedback:   &stubFeedbackRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "json inválido",
			body:       `{"rating":`,
			qrRepo:     &stubVehicleQRRepo{target: sampleTarget()},
			feedback:   &stubFeedbackRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "qr no disponible",
			body:       validBody,
			qrRepo:     &stubVehicleQRRepo{resolveErr: domain.ErrQRNotAvailable},
			feedback:   &stubFeedbackRepo{},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "ya calificó hoy",
			body:       validBody,
			qrRepo:     &stubVehicleQRRepo{target: sampleTarget()},
			feedback:   &stubFeedbackRepo{createErr: domain.ErrFeedbackAlreadySubmitted},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mux := newPublicMux(tc.qrRepo, tc.feedback, nil)
			rec := doPOST(mux, "/api/v1/public/vehicle-feedback/Abcdefghijklmnopqrstuv", tc.body)
			if rec.Code != tc.wantStatus {
				t.Fatalf("esperaba %d, obtuve %d: %s", tc.wantStatus, rec.Code, rec.Body.String())
			}
			if strings.Contains(rec.Body.String(), "driver") {
				t.Fatalf("la respuesta no debe incluir datos del conductor: %s", rec.Body.String())
			}
		})
	}
}

func TestPublicSubmit_CuerpoMayorA4KBDevuelve400(t *testing.T) {
	mux := newPublicMux(&stubVehicleQRRepo{target: sampleTarget()}, &stubFeedbackRepo{}, nil)

	big := `{"rating":4,"comment":"` + strings.Repeat("a", 5000) + `","deviceId":"0123456789abcdef0123456789abcdef"}`
	rec := doPOST(mux, "/api/v1/public/vehicle-feedback/Abcdefghijklmnopqrstuv", big)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuve %d", rec.Code)
	}
}

func TestPublicSubmit_BusLiveSeMarcaConPosicion(t *testing.T) {
	feedback := &stubFeedbackRepo{}
	mux := newPublicMux(
		&stubVehicleQRRepo{target: sampleTarget()},
		feedback,
		&fbStubPositionStore{position: &domain.Position{BusID: "bus-01"}},
	)

	rec := doPOST(mux, "/api/v1/public/vehicle-feedback/Abcdefghijklmnopqrstuv", validBody)
	if rec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuve %d: %s", rec.Code, rec.Body.String())
	}
}

// --- Pruebas de administración ---

func newAdminMux(qrRepo *stubVehicleQRRepo, feedbackRepo *stubFeedbackRepo) *http.ServeMux {
	middleware := NewAuthMiddleware(&stubTokenService{})
	manager := usecases.NewVehicleQRManager(qrRepo, &stubTokenGenerator{token: "Abcdefghijklmnopqrstuv"})
	qrController := NewVehicleQRController(manager, qrRepo)
	feedbackController := NewVehicleFeedbackAdminController(feedbackRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/admin/vehicle-qr", middleware.RequireRole(qrController.List, domain.RoleAdmin))
	mux.HandleFunc("POST /api/v1/admin/vehicle-qr/{fleetId}", middleware.RequireRole(qrController.Generate, domain.RoleAdmin))
	mux.HandleFunc("POST /api/v1/admin/vehicle-qr/{fleetId}/rotate", middleware.RequireRole(qrController.Rotate, domain.RoleAdmin))
	mux.HandleFunc("DELETE /api/v1/admin/vehicle-qr/{fleetId}", middleware.RequireRole(qrController.Revoke, domain.RoleAdmin))
	mux.HandleFunc("GET /api/v1/admin/vehicle-feedback", middleware.RequireRole(feedbackController.List, domain.RoleAdmin))
	mux.HandleFunc("GET /api/v1/admin/vehicle-feedback/summary", middleware.RequireRole(feedbackController.Summary, domain.RoleAdmin))
	mux.HandleFunc("PUT /api/v1/admin/vehicle-feedback/{id}/status", middleware.RequireRole(feedbackController.UpdateStatus, domain.RoleAdmin))
	return mux
}

func TestAdminVehicleFeedback_SinTokenDevuelve401(t *testing.T) {
	mux := newAdminMux(&stubVehicleQRRepo{}, &stubFeedbackRepo{})

	for _, path := range []string{"/api/v1/admin/vehicle-qr", "/api/v1/admin/vehicle-feedback"} {
		rec := doGET(mux, path)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("%s: esperaba 401, obtuve %d", path, rec.Code)
		}
	}
}

func TestAdminVehicleFeedback_ConRolConductorDevuelve403(t *testing.T) {
	mux := newAdminMux(&stubVehicleQRRepo{}, &stubFeedbackRepo{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/vehicle-qr", nil)
	req.Header.Set("Authorization", "Bearer driver-token")
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("esperaba 403, obtuve %d", rec.Code)
	}
}

func TestAdminGenerate_Devuelve201YTokens(t *testing.T) {
	mux := newAdminMux(&stubVehicleQRRepo{}, &stubFeedbackRepo{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/vehicle-qr/fleet-1", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuve %d: %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "password") {
		t.Fatalf("respuesta inesperada: %s", rec.Body.String())
	}
}

func TestAdminGenerate_QRExistenteDevuelve409(t *testing.T) {
	mux := newAdminMux(&stubVehicleQRRepo{createErr: domain.ErrQRAlreadyExists}, &stubFeedbackRepo{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/vehicle-qr/fleet-1", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("esperaba 409, obtuve %d", rec.Code)
	}
}

func TestAdminRevoke_UnidadInexistenteDevuelve404(t *testing.T) {
	mux := newAdminMux(&stubVehicleQRRepo{revokeErr: domain.ErrQRNotAvailable}, &stubFeedbackRepo{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/vehicle-qr/fleet-1", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("esperaba 404, obtuve %d", rec.Code)
	}
}

func TestAdminFeedbackSummary_PeriodoInvalidoDevuelve400(t *testing.T) {
	mux := newAdminMux(&stubVehicleQRRepo{}, &stubFeedbackRepo{})

	for _, days := range []string{"0", "366", "abc"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/vehicle-feedback/summary?days="+days, nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("days=%s: esperaba 400, obtuve %d", days, rec.Code)
		}
	}
}

func TestAdminFeedbackSummary_PorDefectoUsa30Dias(t *testing.T) {
	feedback := &stubFeedbackRepo{}
	middleware := NewAuthMiddleware(&stubTokenService{})
	controller := NewVehicleFeedbackAdminController(feedback)
	controller.now = func() time.Time {
		return time.Date(2026, time.November, 11, 12, 0, 0, 0, time.UTC)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/admin/vehicle-feedback/summary", middleware.RequireRole(controller.Summary, domain.RoleAdmin))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/vehicle-feedback/summary", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperaba 200, obtuve %d", rec.Code)
	}
	if feedback.fromDate != "2026-10-13" {
		t.Fatalf("esperaba fromDate=2026-10-13 (30 días hasta 2026-11-11), obtuve %s", feedback.fromDate)
	}
}

func TestAdminUpdateStatus_EstadoInvalidoDevuelve400(t *testing.T) {
	mux := newAdminMux(&stubVehicleQRRepo{}, &stubFeedbackRepo{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/vehicle-feedback/fb-1/status", strings.NewReader(`{"status":"deleted"}`))
	req.Header.Set("Authorization", "Bearer admin-token")
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuve %d", rec.Code)
	}
}
