package http

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseAllowedOrigins(t *testing.T) {
	got := ParseAllowedOrigins(" https://panel.example.com/ , http://localhost:5173,,")
	want := []string{"https://panel.example.com", "http://localhost:5173"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("esperaba %v, obtuve %v", want, got)
	}
}

func corsRequest(method, origin string) *httptest.ResponseRecorder {
	handler := CORSMiddleware([]string{"https://panel.example.com"}, http.HandlerFunc(okHandler))
	req := httptest.NewRequest(method, "/api/v1/routes", nil)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestCORS_OrigenPermitido(t *testing.T) {
	rec := corsRequest(http.MethodGet, "https://panel.example.com")
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://panel.example.com" {
		t.Fatalf("esperaba el origen permitido, obtuve %q", got)
	}
	if rec.Header().Get("Vary") != "Origin" {
		t.Error("la respuesta debe incluir Vary: Origin")
	}
}

func TestCORS_OrigenNoPermitido(t *testing.T) {
	rec := corsRequest(http.MethodGet, "https://malicioso.example.com")
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("un origen no permitido no debe recibir la cabecera, obtuve %q", got)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("la petición sigue llegando al handler, esperaba 200 y obtuve %d", rec.Code)
	}
}

func TestCORS_PreflightNoLlegaAlHandler(t *testing.T) {
	handler := CORSMiddleware([]string{"https://panel.example.com"}, http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTeapot) },
	))
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/routes", nil)
	req.Header.Set("Origin", "https://panel.example.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("el preflight debe responder 200, obtuve %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("el preflight debe declarar los métodos permitidos")
	}
}
