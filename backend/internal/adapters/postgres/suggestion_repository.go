package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type suggestionRepository struct {
	db *sql.DB
}

func NewSuggestionRepository(db *sql.DB) *suggestionRepository {
	return &suggestionRepository{db: db}
}

func (r *suggestionRepository) Create(ctx context.Context, s *domain.Suggestion) error {
	if err := domain.ValidateSuggestionMessage(s.Message); err != nil {
		return err
	}

	const q = `
		INSERT INTO kuska.suggestions (user_id, message)
		VALUES ($1, $2)
		RETURNING id, status, created_at
	`
	err := r.db.QueryRowContext(ctx, q, s.UserID, s.Message).
		Scan(&s.ID, &s.Status, &s.CreatedAt)
	if err != nil {
		return fmt.Errorf("postgres: error creando sugerencia: %w", err)
	}
	return nil
}

func (r *suggestionRepository) GetAll(ctx context.Context, onlyNew bool) ([]domain.Suggestion, error) {
	const q = `
		SELECT s.id, COALESCE(s.user_id::text, ''), COALESCE(u.full_name, 'Usuario eliminado'),
		       s.message, s.status, s.created_at
		FROM kuska.suggestions s
		LEFT JOIN kuska.users u ON u.id = s.user_id
		WHERE ($1::boolean = false OR s.status = 'new')
		ORDER BY s.created_at DESC
		LIMIT 200
	`
	rows, err := r.db.QueryContext(ctx, q, onlyNew)
	if err != nil {
		return nil, fmt.Errorf("postgres: error listando sugerencias: %w", err)
	}
	defer rows.Close()

	suggestions := []domain.Suggestion{}
	for rows.Next() {
		var s domain.Suggestion
		if err := rows.Scan(&s.ID, &s.UserID, &s.UserName, &s.Message, &s.Status, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres: error leyendo sugerencia: %w", err)
		}
		suggestions = append(suggestions, s)
	}
	return suggestions, rows.Err()
}

func (r *suggestionRepository) MarkRead(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE kuska.suggestions SET status = 'read' WHERE id = $1`, id)
	if err != nil {
		if isInvalidUUID(err) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("postgres: error marcando sugerencia como leída: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrRecordNotFound
	}
	return nil
}
