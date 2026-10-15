package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"kuska-haku-backend/internal/domain"
)

type routeRepository struct {
	db *sql.DB
}

// NewRouteRepository crea una nueva instancia del repositorio de rutas
func NewRouteRepository(db *sql.DB) *routeRepository {
	return &routeRepository{db: db}
}

func (r *routeRepository) Create(ctx context.Context, route *domain.Route) error {
	if len(route.Path) < 2 {
		return fmt.Errorf("una ruta debe tener al menos 2 coordenadas")
	}

	// 1. Transformar []Coordinate de Go a formato WKT (Well-Known Text) de PostGIS
	// Regla de oro en GIS: El formato siempre es "Longitud Latitud" (X Y)
	var points []string
	for _, coord := range route.Path {
		points = append(points, fmt.Sprintf("%f %f", coord.Longitude, coord.Latitude))
	}
	linestring := fmt.Sprintf("LINESTRING(%s)", strings.Join(points, ", "))

	// 2. Insertar usando ST_GeomFromText con SRID 4326 (Estándar GPS WGS 84)
	query := `
		INSERT INTO routes (name, path)
		VALUES ($1, ST_GeomFromText($2, 4326))
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, route.Name, linestring).Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error insertando ruta en bd: %v", err)
	}

	return nil
}

func (r *routeRepository) GetAll(ctx context.Context) ([]domain.Route, error) {
	// 1. Usar ST_AsGeoJSON delega a Postgres la tarea de convertir el binario a JSON legible
	query := `
		SELECT id, name, ST_AsGeoJSON(path) as path_json, created_at, updated_at
		FROM routes
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error consultando rutas: %v", err)
	}
	defer rows.Close()

	var routes []domain.Route
	for rows.Next() {
		var rt domain.Route
		var pathJSON string

		err := rows.Scan(&rt.ID, &rt.Name, &pathJSON, &rt.CreatedAt, &rt.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error escaneando fila de ruta: %v", err)
		}

		// 2. Parsear el GeoJSON nativo de PostGIS a nuestra estructura
		var geo struct {
			Coordinates [][]float64 `json:"coordinates"`
		}
		if err := json.Unmarshal([]byte(pathJSON), &geo); err != nil {
			return nil, fmt.Errorf("error parseando geojson: %v", err)
		}

		// 3. Reconstruir el []Coordinate para el dominio
		for _, c := range geo.Coordinates {
			rt.Path = append(rt.Path, domain.Coordinate{
				Longitude: c[0],
				Latitude:  c[1],
			})
		}

		routes = append(routes, rt)
	}

	return routes, nil
}
