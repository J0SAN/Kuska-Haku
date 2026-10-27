package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type incomeRepository struct {
	db *sql.DB
}

func NewIncomeRepository(db *sql.DB) *incomeRepository {
	return &incomeRepository{db: db}
}

func (r *incomeRepository) Create(ctx context.Context, record *domain.IncomeRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}

	const q = `
		INSERT INTO kuska.income_records
			(fleet_id, bus_code, driver_id, record_date, passengers, amount, notes)
		VALUES ($1, $2, $3, $4::date, $5, $6, $7)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, q,
		record.FleetID, record.BusCode, record.DriverID, record.RecordDate,
		record.Passengers, record.Amount, record.Notes,
	).Scan(&record.ID, &record.CreatedAt)
	if err != nil {
		return fmt.Errorf("postgres: error creando registro de ingreso: %w", err)
	}
	return nil
}

func (r *incomeRepository) GetByDriverAndDate(ctx context.Context, driverID, date string) ([]domain.IncomeRecord, error) {
	const q = `
		SELECT id, COALESCE(fleet_id::text, ''), bus_code, COALESCE(driver_id::text, ''),
		       to_char(record_date, 'YYYY-MM-DD'), passengers, amount::float8, notes, created_at
		FROM kuska.income_records
		WHERE driver_id = $1 AND record_date = $2::date
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, q, driverID, date)
	if err != nil {
		if isInvalidUUID(err) {
			return []domain.IncomeRecord{}, nil
		}
		return nil, fmt.Errorf("postgres: error listando ingresos del conductor: %w", err)
	}
	defer rows.Close()

	records := []domain.IncomeRecord{}
	for rows.Next() {
		var rec domain.IncomeRecord
		if err := rows.Scan(
			&rec.ID, &rec.FleetID, &rec.BusCode, &rec.DriverID,
			&rec.RecordDate, &rec.Passengers, &rec.Amount, &rec.Notes, &rec.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: error escaneando ingreso: %w", err)
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}
