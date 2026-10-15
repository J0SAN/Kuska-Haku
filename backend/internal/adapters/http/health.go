package http

import "net/http"

// HealthCheckHandler responde con el estado del servidor para herramientas de monitoreo
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
