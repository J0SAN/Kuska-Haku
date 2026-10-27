package postgres

import (
	"context"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

// DistancesToStops implementa ports.RouteProgressRepository midiendo la distancia
// real sobre el trazado (LineString) desde la coordenada del bus hasta cada paradero.
func (r *routeRepository) DistancesToStops(ctx context.Context, routeID string, from domain.Coordinate) ([]domain.StopDistance, error) {
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
	// ST_MakePoint recibe (longitud, latitud)
	rows, err := r.db.QueryContext(ctx, q, routeID, from.Longitude, from.Latitude)
	if err != nil {
		if isInvalidUUID(err) {
			return []domain.StopDistance{}, nil
		}
		return nil, fmt.Errorf("postgres: error calculando distancias a paraderos: %w", err)
	}
	defer rows.Close()

	stops := []domain.StopDistance{}
	for rows.Next() {
		var sd domain.StopDistance
		var fp, fb, dist float64
		if err := rows.Scan(&sd.StopID, &sd.Name, &sd.StopOrder, &fp, &fb, &dist); err != nil {
			return nil, fmt.Errorf("postgres: error leyendo distancia a paradero: %w", err)
		}
		sd.Passed, sd.DistanceMeters = domain.ClassifyStop(fp, fb, dist)
		stops = append(stops, sd)
	}
	return stops, rows.Err()
}
