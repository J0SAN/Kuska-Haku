package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type fleetRepository struct {
	db *sql.DB
}

func NewFleetRepository(db *sql.DB) *fleetRepository {
	return &fleetRepository{db: db}
}

func scanFleetRow(scanner interface{ Scan(dest ...any) error }) (domain.Fleet, error) {
	var f domain.Fleet
	var driverID sql.NullString
	err := scanner.Scan(&f.ID, &f.BusID, &f.RouteID, &driverID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return domain.Fleet{}, err
	}
	if driverID.Valid {
		val := driverID.String
		f.DriverID = &val
	}
	return f, nil
}

func (r *fleetRepository) Create(ctx context.Context, fleet *domain.Fleet) error {
	if err := fleet.Validate(); err != nil {
		return err
	}
	if fleet.Status == "" {
		fleet.Status = "active"
	}
	query := `
		INSERT INTO kuska.fleet (bus_id, route_id, driver_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, fleet.BusID, fleet.RouteID, fleet.DriverID, fleet.Status).
		Scan(&fleet.ID, &fleet.CreatedAt, &fleet.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrDuplicateKey
		}
		return fmt.Errorf("error insertando bus en bd: %v", err)
	}
	return nil
}

func (r *fleetRepository) GetAll(ctx context.Context) ([]domain.Fleet, error) {
	query := `
		SELECT id, bus_id, route_id, driver_id, status, created_at, updated_at
		FROM kuska.fleet
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error consultando flota: %v", err)
	}
	defer rows.Close()

	fleets := []domain.Fleet{}
	for rows.Next() {
		f, err := scanFleetRow(rows)
		if err != nil {
			return nil, fmt.Errorf("error escaneando fila de flota: %v", err)
		}
		fleets = append(fleets, f)
	}
	return fleets, rows.Err()
}

func (r *fleetRepository) GetByID(ctx context.Context, id string) (*domain.Fleet, error) {
	query := `
		SELECT id, bus_id, route_id, driver_id, status, created_at, updated_at
		FROM kuska.fleet
		WHERE id = $1
	`
	f, err := scanFleetRow(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || isInvalidUUID(err) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("error consultando bus por id: %v", err)
	}
	return &f, nil
}

func (r *fleetRepository) Update(ctx context.Context, fleet *domain.Fleet) error {
	if err := fleet.Validate(); err != nil {
		return err
	}
	if fleet.Status == "" {
		fleet.Status = "active"
	}
	query := `
		UPDATE kuska.fleet
		SET bus_id = $1, route_id = $2, driver_id = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, query, fleet.BusID, fleet.RouteID, fleet.DriverID, fleet.Status, fleet.ID).
		Scan(&fleet.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || isInvalidUUID(err) {
			return ErrRecordNotFound
		}
		if isUniqueViolation(err) {
			return ErrDuplicateKey
		}
		return fmt.Errorf("error actualizando bus: %v", err)
	}
	return nil
}

func (r *fleetRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM kuska.fleet WHERE id = $1`, id)
	if err != nil {
		if isInvalidUUID(err) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("error eliminando bus: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *fleetRepository) GetByRouteID(ctx context.Context, routeID string) ([]domain.Fleet, error) {
	query := `
		SELECT id, bus_id, route_id, driver_id, status, created_at, updated_at
		FROM kuska.fleet
		WHERE route_id = $1 AND status = 'active'
		ORDER BY bus_id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, routeID)
	if err != nil {
		if isInvalidUUID(err) {
			return []domain.Fleet{}, nil
		}
		return nil, fmt.Errorf("error consultando buses por ruta: %v", err)
	}
	defer rows.Close()

	fleets := []domain.Fleet{}
	for rows.Next() {
		f, err := scanFleetRow(rows)
		if err != nil {
			return nil, err
		}
		fleets = append(fleets, f)
	}
	return fleets, rows.Err()
}

func (r *fleetRepository) GetByDriverID(ctx context.Context, driverID string) (*domain.Fleet, error) {
	query := `
		SELECT id, bus_id, route_id, driver_id, status, created_at, updated_at
		FROM kuska.fleet
		WHERE driver_id = $1 AND status = 'active'
		LIMIT 1
	`
	f, err := scanFleetRow(r.db.QueryRowContext(ctx, query, driverID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || isInvalidUUID(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error consultando bus del conductor: %v", err)
	}
	return &f, nil
}
