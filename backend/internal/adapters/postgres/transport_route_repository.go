package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type routeRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) *routeRepository {
	return &routeRepository{db: db}
}

func buildLineStringGeoJSON(path []domain.Coordinate) (string, error) {
	coords := make([][2]float64, len(path))
	for i, c := range path {
		coords[i] = [2]float64{c.Longitude, c.Latitude}
	}
	payload := map[string]any{
		"type":        "LineString",
		"coordinates": coords,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("postgres: error serializando geojson: %w", err)
	}
	return string(data), nil
}

func parseGeoJSONPath(pathJSON string) ([]domain.Coordinate, error) {
	var geo struct {
		Coordinates [][]float64 `json:"coordinates"`
	}
	if err := json.Unmarshal([]byte(pathJSON), &geo); err != nil {
		return nil, fmt.Errorf("error parseando geojson: %v", err)
	}
	coords := make([]domain.Coordinate, 0, len(geo.Coordinates))
	for _, c := range geo.Coordinates {
		coords = append(coords, domain.Coordinate{
			Longitude: c[0],
			Latitude:  c[1],
		})
	}
	return coords, nil
}

func (r *routeRepository) Create(ctx context.Context, route *domain.Route) error {
	if err := route.Validate(); err != nil {
		return err
	}

	geoJSON, err := buildLineStringGeoJSON(route.Path)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO kuska.routes (name, path)
		VALUES ($1, ST_SetSRID(ST_GeomFromGeoJSON($2), 4326))
		RETURNING id, created_at, updated_at
	`
	if err := r.db.QueryRowContext(ctx, query, route.Name, geoJSON).Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt); err != nil {
		if isUniqueViolation(err) {
			return ErrDuplicateKey
		}
		return fmt.Errorf("error insertando ruta en bd: %v", err)
	}
	return nil
}

func (r *routeRepository) GetAll(ctx context.Context) ([]domain.Route, error) {
	query := `
		SELECT id, name, ST_AsGeoJSON(path) as path_json, created_at, updated_at
		FROM kuska.routes
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error consultando rutas: %v", err)
	}
	defer rows.Close()

	routes := []domain.Route{}
	for rows.Next() {
		var rt domain.Route
		var pathJSON string

		if err := rows.Scan(&rt.ID, &rt.Name, &pathJSON, &rt.CreatedAt, &rt.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando fila de ruta: %v", err)
		}

		coords, err := parseGeoJSONPath(pathJSON)
		if err != nil {
			return nil, err
		}
		rt.Path = coords
		routes = append(routes, rt)
	}
	return routes, rows.Err()
}

func (r *routeRepository) GetByID(ctx context.Context, id string) (*domain.Route, error) {
	query := `
		SELECT id, name, ST_AsGeoJSON(path) as path_json, created_at, updated_at
		FROM kuska.routes
		WHERE id = $1
	`
	var rt domain.Route
	var pathJSON string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&rt.ID, &rt.Name, &pathJSON, &rt.CreatedAt, &rt.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || isInvalidUUID(err) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("error obteniendo ruta por id: %v", err)
	}

	coords, err := parseGeoJSONPath(pathJSON)
	if err != nil {
		return nil, err
	}
	rt.Path = coords
	return &rt, nil
}

func (r *routeRepository) Update(ctx context.Context, route *domain.Route) error {
	if err := route.Validate(); err != nil {
		return err
	}

	geoJSON, err := buildLineStringGeoJSON(route.Path)
	if err != nil {
		return err
	}

	query := `
		UPDATE kuska.routes
		SET name = $1, path = ST_SetSRID(ST_GeomFromGeoJSON($2), 4326), updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING updated_at
	`
	err = r.db.QueryRowContext(ctx, query, route.Name, geoJSON, route.ID).Scan(&route.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || isInvalidUUID(err) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("error actualizando ruta: %v", err)
	}
	return nil
}

func (r *routeRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM kuska.routes WHERE id = $1`, id)
	if err != nil {
		if isInvalidUUID(err) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("error eliminando ruta: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *routeRepository) CountBuses(ctx context.Context, routeID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM kuska.fleet WHERE route_id = $1`, routeID).Scan(&count)
	if isInvalidUUID(err) {
		return 0, nil
	}
	return count, err
}
