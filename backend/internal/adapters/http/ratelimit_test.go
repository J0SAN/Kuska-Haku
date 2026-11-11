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

func TestClientIP_SinCabeceraConfiguradaIgnoraLasCabeceras(t *testing.T) {
	req := requestFrom("10.0.0.1")
	req.Header.Set("CF-Connecting-IP", "203.0.113.9")
	req.Header.Set("X-Forwarded-For", "198.51.100.7")
	if got := clientIP(req, "", 0); got != "10.0.0.1" {
		t.Fatalf("sin cabecera de confianza debe usar el socket, obtuve %s", got)
	}
}

func TestClientIP_UsaLaCabeceraDeConfianza(t *testing.T) {
	req := requestFrom("10.0.0.1")
	req.Header.Set("CF-Connecting-IP", " 203.0.113.9 ")
	if got := clientIP(req, "CF-Connecting-IP", 0); got != "203.0.113.9" {
		t.Fatalf("esperaba la IP de la cabecera, obtuve %s", got)
	}
}

func TestClientIP_CabeceraInvalidaVuelveAlSocket(t *testing.T) {
	req := requestFrom("10.0.0.1")
	req.Header.Set("CF-Connecting-IP", "no-es-una-ip")
	if got := clientIP(req, "CF-Connecting-IP", 0); got != "10.0.0.1" {
		t.Fatalf("con cabecera inválida debe usar el socket, obtuve %s", got)
	}
}

func TestClientIP_ConHopsCeroIgnoraXForwardedFor(t *testing.T) {
	req := requestFrom("10.0.0.1")
	req.Header.Set("X-Forwarded-For", "203.0.113.9")
	if got := clientIP(req, "", 0); got != "10.0.0.1" {
		t.Fatalf("con hops=0 debe usar el socket, obtuve %s", got)
	}
}

func TestClientIP_ConHopsUnoIgnoraElValorFalsificadoALaIzquierda(t *testing.T) {
	req := requestFrom("10.0.0.1")
	// El cliente puede escribir "forjada"; el proxy agrega la IP real a la derecha.
	req.Header.Set("X-Forwarded-For", "198.51.100.7, 203.0.113.9")
	if got := clientIP(req, "", 1); got != "203.0.113.9" {
		t.Fatalf("con hops=1 debe tomar la entrada de la derecha, obtuve %s", got)
	}
}

func TestClientIP_ConHopsDosTomaLaSegundaDesdeLaDerecha(t *testing.T) {
	req := requestFrom("10.0.0.1")
	req.Header.Set("X-Forwarded-For", "198.51.100.7, 203.0.113.9, 192.0.2.10")
	if got := clientIP(req, "", 2); got != "203.0.113.9" {
		t.Fatalf("con hops=2 debe tomar la segunda desde la derecha, obtuve %s", got)
	}
}

func TestClientIP_ConHopsInsuficientesVuelveAlSocket(t *testing.T) {
	req := requestFrom("10.0.0.1")
	req.Header.Set("X-Forwarded-For", "203.0.113.9")
	if got := clientIP(req, "", 3); got != "10.0.0.1" {
		t.Fatalf("con menos entradas que hops debe usar el socket, obtuve %s", got)
	}
}

func TestRateLimiter_ConHopsCadaClienteTieneSuCuota(t *testing.T) {
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	limiter := NewRateLimiter(1, time.Hour).WithTrustedProxyHops(1)
	limiter.now = func() time.Time { return now }
	handler := limiter.Wrap(okHandler)

	for _, ip := range []string{"203.0.113.1", "203.0.113.2"} {
		req := requestFrom("10.0.0.1")
		req.Header.Set("X-Forwarded-For", ip)
		rec := httptest.NewRecorder()
		handler(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("cliente %s: esperaba 200, obtuve %d", ip, rec.Code)
		}
	}
}

func TestRateLimiter_DetrasDelProxyCadaClienteTieneSuCuota(t *testing.T) {
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	limiter := NewRateLimiter(1, time.Hour).WithClientIPHeader("CF-Connecting-IP")
	limiter.now = func() time.Time { return now }
	handler := limiter.Wrap(okHandler)

	// Mismo socket (el proxy) y clientes distintos
	for _, ip := range []string{"203.0.113.1", "203.0.113.2"} {
		req := requestFrom("10.0.0.1")
		req.Header.Set("CF-Connecting-IP", ip)
		rec := httptest.NewRecorder()
		handler(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("cliente %s: esperaba 200, obtuve %d", ip, rec.Code)
		}
	}
}
