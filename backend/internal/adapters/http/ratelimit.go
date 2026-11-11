package http

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimiter aplica una ventana fija por IP. Es intencionalmente simple y usa
// solo la stdlib: sirve para frenar el registro masivo y la fuerza bruta en los
// endpoints públicos de autenticación.
type RateLimiter struct {
	mu          sync.Mutex
	windows     map[string]*rateWindow
	limit       int
	window      time.Duration
	now         func() time.Time // inyectable para pruebas deterministas
	lastCleanup time.Time
}

type rateWindow struct {
	count   int
	resetAt time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		windows: make(map[string]*rateWindow),
		limit:   limit,
		window:  window,
		now:     time.Now,
	}
}

// Wrap protege next con el límite configurado. Responde 429 con Retry-After
// cuando se supera la cuota de la ventana.
func (l *RateLimiter) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retryAfter, allowed := l.allow(clientIP(r))
		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			writeError(w, http.StatusTooManyRequests, "Demasiados intentos. Inténtalo en unos minutos.")
			return
		}
		next(w, r)
	}
}

func (l *RateLimiter) allow(key string) (retryAfter int, allowed bool) {
	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanupLocked(now)

	current, ok := l.windows[key]
	if !ok || !now.Before(current.resetAt) {
		l.windows[key] = &rateWindow{count: 1, resetAt: now.Add(l.window)}
		return 0, true
	}

	current.count++
	if current.count <= l.limit {
		return 0, true
	}

	remaining := current.resetAt.Sub(now)
	seconds := int(math.Ceil(remaining.Seconds()))
	if seconds < 1 {
		seconds = 1
	}
	return seconds, false
}

// cleanupLocked purga las ventanas vencidas como máximo una vez por ventana.
func (l *RateLimiter) cleanupLocked(now time.Time) {
	if now.Sub(l.lastCleanup) < l.window {
		return
	}
	for key, window := range l.windows {
		if !now.Before(window.resetAt) {
			delete(l.windows, key)
		}
	}
	l.lastCleanup = now
}

// clientIP toma la IP del socket. Detrás de un proxy hay que definir los
// proxies de confianza antes de leer X-Forwarded-For (spec 002, §9.5).
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
