package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type statsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *statsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) GetDashboard(ctx context.Context, today string, days int) (*domain.Dashboard, error) {
	d := &domain.Dashboard{
		Date:     today,
		Days:     days,
		Series:   []domain.DayStats{},
		TopBuses: []domain.BusStats{},
	}

	err := r.db.QueryRowContext(ctx, `
		SELECT
		  (SELECT COUNT(*) FROM kuska.fleet WHERE status = 'active'),
		  (SELECT COUNT(*) FROM kuska.routes),
		  (SELECT COUNT(*) FROM kuska.suggestions WHERE status = 'new'),
		  (SELECT COALESCE(SUM(passengers), 0) FROM kuska.income_records WHERE record_date = $1::date),
		  (SELECT COALESCE(SUM(amount), 0)::float8 FROM kuska.income_records WHERE record_date = $1::date)
	`, today).Scan(&d.ActiveBuses, &d.ActiveRoutes, &d.NewSuggestions, &d.PassengersToday, &d.IncomeToday)
	if err != nil {
		return nil, fmt.Errorf("postgres: error calculando KPIs del dashboard: %w", err)
	}

	serieRows, err := r.db.QueryContext(ctx, `
		SELECT to_char(d::date, 'YYYY-MM-DD'),
		       COALESCE(SUM(r.passengers), 0),
		       COALESCE(SUM(r.amount), 0)::float8
		FROM generate_series($1::date - ($2::int - 1), $1::date, interval '1 day') AS d
		LEFT JOIN kuska.income_records r ON r.record_date = d::date
		GROUP BY d
		ORDER BY d ASC
	`, today, days)
	if err != nil {
		return nil, fmt.Errorf("postgres: error calculando serie diaria: %w", err)
	}
	defer serieRows.Close()

	for serieRows.Next() {
		var s domain.DayStats
		if err := serieRows.Scan(&s.Date, &s.Passengers, &s.Amount); err != nil {
			return nil, fmt.Errorf("postgres: error leyendo serie diaria: %w", err)
		}
		d.Series = append(d.Series, s)
	}
	if err := serieRows.Err(); err != nil {
		return nil, err
	}

	topRows, err := r.db.QueryContext(ctx, `
		SELECT bus_code, SUM(passengers), SUM(amount)::float8
		FROM kuska.income_records
		WHERE record_date BETWEEN $1::date - ($2::int - 1) AND $1::date
		GROUP BY bus_code
		ORDER BY SUM(amount) DESC, bus_code ASC
		LIMIT 5
	`, today, days)
	if err != nil {
		return nil, fmt.Errorf("postgres: error calculando top de buses: %w", err)
	}
	defer topRows.Close()

	for topRows.Next() {
		var b domain.BusStats
		if err := topRows.Scan(&b.BusCode, &b.Passengers, &b.Amount); err != nil {
			return nil, fmt.Errorf("postgres: error leyendo top de buses: %w", err)
		}
		d.TopBuses = append(d.TopBuses, b)
	}
	return d, topRows.Err()
}
