package websocket

import (
	"errors"
	"net/http/httptest"
	"testing"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type stubTokens struct{}

func (stubTokens) Generate(domain.User) (string, error) { return "", nil }
func (stubTokens) Parse(token string) (*ports.TokenClaims, error) {
	if token == "valid-token" {
		return &ports.TokenClaims{UserID: "u1", Role: domain.RolePassenger}, nil
	}
	return nil, errors.New("token inválido")
}

func TestServeRoute_WithoutTokenReturns401(t *testing.T) {
	g := NewGateway(nil, nil, stubTokens{})
	rec := httptest.NewRecorder()
	g.ServeRoute(rec, httptest.NewRequest("GET", "/ws/routes/ruta-01", nil))
	if rec.Code != 401 {
		t.Fatalf("esperaba 401 sin token, obtuve %d", rec.Code)
	}
}

func TestServeRoute_InvalidTokenReturns401(t *testing.T) {
	g := NewGateway(nil, nil, stubTokens{})
	rec := httptest.NewRecorder()
	g.ServeRoute(rec, httptest.NewRequest("GET", "/ws/routes/ruta-01?token=fake", nil))
	if rec.Code != 401 {
		t.Fatalf("esperaba 401 con token inválido, obtuve %d", rec.Code)
	}
}
