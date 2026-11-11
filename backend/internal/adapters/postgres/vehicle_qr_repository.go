package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type vehicleQRRepository struct {
	db *sql.DB
}

func NewVehicleQRRepository(db *sql.DB) *vehicleQRRepository {
	return &vehicleQRRepository{db: db}
}

func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func (r *vehicleQRRepository) ListWithFleet(ctx context.Context) ([]domain.VehicleQRSummary, error) {
	const q = `
		SELECT f.id::text, f.bus_id, f.route_id::text, COALESCE(rt.name, f.route_id::text),
		       f.driver_id::text, COALESCE(u.full_name, ''),
		       q.id::text, q.token, q.created_at
		FROM kuska.fleet f
		LEFT JOIN kuska.routes rt ON rt.id = f.route_id
		LEFT JOIN kuska.users u ON u.id = f.driver_id
		LEFT JOIN kuska.vehicle_qr_codes q
		       ON q.fleet_id = f.id AND q.revoked_at IS NULL
		ORDER BY f.bus_id ASC
	`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("postgres: error listando QR de vehículos: %w", err)
	}
	defer rows.Close()

	summary := []domain.VehicleQRSummary{}
	for rows.Next() {
		var item domain.VehicleQRSummary
		var driverID, qrID, token sql.NullString
		var createdAt sql.NullTime
		if err := rows.Scan(
			&item.FleetID, &item.BusCode, &item.RouteID, &item.RouteName,
			&driverID, &item.DriverName, &qrID, &token, &createdAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: error leyendo QR de vehículo: %w", err)
		}
		if driverID.Valid {
			value := driverID.String
			item.DriverID = &value
		}
		if qrID.Valid {
			value := qrID.String
			item.QRCodeID = &value
		}
		if token.Valid {
			value := token.String
			item.Token = &value
		}
		if createdAt.Valid {
			value := createdAt.Time
			item.CreatedAt = &value
		}
		summary = append(summary, item)
	}
	return summary, rows.Err()
}

func (r *vehicleQRRepository) Create(ctx context.Context, fleetID, token, createdBy string) (*domain.VehicleQR, error) {
	if err := domain.ValidateQRToken(token); err != nil {
		return nil, err
	}

	const q = `
		INSERT INTO kuska.vehicle_qr_codes (fleet_id, token, created_by)
		VALUES ($1, $2, $3)
		RETURNING id::text, fleet_id::text, token, created_at
	`
	var qr domain.VehicleQR
	err := r.db.QueryRowContext(ctx, q, fleetID, token, nullString(createdBy)).
		Scan(&qr.ID, &qr.FleetID, &qr.Token, &qr.CreatedAt)
	if err != nil {
		if isConstraintViolation(err, "uq_vehicle_qr_active") {
			return nil, domain.ErrQRAlreadyExists
		}
		if isUniqueViolation(err) {
			return nil, domain.ErrQRAlreadyExists
		}
		if isForeignKeyViolation(err) {
			return nil, ErrRecordNotFound
		}
		if isInvalidUUID(err) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("postgres: error creando QR de vehículo: %w", err)
	}
	if createdBy != "" {
		qr.CreatedBy = &createdBy
	}
	return &qr, nil
}

func (r *vehicleQRRepository) Rotate(ctx context.Context, fleetID, token, createdBy string) (*domain.VehicleQR, error) {
	if err := domain.ValidateQRToken(token); err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("postgres: error iniciando transacción de rotación: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx,
		`UPDATE kuska.vehicle_qr_codes SET revoked_at = CURRENT_TIMESTAMP
		 WHERE fleet_id = $1 AND revoked_at IS NULL`, fleetID); err != nil {
		if isInvalidUUID(err) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("postgres: error revocando QR anterior: %w", err)
	}

	const q = `
		INSERT INTO kuska.vehicle_qr_codes (fleet_id, token, created_by)
		VALUES ($1, $2, $3)
		RETURNING id::text, fleet_id::text, token, created_at
	`
	var qr domain.VehicleQR
	if err := tx.QueryRowContext(ctx, q, fleetID, token, nullString(createdBy)).
		Scan(&qr.ID, &qr.FleetID, &qr.Token, &qr.CreatedAt); err != nil {
		if isForeignKeyViolation(err) || isInvalidUUID(err) {
			return nil, ErrRecordNotFound
		}
		if isUniqueViolation(err) {
			return nil, domain.ErrQRAlreadyExists
		}
		return nil, fmt.Errorf("postgres: error creando QR de rotación: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("postgres: error confirmando rotación de QR: %w", err)
	}
	if createdBy != "" {
		qr.CreatedBy = &createdBy
	}
	return &qr, nil
}

func (r *vehicleQRRepository) Revoke(ctx context.Context, fleetID string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE kuska.vehicle_qr_codes SET revoked_at = CURRENT_TIMESTAMP
		 WHERE fleet_id = $1 AND revoked_at IS NULL`, fleetID)
	if err != nil {
		if isInvalidUUID(err) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("postgres: error revocando QR de vehículo: %w", err)
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *vehicleQRRepository) ResolveActive(ctx context.Context, token string) (*domain.VehicleQRTarget, error) {
	const q = `
		SELECT q.id::text, f.id::text, f.bus_id, COALESCE(rt.name, f.route_id::text),
		       f.driver_id::text
		FROM kuska.vehicle_qr_codes q
		JOIN kuska.fleet f ON f.id = q.fleet_id
		LEFT JOIN kuska.routes rt ON rt.id = f.route_id
		WHERE q.token = $1 AND q.revoked_at IS NULL AND f.status = 'active'
	`
	var target domain.VehicleQRTarget
	var driverID sql.NullString
	err := r.db.QueryRowContext(ctx, q, token).
		Scan(&target.QRCodeID, &target.FleetID, &target.BusCode, &target.RouteName, &driverID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrQRNotAvailable
		}
		return nil, fmt.Errorf("postgres: error resolviendo QR: %w", err)
	}
	if driverID.Valid {
		value := driverID.String
		target.DriverID = &value
	}
	return &target, nil
}
