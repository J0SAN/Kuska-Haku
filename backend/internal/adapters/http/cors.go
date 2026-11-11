package http

import (
	"net/http"
	"strings"
)

// DefaultAllowedOrigins es el origen del panel en desarrollo (Vite).
const DefaultAllowedOrigins = "http://localhost:5173"

// ParseAllowedOrigins convierte una lista separada por comas (CORS_ALLOWED_ORIGINS)
// en orígenes normalizados, sin espacios ni barra final.
func ParseAllowedOrigins(raw string) []string {
	origins := []string{}
	for _, part := range strings.Split(raw, ",") {
		origin := strings.TrimRight(strings.TrimSpace(part), "/")
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

// CORSMiddleware solo responde con Access-Control-Allow-Origin cuando el origen
// de la petición está en la lista permitida. Un origen desconocido no recibe la
// cabecera y el navegador bloquea la respuesta.
func CORSMiddleware(allowedOrigins []string, next http.Handler) http.Handler {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		if origin := r.Header.Get("Origin"); allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Manejo de la petición preflight (OPTIONS) del navegador
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
