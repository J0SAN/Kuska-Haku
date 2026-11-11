package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"kuska-haku-backend/internal/domain"

	redisclient "github.com/redis/go-redis/v9"
)

const lastPositionTTL = 2 * time.Minute

func positionKey(busID string) string     { return "position:" + busID }
func routeBusesKey(routeID string) string { return "route:" + routeID + ":buses" }
func routeChannel(routeID string) string  { return "route:" + routeID + ":positions" }

// RedisPublisher implementa PositionPublisher, PositionStore, PositionBroker y HealthPinger.
type RedisPublisher struct {
	client *redisclient.Client
}

func NewRedisPublisher(addr string) *RedisPublisher {
	rdb := redisclient.NewClient(&redisclient.Options{
		Addr: addr,
	})
	return &RedisPublisher{client: rdb}
}

// NewRedisPublisherFromURL acepta una URL redis:// o rediss:// (por ejemplo, la
// cadena de conexión interna de Render Key Value), incluida la autenticación.
func NewRedisPublisherFromURL(rawURL string) (*RedisPublisher, error) {
	opts, err := redisclient.ParseURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("redis: URL de conexión inválida: %w", err)
	}
	return &RedisPublisher{client: redisclient.NewClient(opts)}, nil
}

func (r *RedisPublisher) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisPublisher) Name() string {
	return "redis"
}

func (r *RedisPublisher) SaveLastPosition(ctx context.Context, pos domain.Position) error {
	data, err := json.Marshal(pos)
	if err != nil {
		return fmt.Errorf("redis: error serializando posición: %w", err)
	}
	if err := r.client.Set(ctx, positionKey(pos.BusID), data, lastPositionTTL).Err(); err != nil {
		return fmt.Errorf("redis: error guardando última posición: %w", err)
	}
	return r.client.SAdd(ctx, routeBusesKey(pos.RouteID), pos.BusID).Err()
}

func (r *RedisPublisher) GetLastPosition(ctx context.Context, busID string) (*domain.Position, error) {
	data, err := r.client.Get(ctx, positionKey(busID)).Bytes()
	if err == redisclient.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis: error obteniendo posición: %w", err)
	}
	var pos domain.Position
	if err := json.Unmarshal(data, &pos); err != nil {
		return nil, fmt.Errorf("redis: error deserializando posición: %w", err)
	}
	return &pos, nil
}

func (r *RedisPublisher) ListActiveBusIDs(ctx context.Context, routeID string) ([]string, error) {
	ids, err := r.client.SMembers(ctx, routeBusesKey(routeID)).Result()
	if err != nil {
		return nil, fmt.Errorf("redis: error listando buses activos de ruta %s: %w", routeID, err)
	}
	return ids, nil
}

// PublishToRoute envía la posición al canal específico de la ruta en Redis Pub/Sub.
func (r *RedisPublisher) PublishToRoute(ctx context.Context, routeID string, pos domain.Position) error {
	data, err := json.Marshal(pos)
	if err != nil {
		return fmt.Errorf("redis: error serializando posición: %w", err)
	}
	return r.client.Publish(ctx, routeChannel(routeID), data).Err()
}

// Subscribe abre una suscripción en Redis Pub/Sub exclusiva para la ruta solicitada.
func (r *RedisPublisher) Subscribe(ctx context.Context, routeID string) (<-chan domain.Position, func(), error) {
	sub := r.client.Subscribe(ctx, routeChannel(routeID))
	if _, err := sub.Receive(ctx); err != nil {
		sub.Close()
		return nil, nil, fmt.Errorf("redis: error suscribiendo a ruta %s: %w", routeID, err)
	}

	out := make(chan domain.Position)
	go func() {
		defer close(out)
		for msg := range sub.Channel() {
			var pos domain.Position
			if err := json.Unmarshal([]byte(msg.Payload), &pos); err != nil {
				continue
			}
			select {
			case out <- pos:
			case <-ctx.Done():
				return
			}
		}
	}()

	cleanup := func() { sub.Close() }
	return out, cleanup, nil
}

// Publish implementa PositionPublisher (llamado por el adaptador MQTT):
// guarda la última posición para el ETA y publica en el canal de la ruta.
func (r *RedisPublisher) Publish(ctx context.Context, pos domain.Position) error {
	if err := r.SaveLastPosition(ctx, pos); err != nil {
		return err
	}
	return r.PublishToRoute(ctx, pos.RouteID, pos)
}
