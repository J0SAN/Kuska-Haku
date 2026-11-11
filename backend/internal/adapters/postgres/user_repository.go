package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kuska-haku-backend/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO kuska.users (email, password_hash, full_name, role, dni, terms_accepted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		user.Email, user.PasswordHash, user.FullName, user.Role, user.DNI, user.TermsAcceptedAt,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case isConstraintViolation(err, "uq_users_dni"):
			return domain.ErrDNIAlreadyRegistered
		case isConstraintViolation(err, "users_email_key"):
			return domain.ErrEmailAlreadyRegistered
		case isUniqueViolation(err):
			return ErrDuplicateKey
		}
		return fmt.Errorf("postgres: error creando usuario: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, dni, terms_accepted_at, created_at
		FROM kuska.users
		WHERE email = $1`

	return scanUser(r.db.QueryRowContext(ctx, query, email))
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, dni, terms_accepted_at, created_at
		FROM kuska.users
		WHERE id = $1`

	return scanUser(r.db.QueryRowContext(ctx, query, id))
}

func scanUser(row *sql.Row) (*domain.User, error) {
	user := &domain.User{}
	var dni sql.NullString
	var termsAcceptedAt sql.NullTime

	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role,
		&dni, &termsAcceptedAt, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	if dni.Valid {
		value := dni.String
		user.DNI = &value
	}
	if termsAcceptedAt.Valid {
		value := termsAcceptedAt.Time
		user.TermsAcceptedAt = &value
	}
	return user, nil
}

// ListByRole lista los usuarios de un rol específico sin incluir password_hash.
func (r *UserRepository) ListByRole(ctx context.Context, role domain.Role) ([]domain.User, error) {
	const query = `
		SELECT id, email, full_name, role, created_at
		FROM kuska.users
		WHERE role = $1
		ORDER BY full_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, string(role))
	if err != nil {
		return nil, fmt.Errorf("postgres: error listando usuarios por rol: %w", err)
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres: error escaneando usuario: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
