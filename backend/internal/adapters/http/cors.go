package http

import "net/http"

// CORSMiddleware envuelve los handlers para permitir peticiones desde el frontend web
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Autorizar explícitamente el origen de nuestro entorno de desarrollo Vite
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

		// Autorizar métodos REST
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Autorizar cabeceras de seguridad necesarias para JWT
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Manejo seguro de la petición Preflight (OPTIONS) del navegador
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
