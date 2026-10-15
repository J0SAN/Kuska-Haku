package migrations

import (
	"database/sql"
	"log"
	"os"
)

func Run(db *sql.DB) error {
	// Definimos el orden exacto en el que deben ejecutarse las migraciones
	migrationFiles := []string{
		"internal/migrations/sql/001_create_users.sql",
		"internal/migrations/sql/002_create_fleet.sql",
	}

	// Recorremos y ejecutamos cada archivo
	for _, file := range migrationFiles {
		sqlFile, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error leyendo archivo de migración %s: %v", file, err)
			return err
		}

		_, err = db.Exec(string(sqlFile))
		if err != nil {
			log.Printf("Error ejecutando migración en BD %s: %v", file, err)
			return err
		}

		log.Printf("Migración aplicada: %s", file)
	}

	log.Println("Todas las migraciones se verificaron/ejecutaron correctamente")
	return nil
}
