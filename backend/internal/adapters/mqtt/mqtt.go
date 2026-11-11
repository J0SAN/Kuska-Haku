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

var topics = map[string]byte{
	"gps/positions":          1,
	"kuska/+/bus/+/position": 1,
}

type MqttAdapter struct {
	client   mqttclient.Client
	ingestor ports.PositionIngestor
}

// NewMqttAdapter prepara el cliente sin conectarse todavía. La conexión se hace
// en Start y se reintenta en segundo plano, así un broker caído no impide que la
// API arranque: la ingesta HTTP del conductor sigue funcionando.
func NewMqttAdapter(brokerURL string, ingestor ports.PositionIngestor) *MqttAdapter {
	adapter := &MqttAdapter{ingestor: ingestor}

	opts := mqttclient.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("kuska-haku-backend").
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(2 * time.Second).
		// Con sesión limpia el broker olvida las suscripciones al reconectar,
		// por eso se vuelven a crear en cada conexión.
		SetOnConnectHandler(adapter.subscribe).
		SetConnectionLostHandler(func(_ mqttclient.Client, err error) {
			log.Printf("Conexión MQTT perdida, reintentando: %v", err)
		})
	adapter.client = mqttclient.NewClient(opts)
	return adapter
}

// Start inicia la conexión al broker. No bloquea: si el broker no responde,
// el cliente reintenta cada pocos segundos.
func (m *MqttAdapter) Start() {
	m.client.Connect()
}

func (m *MqttAdapter) subscribe(client mqttclient.Client) {
	token := client.SubscribeMultiple(topics, m.messageHandler)
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
