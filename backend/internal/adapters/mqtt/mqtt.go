package mqtt

import (
	"context"
	"encoding/json"
	"log"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"

	mqttclient "github.com/eclipse/paho.mqtt.golang"
)

type MqttAdapter struct {
	client    mqttclient.Client
	publisher ports.PositionPublisher
}

func NewMqttAdapter(brokerURL string, publisher ports.PositionPublisher) *MqttAdapter {
	opts := mqttclient.NewClientOptions().AddBroker(brokerURL).SetClientID("kuska-haku-backend")
	client := mqttclient.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error conectando a MQTT: %v", token.Error())
	}

	return &MqttAdapter{
		client:    client,
		publisher: publisher,
	}
}

// Start arranca la escucha de los mensajes del simulador de buses.
func (m *MqttAdapter) Start() {
	topic := "gps/positions"
	token := m.client.Subscribe(topic, 1, m.messageHandler)
	token.Wait()
	log.Printf("Suscrito a MQTT en el topic: %s", topic)
}

func (m *MqttAdapter) messageHandler(client mqttclient.Client, msg mqttclient.Message) {
	var pos domain.Position
	if err := json.Unmarshal(msg.Payload(), &pos); err != nil {
		log.Printf("Error decodificando mensaje MQTT: %v", err)
		return
	}

	if err := pos.Validate(); err != nil {
		log.Printf("Posición inválida recibida: %v", err)
		return
	}

	// Enviamos la posición válida al puerto (que la mandará a Redis)
	if err := m.publisher.Publish(context.Background(), pos); err != nil {
		log.Printf("Error publicando posición: %v", err)
	}
}
