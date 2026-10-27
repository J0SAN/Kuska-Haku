package http

import (
	"encoding/json"
	"math"
	"net/http"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type DriverController struct {
	fleet   ports.FleetRepository
	routes  ports.RouteRepository
	incomes ports.IncomeRepository
}

func NewDriverController(
	f ports.FleetRepository,
	r ports.RouteRepository,
	i ports.IncomeRepository,
) *DriverController {
	return &DriverController{
		fleet:   f,
		routes:  r,
		incomes: i,
	}
}

type myBusRouteSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type myBusResponse struct {
	Bus   domain.Fleet       `json:"bus"`
	Route *myBusRouteSummary `json:"route"`
}

// MyBus devuelve el bus activo asignado al conductor autenticado y su ruta.
func (h *DriverController) MyBus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	bus, err := h.fleet.GetByDriverID(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error buscando tu bus asignado")
		return
	}
	if bus == nil {
		writeError(w, http.StatusNotFound, "Todavía no tienes un bus asignado; solicita al administrador que te asigne uno")
		return
	}

	resp := myBusResponse{Bus: *bus}
	if route, err := h.routes.GetByID(r.Context(), bus.RouteID); err == nil && route != nil {
		resp.Route = &myBusRouteSummary{ID: route.ID, Name: route.Name}
	}
	writeJSON(w, http.StatusOK, resp)
}

type incomeRequest struct {
	Passengers int     `json:"passengers"`
	Amount     float64 `json:"amount"`
	Notes      string  `json:"notes"`
}

type todayIncomesResponse struct {
	Date            string                `json:"date"`
	Records         []domain.IncomeRecord `json:"records"`
	TotalPassengers int                   `json:"totalPassengers"`
	TotalAmount     float64               `json:"totalAmount"`
}

// Incomes maneja GET (listar ingresos de hoy) y POST (registrar nuevo ingreso del turno).
func (h *DriverController) Incomes(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	switch r.Method {
	case http.MethodGet:
		today := domain.TodayLima()
		records, err := h.incomes.GetByDriverAndDate(r.Context(), claims.UserID, today)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error listando tus registros de hoy")
			return
		}
		if records == nil {
			records = []domain.IncomeRecord{}
		}
		resp := todayIncomesResponse{Date: today, Records: records}
		for _, rec := range records {
			resp.TotalPassengers += rec.Passengers
			resp.TotalAmount += rec.Amount
		}
		resp.TotalAmount = math.Round(resp.TotalAmount*100) / 100
		writeJSON(w, http.StatusOK, resp)

	case http.MethodPost:
		var req incomeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Cuerpo de petición inválido")
			return
		}

		// El bus y la fecha se obtienen exclusivamente en el servidor a partir del JWT del conductor.
		bus, err := h.fleet.GetByDriverID(r.Context(), claims.UserID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error buscando tu bus asignado")
			return
		}
		if bus == nil {
			writeError(w, http.StatusConflict, "No tienes un bus asignado; solicita al administrador que te asigne uno")
			return
		}

		record := domain.IncomeRecord{
			FleetID:    bus.ID,
			BusCode:    bus.BusID,
			DriverID:   claims.UserID,
			RecordDate: domain.TodayLima(),
			Passengers: req.Passengers,
			Amount:     req.Amount,
			Notes:      req.Notes,
		}
		if err := record.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := h.incomes.Create(r.Context(), &record); err != nil {
			writeError(w, http.StatusInternalServerError, "Error guardando el registro de ingreso")
			return
		}
		writeJSON(w, http.StatusCreated, record)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}
