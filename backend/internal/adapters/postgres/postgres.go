package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// Errores estandarizados de la capa de persistencia
var (
	ErrRecordNotFound        = errors.New("registro no encontrado")
	ErrDuplicateKey          = errors.New("registro duplicado (violación de unicidad)")
	ErrDriverAlreadyAssigned = errors.New("ese conductor ya tiene otro bus activo asignado")
	ErrForeignKeyViolation   = errors.New("la ruta o el conductor indicados no existen")
	ErrRouteHasBuses         = errors.New("la ruta tiene buses asignados")
)

// isUniqueViolation detecta el código 23505 (unique_violation) de PostgreSQL.
func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

// isConstraintViolation detecta si una violación 23505 corresponde a un constraint específico.
func isConstraintViolation(err error, constraint string) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" && pqErr.Constraint == constraint
	}
	return err != nil && strings.Contains(err.Error(), constraint)
}

// isForeignKeyViolation detecta el código 23503 (foreign_key_violation) de PostgreSQL.
func isForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23503"
	}
	return err != nil && strings.Contains(err.Error(), "violates foreign key constraint")
}

// isInvalidUUID detecta cuando se consulta con un string que no es un UUID válido,
// permitiendo tratarlo como un 404 (ErrRecordNotFound) en vez de un 500.
func isInvalidUUID(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "22P02"
	}
	return err != nil && strings.Contains(err.Error(), "invalid input syntax for type uuid")
}

// DB envuelve la conexión y es el punto único desde el que se construirán
// los repositorios concretos.
type DB struct {
	conn *sql.DB
}

func New(dsn string) (*DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: error abriendo conexión: %w", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)

	return &DB{conn: conn}, nil
}

// Ping implementa ports.HealthPinger.
func (d *DB) Ping(ctx context.Context) error {
	return d.conn.PingContext(ctx)
}

// Name implementa ports.HealthPinger.
func (d *DB) Name() string {
	return "postgres"
}

func (d *DB) Conn() *sql.DB {
	return d.conn
}

func (d *DB) Close() error {
	return d.conn.Close()
}
