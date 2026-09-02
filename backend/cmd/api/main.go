package main

import (
	"log"
	"net/http"

	// Usamos los adaptadores que programaste en el paso 3
	"kuska-haku-backend/internal/adapters/mqtt"
	"kuska-haku-backend/internal/adapters/redis"
	"kuska-haku-backend/internal/adapters/websocket"
)

func main() {
	// 1. Inicializamos Redis (nuestro PositionPublisher)
	redisPublisher := redis.NewRedisPublisher("localhost:6379")

	// 2. Inicializamos MQTT pasándole Redis para que inyecte los datos allí
	mqttAdapter := mqtt.NewMqttAdapter("tcp://localhost:1883", redisPublisher)

	// Arrancamos la escucha de MQTT en segundo plano
	go mqttAdapter.Start()

	// 3. Configuramos los endpoints HTTP
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`)) // Versión simple para la Fase 1
	})

	http.HandleFunc("/ws", websocket.HandleWS)

	// 4. Levantamos el servidor
	log.Println("kuska-haku backend escuchando en :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
