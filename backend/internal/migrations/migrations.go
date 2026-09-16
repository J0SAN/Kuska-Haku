package migrations

import (
	"database/sql"
	"log"
	"os"
)

func Run(db *sql.DB) error {
	// Lee el archivo asumiendo que ejecutas el servidor desde la carpeta "backend"
	sqlFile, err := os.ReadFile("internal/migrations/sql/001_create_users.sql")
	if err != nil {
		log.Printf("Error leyendo archivo de migración: %v", err)
		return err
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		log.Printf("Error ejecutando migración en BD: %v", err)
		return err
	}

	log.Println("Migraciones ejecutadas correctamente")
	return nil
}
