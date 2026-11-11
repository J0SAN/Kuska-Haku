package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func requestFrom(ip string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
	req.RemoteAddr = ip + ":54321"
	return req
}

func TestRateLimiter_SextaPeticionDevuelve429(t *testing.T) {
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	limiter := NewRateLimiter(5, time.Hour)
	limiter.now = func() time.Time { return now }
	handler := limiter.Wrap(okHandler)

	for i := 0; i < 5; i++ {
		rec := httptest.NewRecorder()
		handler(rec, requestFrom("10.0.0.1"))
		if rec.Code != http.StatusOK {
			t.Fatalf("petición %d: esperaba 200, obtuve %d", i+1, rec.Code)
		}
	}

	rec := httptest.NewRecorder()
	handler(rec, requestFrom("10.0.0.1"))
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("esperaba 429, obtuve %d", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Error("la respuesta 429 debe incluir Retry-After")
	}
}

func TestRateLimiter_SeReiniciaAlPasarLaVentana(t *testing.T) {
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	limiter := NewRateLimiter(1, 10*time.Minute)
	limiter.now = func() time.Time { return now }
	handler := limiter.Wrap(okHandler)

	rec := httptest.NewRecorder()
	handler(rec, requestFrom("10.0.0.1"))
	if rec.Code != http.StatusOK {
		t.Fatalf("primera petición: esperaba 200, obtuve %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	handler(rec, requestFrom("10.0.0.1"))
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("segunda petición: esperaba 429, obtuve %d", rec.Code)
	}

	now = now.Add(10*time.Minute + time.Second)

	rec = httptest.NewRecorder()
	handler(rec, requestFrom("10.0.0.1"))
	if rec.Code != http.StatusOK {
		t.Fatalf("tras la ventana: esperaba 200, obtuve %d", rec.Code)
	}
}

func TestRateLimiter_IPsDistintasNoSeAfectan(t *testing.T) {
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	limiter := NewRateLimiter(1, time.Hour)
	limiter.now = func() time.Time { return now }
	handler := limiter.Wrap(okHandler)

	for _, ip := range []string{"10.0.0.1", "10.0.0.2"} {
		rec := httptest.NewRecorder()
		handler(rec, requestFrom(ip))
		if rec.Code != http.StatusOK {
			t.Fatalf("IP %s: esperaba 200, obtuve %d", ip, rec.Code)
		}
	}
}
