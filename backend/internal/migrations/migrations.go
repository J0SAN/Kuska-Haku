// Package migrations aplica los archivos .sql de sql/ en orden, una sola
// vez cada uno, contra la base de datos. Con go:embed los archivos .sql
// viajan dentro del binario compilado.
package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
)

//go:embed sql/*.sql
var migrationFiles embed.FS

// Run aplica todas las migraciones pendientes en orden alfabético por
// nombre de archivo (001_, 002_, 003_, ...). Es idempotente: correrla de
// nuevo sin migraciones nuevas no hace nada.
func Run(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS kuska;
		CREATE TABLE IF NOT EXISTS kuska.schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`); err != nil {
		return fmt.Errorf("migrations: no se pudo crear la tabla de control: %w", err)
	}

	entries, err := fs.ReadDir(migrationFiles, "sql")
	if err != nil {
		return fmt.Errorf("migrations: no se pudo leer el directorio embebido: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		var alreadyApplied bool
		err := db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM kuska.schema_migrations WHERE version = $1)`,
			name,
		).Scan(&alreadyApplied)
		if err != nil {
			return fmt.Errorf("migrations: error verificando %s: %w", name, err)
		}
		if alreadyApplied {
			continue
		}

		content, err := migrationFiles.ReadFile("sql/" + name)
		if err != nil {
			return fmt.Errorf("migrations: error leyendo %s: %w", name, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("migrations: error iniciando transacción para %s: %w", name, err)
		}
		if _, err := tx.Exec(string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migrations: error aplicando %s: %w", name, err)
		}
		if _, err := tx.Exec(
			`INSERT INTO kuska.schema_migrations (version) VALUES ($1)`, name,
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migrations: error registrando %s: %w", name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("migrations: error confirmando %s: %w", name, err)
		}

		log.Printf("migrations: aplicada %s", name)
	}

	log.Println("Todas las migraciones se verificaron/ejecutaron correctamente")
	return nil
}
