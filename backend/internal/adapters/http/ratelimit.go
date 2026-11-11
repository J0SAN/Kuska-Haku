package http

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
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
	ipHeader    string // cabecera de confianza con la IP del cliente; vacía = socket
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

// WithClientIPHeader indica la cabecera que el proxy de la plataforma escribe con
// la IP real del cliente (en Render, CF-Connecting-IP, que Cloudflare sobrescribe
// en cada petición). Solo debe configurarse detrás de un proxy que la garantice:
// sin proxy, el cliente podría falsificarla.
func (l *RateLimiter) WithClientIPHeader(header string) *RateLimiter {
	l.ipHeader = header
	return l
}

// Wrap protege next con el límite configurado. Responde 429 con Retry-After
// cuando se supera la cuota de la ventana.
func (l *RateLimiter) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retryAfter, allowed := l.allow(clientIP(r, l.ipHeader))
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

// clientIP usa la cabecera de confianza si está configurada y contiene una IP
// válida; si no, la IP del socket. No se lee X-Forwarded-For porque el proxy de
// Render conserva lo que envía el cliente y solo agrega su propia dirección.
func clientIP(r *http.Request, trustedHeader string) string {
	if trustedHeader != "" {
		if ip := net.ParseIP(strings.TrimSpace(r.Header.Get(trustedHeader))); ip != nil {
			return ip.String()
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
