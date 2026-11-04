package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

// Progress implementa ports.RouteProgressRepository midiendo el desvío del bus respecto
// al trazado de la ruta y la distancia real sobre las calles hasta cada paradero.
func (r *routeRepository) Progress(ctx context.Context, routeID string, from domain.Coordinate) (*domain.BusProgress, error) {
	out := &domain.BusProgress{Stops: []domain.StopDistance{}}

	// 1. Medir desvío en metros reales (ST_MakePoint recibe longitud, latitud)
	err := r.db.QueryRowContext(ctx, `
		SELECT ST_Distance(path::geography, ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography)
		FROM kuska.routes
		WHERE id = $1
	`, routeID, from.Longitude, from.Latitude).Scan(&out.OffRouteMeters)
	if err != nil {
		if isInvalidUUID(err) || errors.Is(err, sql.ErrNoRows) {
			return out, nil
		}
		return nil, fmt.Errorf("postgres: error midiendo desvío de ruta: %w", err)
	}

	// 2. Medir distancia sobre el trazado a cada paradero
	const q = `
		WITH route_geom AS (
			SELECT path AS g FROM kuska.routes WHERE id = $1
		),
		bus_proj AS (
			SELECT ST_LineLocatePoint(route_geom.g, ST_SetSRID(ST_MakePoint($2, $3), 4326)) AS f
			FROM route_geom
		)
		SELECT s.id, s.name, s.stop_order,
		       ST_LineLocatePoint(route_geom.g, s.location) AS fp,
		       bus_proj.f AS fb,
		       ST_Length(ST_LineSubstring(
		           route_geom.g,
		           LEAST(bus_proj.f, ST_LineLocatePoint(route_geom.g, s.location)),
		           GREATEST(bus_proj.f, ST_LineLocatePoint(route_geom.g, s.location))
		       )::geography) AS dist
		FROM kuska.stops s
		CROSS JOIN route_geom
		CROSS JOIN bus_proj
		WHERE s.route_id = $1
		ORDER BY s.stop_order ASC
	`
	rows, err := r.db.QueryContext(ctx, q, routeID, from.Longitude, from.Latitude)
	if err != nil {
		return nil, fmt.Errorf("postgres: error calculando distancias a paraderos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sd domain.StopDistance
		var fp, fb, dist float64
		if err := rows.Scan(&sd.StopID, &sd.Name, &sd.StopOrder, &fp, &fb, &dist); err != nil {
			return nil, fmt.Errorf("postgres: error leyendo distancia a paradero: %w", err)
		}
		sd.Passed, sd.DistanceMeters = domain.ClassifyStop(fp, fb, dist)
		out.Stops = append(out.Stops, sd)
	}
	return out, rows.Err()
}
