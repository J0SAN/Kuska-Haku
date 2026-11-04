package http

import (
	"math"
	"net/http"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type PassengerController struct {
	routes   ports.RouteRepository
	fleet    ports.FleetRepository
	store    ports.PositionStore
	progress ports.RouteProgressRepository
}

func NewPassengerController(
	r ports.RouteRepository,
	f ports.FleetRepository,
	s ports.PositionStore,
	p ports.RouteProgressRepository,
) *PassengerController {
	return &PassengerController{
		routes:   r,
		fleet:    f,
		store:    s,
		progress: p,
	}
}

func (h *PassengerController) ListRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}
	routes, err := h.routes.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error listando rutas")
		return
	}
	if routes == nil {
		routes = []domain.Route{}
	}
	writeJSON(w, http.StatusOK, routes)
}

type etaStop struct {
	StopID         string  `json:"stopId"`
	Name           string  `json:"name"`
	StopOrder      int     `json:"stopOrder"`
	DistanceMeters float64 `json:"distanceMeters"`
	ETASeconds     *int    `json:"etaSeconds"`
	Passed         bool    `json:"passed"`
}

type etaBus struct {
	BusID             string            `json:"busId"`
	Position          domain.Coordinate `json:"position"`
	SpeedKmh          float64           `json:"speedKmh"`
	UpdatedSecondsAgo int               `json:"updatedSecondsAgo"`
	OffRoute          bool              `json:"offRoute"`
	OffRouteMeters    float64           `json:"offRouteMeters"`
	Stops             []etaStop         `json:"stops"`
}

type etaResponse struct {
	RouteID string   `json:"routeId"`
	Buses   []etaBus `json:"buses"`
}

// RouteItem atiende GET /api/v1/routes/{id}/eta
func (h *PassengerController) RouteItem(w http.ResponseWriter, r *http.Request) {
	id, rest := idFromPath(r.URL.Path, "/api/v1/routes/")
	if id == "" || rest != "eta" {
		writeError(w, http.StatusNotFound, "Recurso no encontrado")
		return
	}
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}
	h.eta(w, r, id)
}

func (h *PassengerController) eta(w http.ResponseWriter, r *http.Request, routeID string) {
	ctx := r.Context()

	route, err := h.routes.GetByID(ctx, routeID)
	if err != nil || route == nil {
		writeError(w, http.StatusNotFound, "Ruta no encontrada")
		return
	}

	buses, err := h.fleet.GetByRouteID(ctx, routeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error listando los buses de la ruta")
		return
	}

	resp := etaResponse{RouteID: routeID, Buses: []etaBus{}}
	now := time.Now()

	for _, bus := range buses {
		if bus.Status != "" && bus.Status != "active" {
			continue
		}
		pos, err := h.store.GetLastPosition(ctx, bus.BusID)
		if err != nil || pos == nil {
			continue
		}

		prog, err := h.progress.Progress(ctx, routeID, domain.Coordinate{
			Latitude:  pos.Latitude,
			Longitude: pos.Longitude,
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error calculando el ETA")
			return
		}

		offRoute := domain.IsOffRoute(prog.OffRouteMeters)
		ageSeconds := int(now.Sub(pos.Timestamp).Seconds())
		if ageSeconds < 0 {
			ageSeconds = 0
		}

		eb := etaBus{
			BusID:             bus.BusID,
			Position:          domain.Coordinate{Latitude: pos.Latitude, Longitude: pos.Longitude},
			SpeedKmh:          pos.SpeedKmh,
			UpdatedSecondsAgo: ageSeconds,
			OffRoute:          offRoute,
			OffRouteMeters:    math.Round(prog.OffRouteMeters),
			Stops:             make([]etaStop, 0, len(prog.Stops)),
		}

		for _, s := range prog.Stops {
			es := etaStop{
				StopID:         s.StopID,
				Name:           s.Name,
				StopOrder:      s.StopOrder,
				DistanceMeters: math.Round(s.DistanceMeters),
				Passed:         s.Passed,
			}
			if !s.Passed && !offRoute {
				secs := int(math.Round(domain.EstimateETA(s.DistanceMeters, pos.SpeedKmh).Seconds()))
				es.ETASeconds = &secs
			}
			eb.Stops = append(eb.Stops, es)
		}
		resp.Buses = append(resp.Buses, eb)
	}

	writeJSON(w, http.StatusOK, resp)
}
