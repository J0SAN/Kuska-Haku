package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	host := flag.String("host", "localhost:8080", "host:puerto del backend")
	routeID := flag.String("route", "ruta-01", "ID o código de la ruta a escuchar")
	token := flag.String("token", "", "JWT del usuario (obligatorio desde la Fase 4)")
	flag.Parse()

	q := url.Values{}
	q.Set("token", *token)
	u := url.URL{
		Scheme:   "ws",
		Host:     *host,
		Path:     "/ws/routes/" + *routeID,
		RawQuery: q.Encode(),
	}

	// No se imprime RawQuery en consola para no exponer el token JWT en los logs.
	log.Printf("wsclient: conectando a ws://%s%s", u.Host, u.Path)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("wsclient: fallo al conectar: %v", err)
	}
	defer conn.Close()

	log.Printf("wsclient: conectado. Esperando posiciones de la ruta %s...", *routeID)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("wsclient: conexión cerrada: %v", err)
			return
		}
		log.Printf("← %s", string(message))
	}
}
