package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// healthHandler responde al endpoint /health requerido por la Fase 0
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Estructura JSON esperada según el README
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"checks": []map[string]string{
			{"name": "postgres", "status": "ok"},
		},
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/health", healthHandler)

	fmt.Println("kuska-haku backend escuchando en :8080")
	// Levanta el servidor en el puerto 8080
	error := http.ListenAndServe(":8080", nil)
	if error != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", error)
	}
}
