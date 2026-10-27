package domain

import (
	"errors"
	"strings"
	"time"
)

const minRoutePoints = 2

// Coordinate representa un punto geográfico individual en el trazado de una ruta o paradero.
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c Coordinate) Validate() error {
	if c.Latitude < -90 || c.Latitude > 90 {
		return errors.New("latitud fuera de rango válido (-90 a 90)")
	}
	if c.Longitude < -180 || c.Longitude > 180 {
		return errors.New("longitud fuera de rango válido (-180 a 180)")
	}
	return nil
}

// Route representa el recorrido operativo persistido como LineString en PostGIS.
type Route struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Path      []Coordinate `json:"path"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

func (r Route) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("el nombre de la ruta es requerido")
	}
	if len(r.Path) < minRoutePoints {
		return errors.New("el trazado debe tener al menos 2 puntos")
	}
	for _, c := range r.Path {
		if err := c.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Stop representa un paradero ordenado perteneciente a una ruta (necesario para el cálculo de ETA).
type Stop struct {
	ID        string     `json:"id"`
	RouteID   string     `json:"routeId"`
	Name      string     `json:"name"`
	Location  Coordinate `json:"location"`
	StopOrder int        `json:"stopOrder"`
	CreatedAt time.Time  `json:"createdAt"`
}

func (s Stop) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("el nombre del paradero es requerido")
	}
	if strings.TrimSpace(s.RouteID) == "" {
		return errors.New("el paradero debe pertenecer a una ruta")
	}
	if s.StopOrder < 1 {
		return errors.New("el orden del paradero debe ser mayor o igual a 1")
	}
	return s.Location.Validate()
}

// Fleet representa un vehículo (bus) físico o simulado registrado en la plataforma.
type Fleet struct {
	ID        string    `json:"id"`
	BusID     string    `json:"busId"`
	RouteID   string    `json:"routeId"`
	DriverID  *string   `json:"driverId,omitempty"`
	Status    string    `json:"status"` // ej: "active", "maintenance"
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (f Fleet) Validate() error {
	if strings.TrimSpace(f.BusID) == "" {
		return errors.New("el identificador del bus (busId) es requerido")
	}
	if strings.TrimSpace(f.RouteID) == "" {
		return errors.New("el bus debe estar asignado a una ruta (routeId)")
	}
	return nil
}
