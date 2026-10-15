package domain

import (
	"time"
)

// Coordinate representa un punto geográfico individual en el trazado de una ruta
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Route representa el recorrido operativo.
// El campo Path abstraerá la complejidad geométrica (PostGIS) en un arreglo simple de coordenadas para el frontend.
type Route struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Path      []Coordinate `json:"path"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

// Fleet representa un vehículo (bus) físico o simulado registrado en la plataforma
type Fleet struct {
	ID        string    `json:"id"`
	BusID     string    `json:"busId"`
	RouteID   string    `json:"routeId"`
	Status    string    `json:"status"` // ej: "active", "maintenance"
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
