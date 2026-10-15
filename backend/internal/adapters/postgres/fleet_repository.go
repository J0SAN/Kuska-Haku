package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type fleetRepository struct {
	db *sql.DB
}

// NewFleetRepository crea una nueva instancia del repositorio de flota
func NewFleetRepository(db *sql.DB) *fleetRepository {
	return &fleetRepository{db: db}
}

func (r *fleetRepository) Create(ctx context.Context, fleet *domain.Fleet) error {
	// Insertamos el registro. La base de datos generará el UUID y las fechas automáticamente
	query := `
		INSERT INTO fleet (bus_id, route_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, fleet.BusID, fleet.RouteID, fleet.Status).
		Scan(&fleet.ID, &fleet.CreatedAt, &fleet.UpdatedAt)

	if err != nil {
		// Aquí capturaremos errores como violaciones de la restricción UNIQUE del bus_id
		return fmt.Errorf("error insertando bus en bd: %v", err)
	}

	return nil
}

func (r *fleetRepository) GetAll(ctx context.Context) ([]domain.Fleet, error) {
	query := `
		SELECT id, bus_id, route_id, status, created_at, updated_at
		FROM fleet
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error consultando flota: %v", err)
	}
	defer rows.Close()

	var fleets []domain.Fleet
	for rows.Next() {
		var f domain.Fleet
		err := rows.Scan(&f.ID, &f.BusID, &f.RouteID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error escaneando fila de flota: %v", err)
		}
		fleets = append(fleets, f)
	}

	return fleets, nil
}
