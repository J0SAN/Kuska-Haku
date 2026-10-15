package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // driver de Postgres, registrado por side-effect
)

// DB envuelve la conexión y es el punto único desde el que se construirán
// los repositorios concretos en fases posteriores.
type DB struct {
	conn *sql.DB
}

// New abre la conexión contra Postgres usando una cadena DSN estándar,
// ej: postgres://user:pass@localhost:5432/kuska_haku?sslmode=disable
func New(dsn string) (*DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: error abriendo conexión: %w", err)
	}

	// Pool conservador para el MVP de una sola ruta; se ajustará con datos
	// reales de carga una vez en staging.
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

// Conn expone la conexión subyacente para quien necesite correr
// migraciones o construir repositorios concretos (ver user_repository.go).
func (d *DB) Conn() *sql.DB {
	return d.conn
}

// Close libera la conexión. Se llama desde main() al apagar el servicio.
func (d *DB) Close() error {
	return d.conn.Close()
}
