package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:8080/ws"
	log.Printf("Conectando al WebSocket del backend en: %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Fallo al conectar: %v", err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Error leyendo mensaje: %v", err)
		}
		log.Printf("Recibido del servidor: %s", string(message))
	}
}
