package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"kuska-haku-backend/internal/domain"

	redisclient "github.com/redis/go-redis/v9"
)

// RedisPublisher implementa el puerto PositionPublisher usando Redis Pub/Sub.
type RedisPublisher struct {
	client *redisclient.Client
}

func NewRedisPublisher(addr string) *RedisPublisher {
	rdb := redisclient.NewClient(&redisclient.Options{
		Addr: addr,
	})
	return &RedisPublisher{client: rdb}
}

// Publish serializa la posición y la envía al canal 'gps_positions'.
func (r *RedisPublisher) Publish(ctx context.Context, pos domain.Position) error {
	data, err := json.Marshal(pos)
	if err != nil {
		return fmt.Errorf("error serializando posición: %w", err)
	}

	err = r.client.Publish(ctx, "gps_positions", data).Err()
	if err != nil {
		return fmt.Errorf("error publicando en redis: %w", err)
	}

	return nil
}
