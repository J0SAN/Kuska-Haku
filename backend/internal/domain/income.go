package domain

import (
	"errors"
	"math"
	"time"
	"unicode/utf8"
)

// Perú mantiene UTC-5 todo el año; usar FixedZone evita depender de tzdata en contenedores mínimos.
var limaZone = time.FixedZone("PET", -5*60*60)

// TodayLima devuelve la fecha actual en Perú con formato YYYY-MM-DD.
func TodayLima() string {
	return time.Now().In(limaZone).Format("2006-01-02")
}

const (
	maxPassengersPerRecord = 1000
	maxAmountPerRecord     = 100000.0
	maxNotesRunes          = 300
)

// IncomeRecord representa un registro manual de ingresos realizado por el conductor.
type IncomeRecord struct {
	ID         string    `json:"id"`
	FleetID    string    `json:"fleetId"`
	BusCode    string    `json:"busCode"`
	DriverID   string    `json:"driverId"`
	RecordDate string    `json:"recordDate"`
	Passengers int       `json:"passengers"`
	Amount     float64   `json:"amount"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (r IncomeRecord) Validate() error {
	if r.Passengers < 0 {
		return errors.New("el número de pasajeros no puede ser negativo")
	}
	if r.Passengers > maxPassengersPerRecord {
		return errors.New("el número de pasajeros por registro es demasiado alto")
	}
	if r.Amount < 0 {
		return errors.New("el monto no puede ser negativo")
	}
	if r.Amount > maxAmountPerRecord {
		return errors.New("el monto por registro es demasiado alto")
	}
	if math.Abs(r.Amount*100-math.Round(r.Amount*100)) > 1e-6 {
		return errors.New("el monto admite como máximo 2 decimales")
	}
	if r.Passengers == 0 && r.Amount == 0 {
		return errors.New("registra al menos un pasajero o un monto mayor a cero")
	}
	if utf8.RuneCountInString(r.Notes) > maxNotesRunes {
		return errors.New("las notas no pueden superar los 300 caracteres")
	}
	return nil
}
