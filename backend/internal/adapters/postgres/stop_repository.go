package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type stopRepository struct {
	db *sql.DB
}

func NewStopRepository(db *sql.DB) *stopRepository {
	return &stopRepository{db: db}
}

func (r *stopRepository) Create(ctx context.Context, stop *domain.Stop) error {
	if err := stop.Validate(); err != nil {
		return err
	}

	// ST_MakePoint recibe (longitud, latitud) en ese orden estricto.
	const q = `
		INSERT INTO kuska.stops (route_id, name, location, stop_order)
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, q,
		stop.RouteID, stop.Name, stop.Location.Longitude, stop.Location.Latitude, stop.StopOrder,
	).Scan(&stop.ID, &stop.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrDuplicateKey
		}
		if isForeignKeyViolation(err) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("postgres: error creando paradero: %w", err)
	}
	return nil
}

func (r *stopRepository) GetByRouteID(ctx context.Context, routeID string) ([]domain.Stop, error) {
	const q = `
		SELECT id, route_id, name,
		       ST_Y(location) AS lat,
		       ST_X(location) AS lng,
		       stop_order, created_at
		FROM kuska.stops
		WHERE route_id = $1
		ORDER BY stop_order ASC
	`
	rows, err := r.db.QueryContext(ctx, q, routeID)
	if err != nil {
		if isInvalidUUID(err) {
			return []domain.Stop{}, nil
		}
		return nil, fmt.Errorf("postgres: error listando paraderos: %w", err)
	}
	defer rows.Close()

	stops := []domain.Stop{}
	for rows.Next() {
		var s domain.Stop
		if err := rows.Scan(
			&s.ID, &s.RouteID, &s.Name,
			&s.Location.Latitude, &s.Location.Longitude,
			&s.StopOrder, &s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: error escaneando paradero: %w", err)
		}
		stops = append(stops, s)
	}
	return stops, rows.Err()
}

func (r *stopRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM kuska.stops WHERE id = $1`, id)
	if err != nil {
		if isInvalidUUID(err) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("postgres: error eliminando paradero: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
