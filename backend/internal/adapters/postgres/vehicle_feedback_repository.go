package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type vehicleFeedbackRepository struct {
	db *sql.DB
}

func NewVehicleFeedbackRepository(db *sql.DB) *vehicleFeedbackRepository {
	return &vehicleFeedbackRepository{db: db}
}

func (r *vehicleFeedbackRepository) Create(ctx context.Context, feedback *domain.VehicleFeedback) error {
	const q = `
		INSERT INTO kuska.vehicle_feedback
			(fleet_id, qr_code_id, bus_code, route_name, driver_id, rating, comment,
			 device_hash, service_date, bus_live)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id::text, status, created_at
	`
	err := r.db.QueryRowContext(ctx, q,
		feedback.FleetID, feedback.QRCodeID, feedback.BusCode, feedback.RouteName,
		feedback.DriverID, feedback.Rating, feedback.Comment, feedback.DeviceHash,
		feedback.ServiceDate, feedback.BusLive,
	).Scan(&feedback.ID, &feedback.Status, &feedback.CreatedAt)
	if err != nil {
		if isConstraintViolation(err, "uq_vehicle_feedback_device_day") {
			return domain.ErrFeedbackAlreadySubmitted
		}
		return fmt.Errorf("postgres: error guardando calificación: %w", err)
	}
	return nil
}

func (r *vehicleFeedbackRepository) List(ctx context.Context, filter domain.VehicleFeedbackFilter) ([]domain.VehicleFeedback, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 200 {
		limit = 200
	}

	const q = `
		SELECT vf.id::text, vf.fleet_id::text, vf.qr_code_id::text, vf.bus_code, vf.route_name,
		       vf.driver_id::text, COALESCE(u.full_name, ''), vf.rating, vf.comment,
		       vf.service_date::text, vf.bus_live, vf.status, vf.created_at
		FROM kuska.vehicle_feedback vf
		LEFT JOIN kuska.users u ON u.id = vf.driver_id
		WHERE ($1::text = '' OR vf.fleet_id::text = $1)
		  AND ($2::text = '' OR vf.status = $2)
		ORDER BY vf.created_at DESC
		LIMIT $3
	`
	rows, err := r.db.QueryContext(ctx, q, filter.FleetID, filter.Status, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres: error listando calificaciones: %w", err)
	}
	defer rows.Close()

	list := []domain.VehicleFeedback{}
	for rows.Next() {
		var item domain.VehicleFeedback
		var fleetID, qrCodeID, driverID sql.NullString
		if err := rows.Scan(
			&item.ID, &fleetID, &qrCodeID, &item.BusCode, &item.RouteName,
			&driverID, &item.DriverName, &item.Rating, &item.Comment,
			&item.ServiceDate, &item.BusLive, &item.Status, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: error leyendo calificación: %w", err)
		}
		if fleetID.Valid {
			value := fleetID.String
			item.FleetID = &value
		}
		if qrCodeID.Valid {
			value := qrCodeID.String
			item.QRCodeID = &value
		}
		if driverID.Valid {
			value := driverID.String
			item.DriverID = &value
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *vehicleFeedbackRepository) Summaries(ctx context.Context, fromDate string) ([]domain.VehicleRatingSummary, error) {
	const q = `
		SELECT vf.fleet_id::text, vf.bus_code, MAX(vf.route_name),
		       ROUND(AVG(vf.rating)::numeric, 2)::float8, COUNT(*)::int
		FROM kuska.vehicle_feedback vf
		WHERE vf.status <> 'hidden' AND vf.service_date >= $1::date
		GROUP BY vf.fleet_id, vf.bus_code
		ORDER BY AVG(vf.rating) ASC, vf.bus_code ASC
	`
	rows, err := r.db.QueryContext(ctx, q, fromDate)
	if err != nil {
		return nil, fmt.Errorf("postgres: error calculando promedios: %w", err)
	}
	defer rows.Close()

	summaries := []domain.VehicleRatingSummary{}
	for rows.Next() {
		var item domain.VehicleRatingSummary
		var fleetID sql.NullString
		if err := rows.Scan(&fleetID, &item.BusCode, &item.RouteName, &item.Average, &item.Total); err != nil {
			return nil, fmt.Errorf("postgres: error leyendo promedio: %w", err)
		}
		if fleetID.Valid {
			value := fleetID.String
			item.FleetID = &value
		}
		summaries = append(summaries, item)
	}
	return summaries, rows.Err()
}

func (r *vehicleFeedbackRepository) UpdateStatus(ctx context.Context, id, status string) error {
	if err := domain.ValidateFeedbackStatus(status); err != nil {
		return err
	}
	res, err := r.db.ExecContext(ctx,
		`UPDATE kuska.vehicle_feedback SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		if isInvalidUUID(err) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("postgres: error actualizando estado de calificación: %w", err)
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
