package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"kuska-haku-backend/internal/domain"

	mqttclient "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	busID := flag.String("bus", "bus-01", "ID del bus")
	routeID := flag.String("route", "ruta-01", "ID de la ruta")
	flag.Parse()

	opts := mqttclient.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID(*busID + "-sim")
	client := mqttclient.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error conectando simulador: %v", token.Error())
	}

	fmt.Printf("Simulador iniciado para %s en %s. Presiona Ctrl+C para salir.\n", *busID, *routeID)

	lat, lon := -7.162, -78.500 // Coordenadas base iniciales
	for {
		pos := domain.Position{
			BusID:     *busID,
			RouteID:   *routeID,
			Latitude:  lat,
			Longitude: lon,
			Timestamp: time.Now(),
		}

		data, _ := json.Marshal(pos)
		client.Publish("gps/positions", 1, false, data)
		log.Printf("Emitiendo: Lat: %.4f, Lon: %.4f", lat, lon)

		lat += 0.001 // Simulamos movimiento
		time.Sleep(3 * time.Second)
	}
}
