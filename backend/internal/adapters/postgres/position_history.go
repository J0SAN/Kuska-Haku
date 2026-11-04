package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kuska-haku-backend/internal/domain"

	"github.com/lib/pq"
)

type PositionHistoryRepository struct {
	db *sql.DB
}

func NewPositionHistoryRepository(db *sql.DB) *PositionHistoryRepository {
	return &PositionHistoryRepository{db: db}
}

// SaveBatch inserta todo el lote en una sola sentencia SQL mediante unnest.
// ON CONFLICT (bus_id, ts) DO NOTHING garantiza que reenviar un lote ya guardado
// no duplique filas, devolviendo únicamente el conteo de filas nuevas insertadas.
func (r *PositionHistoryRepository) SaveBatch(ctx context.Context, ps []domain.Position) (int, error) {
	if len(ps) == 0 {
		return 0, nil
	}
	n := len(ps)
	buses := make([]string, n)
	routes := make([]string, n)
	timestamps := make([]string, n)
	lats := make([]float64, n)
	lngs := make([]float64, n)
	speeds := make([]float64, n)
	headings := make([]float64, n)

	for i, p := range ps {
		buses[i] = p.BusID
		routes[i] = p.RouteID
		timestamps[i] = p.Timestamp.UTC().Format(time.RFC3339Nano)
		lats[i] = p.Latitude
		lngs[i] = p.Longitude
		speeds[i] = p.SpeedKmh
		headings[i] = p.Heading
	}

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO kuska.gps_positions
			(bus_id, route_id, ts, latitude, longitude, speed_kmh, heading)
		SELECT * FROM unnest(
			$1::text[], $2::text[], $3::timestamptz[],
			$4::float8[], $5::float8[], $6::float8[], $7::float8[])
		ON CONFLICT (bus_id, ts) DO NOTHING
	`, pq.Array(buses), pq.Array(routes), pq.Array(timestamps),
		pq.Array(lats), pq.Array(lngs), pq.Array(speeds), pq.Array(headings))
	if err != nil {
		return 0, fmt.Errorf("postgres: error guardando historial de posiciones: %w", err)
	}
	inserted, _ := res.RowsAffected()
	return int(inserted), nil
}
