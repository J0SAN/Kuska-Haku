package postgres

import (
	"context"
	"database/sql"
	"errors"

	// Ajusta el nombre de tu módulo
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
		INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.FullName, user.Role).Scan(&user.ID)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role
		FROM users
		WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return user, nil
}
