package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"kuska-haku-backend/internal/adapters/postgres"
	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, errorResponse{Error: message, Message: message})
}

func idFromPath(path, prefix string) (id, rest string) {
	trimmed := strings.TrimPrefix(path, prefix)
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		return "", ""
	}
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

// AdminController agrupa el CRUD administrativo de rutas, paraderos y flota.
type AdminController struct {
	routeRepo ports.RouteRepository
	fleetRepo ports.FleetRepository
	stopRepo  ports.StopRepository
}

func NewAdminController(rr ports.RouteRepository, fr ports.FleetRepository, sr ...ports.StopRepository) *AdminController {
	var stopRepo ports.StopRepository
	if len(sr) > 0 {
		stopRepo = sr[0]
	}
	return &AdminController{
		routeRepo: rr,
		fleetRepo: fr,
		stopRepo:  stopRepo,
	}
}

// --- Rutas (Routes) ---

func (c *AdminController) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var route domain.Route
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := route.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.routeRepo.Create(r.Context(), &route); err != nil {
		log.Printf("ERROR CREANDO RUTA EN BD: %v", err)
		if errors.Is(err, postgres.ErrDuplicateKey) {
			writeError(w, http.StatusConflict, "Ya existe una ruta con ese nombre o código")
			return
		}
		writeError(w, http.StatusInternalServerError, "Error al crear la ruta")
		return
	}

	writeJSON(w, http.StatusCreated, route)
}

func (c *AdminController) GetRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := c.routeRepo.GetAll(r.Context())
	if err != nil {
		log.Printf("ERROR OBTENIENDO RUTAS DE BD: %v", err)
		writeError(w, http.StatusInternalServerError, "Error al obtener las rutas")
		return
	}
	if routes == nil {
		routes = []domain.Route{}
	}
	writeJSON(w, http.StatusOK, routes)
}

// RouteItem maneja GET, PUT, DELETE en /api/v1/admin/routes/{id} y /api/v1/admin/routes/{id}/stops
func (c *AdminController) RouteItem(w http.ResponseWriter, r *http.Request) {
	id, rest := idFromPath(r.URL.Path, "/api/v1/admin/routes/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de ruta requerido")
		return
	}
	if rest == "stops" {
		c.handleRouteStops(w, r, id)
		return
	}

	switch r.Method {
	case http.MethodGet:
		route, err := c.routeRepo.GetByID(r.Context(), id)
		if err != nil || route == nil {
			writeError(w, http.StatusNotFound, "Ruta no encontrada")
			return
		}
		writeJSON(w, http.StatusOK, route)

	case http.MethodPut:
		var route domain.Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		route.ID = id
		if err := route.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := c.routeRepo.Update(r.Context(), &route); err != nil {
			if errors.Is(err, postgres.ErrRecordNotFound) {
				writeError(w, http.StatusNotFound, "Ruta no encontrada")
				return
			}
			writeError(w, http.StatusInternalServerError, "Error al actualizar la ruta")
			return
		}
		writeJSON(w, http.StatusOK, route)

	case http.MethodDelete:
		count, err := c.routeRepo.CountBuses(r.Context(), id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error verificando dependencias de la ruta")
			return
		}
		if count > 0 {
			writeError(w, http.StatusConflict, "No se puede eliminar: la ruta tiene buses asignados. Reasígnalos o elimínalos primero.")
			return
		}
		if err := c.routeRepo.Delete(r.Context(), id); err != nil {
			if errors.Is(err, postgres.ErrRecordNotFound) {
				writeError(w, http.StatusNotFound, "Ruta no encontrada")
				return
			}
			writeError(w, http.StatusInternalServerError, "Error al eliminar la ruta")
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}

// --- Paraderos (Stops) ---

func (c *AdminController) handleRouteStops(w http.ResponseWriter, r *http.Request, routeID string) {
	if c.stopRepo == nil {
		writeError(w, http.StatusNotImplemented, "Repositorio de paraderos no configurado")
		return
	}

	switch r.Method {
	case http.MethodGet:
		stops, err := c.stopRepo.GetByRouteID(r.Context(), routeID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error listando paraderos")
			return
		}
		if stops == nil {
			stops = []domain.Stop{}
		}
		writeJSON(w, http.StatusOK, stops)

	case http.MethodPost:
		var stop domain.Stop
		if err := json.NewDecoder(r.Body).Decode(&stop); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		stop.RouteID = routeID
		if err := stop.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := c.stopRepo.Create(r.Context(), &stop); err != nil {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, stop)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}

func (c *AdminController) DeleteStop(w http.ResponseWriter, r *http.Request) {
	if c.stopRepo == nil {
		writeError(w, http.StatusNotImplemented, "Repositorio de paraderos no configurado")
		return
	}
	id, _ := idFromPath(r.URL.Path, "/api/v1/admin/stops/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de paradero requerido")
		return
	}
	if err := c.stopRepo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, postgres.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "Paradero no encontrado")
			return
		}
		writeError(w, http.StatusInternalServerError, "Error al eliminar el paradero")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Flota (Fleet) ---

func (c *AdminController) CreateFleet(w http.ResponseWriter, r *http.Request) {
	var fleet domain.Fleet
	if err := json.NewDecoder(r.Body).Decode(&fleet); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if fleet.Status == "" {
		fleet.Status = "active"
	}
	if err := fleet.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.fleetRepo.Create(r.Context(), &fleet); err != nil {
		writeError(w, http.StatusConflict, "Error al registrar el vehículo. Verifique que el Bus ID no esté duplicado.")
		return
	}
	writeJSON(w, http.StatusCreated, fleet)
}

func (c *AdminController) GetFleet(w http.ResponseWriter, r *http.Request) {
	fleet, err := c.fleetRepo.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error al obtener la flota")
		return
	}
	if fleet == nil {
		fleet = []domain.Fleet{}
	}
	writeJSON(w, http.StatusOK, fleet)
}

// FleetItem maneja GET, PUT, DELETE en /api/v1/admin/fleet/{id}
func (c *AdminController) FleetItem(w http.ResponseWriter, r *http.Request) {
	id, _ := idFromPath(r.URL.Path, "/api/v1/admin/fleet/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de unidad requerido")
		return
	}

	switch r.Method {
	case http.MethodGet:
		bus, err := c.fleetRepo.GetByID(r.Context(), id)
		if err != nil || bus == nil {
			writeError(w, http.StatusNotFound, "Unidad no encontrada")
			return
		}
		writeJSON(w, http.StatusOK, bus)

	case http.MethodPut:
		var fleet domain.Fleet
		if err := json.NewDecoder(r.Body).Decode(&fleet); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		fleet.ID = id
		if fleet.Status == "" {
			fleet.Status = "active"
		}
		if err := fleet.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := c.fleetRepo.Update(r.Context(), &fleet); err != nil {
			if errors.Is(err, postgres.ErrRecordNotFound) {
				writeError(w, http.StatusNotFound, "Unidad no encontrada")
				return
			}
			writeError(w, http.StatusInternalServerError, "Error al actualizar la unidad")
			return
		}
		writeJSON(w, http.StatusOK, fleet)

	case http.MethodDelete:
		if err := c.fleetRepo.Delete(r.Context(), id); err != nil {
			if errors.Is(err, postgres.ErrRecordNotFound) {
				writeError(w, http.StatusNotFound, "Unidad no encontrada")
				return
			}
			writeError(w, http.StatusInternalServerError, "Error al eliminar la unidad")
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}
