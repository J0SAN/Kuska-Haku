package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Se esperaba código %v, se obtuvo %v", http.StatusOK, status)
	}

	expected := `{"status":"ok"}`
	if rr.Body.String() != expected {
		t.Errorf("Se esperaba cuerpo %v, se obtuvo %v", expected, rr.Body.String())
	}
}
