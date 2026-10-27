package postgres

import (
	"context"
	"database/sql"
	"errors"

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
		INSERT INTO kuska.users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.FullName, user.Role).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, created_at
		FROM kuska.users
		WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, created_at
		FROM kuska.users
		WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return user, nil
}
