package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader eleva una conexión HTTP normal a una conexión WebSocket en tiempo real.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // En producción se restringe a dominios específicos
	},
}

// HandleWS gestionará la conexión con la app móvil.
func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error estableciendo WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Nuevo cliente móvil conectado por WebSocket")

	// Aquí, más adelante, nos suscribiremos a Redis para enviarle los datos a esta conexión (conn.WriteJSON)
	// Por ahora, dejamos la conexión abierta para validar que la app móvil logra conectarse.
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Cliente desconectado")
			break
		}
	}
}
