package redis

import (
	"context"
	"log"

	"kuska-haku-backend/internal/adapters/websocket"

	redisclient "github.com/redis/go-redis/v9"
)

type RedisSubscriber struct {
	client *redisclient.Client
}

func NewRedisSubscriber(addr string) *RedisSubscriber {
	rdb := redisclient.NewClient(&redisclient.Options{
		Addr: addr,
	})
	return &RedisSubscriber{client: rdb}
}

func (s *RedisSubscriber) Start() {
	ctx := context.Background()
	pubsub := s.client.Subscribe(ctx, "gps_positions")
	defer pubsub.Close()

	log.Println("Suscriptor de Redis escuchando en 'gps_positions'...")

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error recibiendo mensaje de Redis: %v", err)
			continue
		}

		// Enviar las coordenadas a todos los celulares conectados
		websocket.BroadcastMessage([]byte(msg.Payload))
	}
}
