package http

import (
	"encoding/json"
	"log"
	"net/http"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

// AdminController agrupa los endpoints exclusivos para administradores (Fase 3)
type AdminController struct {
	routeRepo ports.RouteRepository
	fleetRepo ports.FleetRepository
}

// NewAdminController inyecta las dependencias de base de datos en el controlador
func NewAdminController(rr ports.RouteRepository, fr ports.FleetRepository) *AdminController {
	return &AdminController{
		routeRepo: rr,
		fleetRepo: fr,
	}
}

// CreateRoute recibe el polígono (trazado) desde el frontend y lo guarda
func (c *AdminController) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var route domain.Route

	// Parsear el JSON entrante
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		http.Error(w, `{"message": "JSON inválido"}`, http.StatusBadRequest)
		return
	}

	// Delegar la lógica compleja de base de datos (PostGIS) al repositorio
	if err := c.routeRepo.Create(r.Context(), &route); err != nil {
		log.Printf("ERROR CREANDO RUTA EN BD: %v", err) // <-- Imprime el error real
		http.Error(w, `{"message": "Error al crear la ruta"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(route)
}

// GetRoutes devuelve todas las rutas trazadas
func (c *AdminController) GetRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := c.routeRepo.GetAll(r.Context())
	if err != nil {
		log.Printf("ERROR OBTENIENDO RUTAS DE BD: %v", err) // <-- Imprime el error real
		http.Error(w, `{"message": "Error al obtener las rutas"}`, http.StatusInternalServerError)
		return
	}

	// Si no hay rutas, inicializar un arreglo vacío en lugar de null para que React no falle
	if routes == nil {
		routes = []domain.Route{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}

// CreateFleet registra un nuevo bus y lo asocia a una ruta
func (c *AdminController) CreateFleet(w http.ResponseWriter, r *http.Request) {
	var fleet domain.Fleet

	if err := json.NewDecoder(r.Body).Decode(&fleet); err != nil {
		http.Error(w, `{"message": "JSON inválido"}`, http.StatusBadRequest)
		return
	}

	// Por defecto, un bus nuevo entra en estado activo
	if fleet.Status == "" {
		fleet.Status = "active"
	}

	if err := c.fleetRepo.Create(r.Context(), &fleet); err != nil {
		http.Error(w, `{"message": "Error al registrar el vehículo. Verifique que el Bus ID no esté duplicado."}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fleet)
}

// GetFleet devuelve todos los buses registrados
func (c *AdminController) GetFleet(w http.ResponseWriter, r *http.Request) {
	fleet, err := c.fleetRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, `{"message": "Error al obtener la flota"}`, http.StatusInternalServerError)
		return
	}

	if fleet == nil {
		fleet = []domain.Fleet{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fleet)
}
