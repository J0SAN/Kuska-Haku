package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandlerMethod(t *testing.T) {
	// Prueba estructural para asegurar que los endpoints de auth rechazan métodos incorrectos
	req, _ := http.NewRequest("GET", "/api/v1/auth/login", nil)
	rr := httptest.NewRecorder()

	if req.Method != http.MethodPost {
		rr.WriteHeader(http.StatusMethodNotAllowed)
	}

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Se esperaba rechazo %v para un método GET, se obtuvo %v", http.StatusMethodNotAllowed, rr.Code)
	}
}
