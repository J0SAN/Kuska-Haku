package mqtt

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"

	mqttclient "github.com/eclipse/paho.mqtt.golang"
)

type MqttAdapter struct {
	client   mqttclient.Client
	ingestor ports.PositionIngestor
}

func NewMqttAdapter(brokerURL string, ingestor ports.PositionIngestor) *MqttAdapter {
	opts := mqttclient.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("kuska-haku-backend").
		SetAutoReconnect(true).
		SetConnectRetryInterval(2 * time.Second)
	client := mqttclient.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error conectando a MQTT: %v", token.Error())
	}

	return &MqttAdapter{
		client:   client,
		ingestor: ingestor,
	}
}

// Start arranca la escucha de los mensajes del simulador de buses y rastreadores.
func (m *MqttAdapter) Start() {
	topics := map[string]byte{
		"gps/positions":          1,
		"kuska/+/bus/+/position": 1,
	}
	token := m.client.SubscribeMultiple(topics, m.messageHandler)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Error suscribiendo a MQTT: %v", token.Error())
		return
	}
	log.Printf("Suscrito a MQTT en los topics: gps/positions y kuska/+/bus/+/position")
}

func (m *MqttAdapter) messageHandler(_ mqttclient.Client, msg mqttclient.Message) {
	var pos domain.Position
	if err := json.Unmarshal(msg.Payload(), &pos); err != nil {
		log.Printf("Error decodificando mensaje MQTT: %v", err)
		return
	}
	if pos.Timestamp.IsZero() {
		pos.Timestamp = time.Now().UTC()
	}

	res, err := m.ingestor.Ingest(context.Background(), []domain.Position{pos})
	if err != nil {
		log.Printf("Error procesando posición MQTT de %s: %v", pos.BusID, err)
		return
	}
	if res.Discarded > 0 {
		log.Printf("Posición MQTT de %s descartada: %s", pos.BusID, res.FirstDiscardReason)
	}
}
