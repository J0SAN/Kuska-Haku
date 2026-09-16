package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Variables globales para mantener un registro de los celulares conectados
var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

// BroadcastMessage permite que otros paquetes (como Redis) envíen datos al celular
func BroadcastMessage(msg []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error enviando a cliente, desconectando: %v", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error estableciendo WebSocket: %v", err)
		return
	}

	// Registrar al nuevo celular en nuestro directorio
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	log.Println("Nuevo cliente móvil conectado por WebSocket")

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
	}()

	// Mantener la conexión abierta
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Cliente desconectado")
			break
		}
	}
}
