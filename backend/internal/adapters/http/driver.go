package http

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

const (
	// maxPositionsPerBatch acota el tamaño máximo por envío (100 puntos ~ 5 min de señal perdida).
	maxPositionsPerBatch = 100
	// maxPositionsBodyBytes limita el tamaño del body HTTP a 64 KB para mitigar ataques DoS.
	maxPositionsBodyBytes = 64 << 10
)

type DriverController struct {
	fleet    ports.FleetRepository
	routes   ports.RouteRepository
	incomes  ports.IncomeRepository
	ingestor ports.PositionIngestor
}

func NewDriverController(
	f ports.FleetRepository,
	r ports.RouteRepository,
	i ports.IncomeRepository,
	ing ...ports.PositionIngestor,
) *DriverController {
	var ingestor ports.PositionIngestor
	if len(ing) > 0 {
		ingestor = ing[0]
	}
	return &DriverController{
		fleet:    f,
		routes:   r,
		incomes:  i,
		ingestor: ingestor,
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

type positionPointRequest struct {
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
	Lat       *float64  `json:"lat"`
	Lng       *float64  `json:"lng"`
	SpeedKmh  float64   `json:"speedKmh"`
	Heading   float64   `json:"heading"`
	Timestamp time.Time `json:"timestamp"`
}

func (p positionPointRequest) coordinates() (lat, lng float64) {
	if p.Latitude != nil {
		lat = *p.Latitude
	} else if p.Lat != nil {
		lat = *p.Lat
	}
	if p.Longitude != nil {
		lng = *p.Longitude
	} else if p.Lng != nil {
		lng = *p.Lng
	}
	return lat, lng
}

type positionsBatchRequest struct {
	Positions []positionPointRequest `json:"positions"`
}

// Positions recibe un lote de hasta 100 puntos GPS del conductor autenticado.
// SEGURIDAD: ignora cualquier busId/routeId enviado por el cliente y utiliza
// exclusivamente el bus y la ruta asignados al conductor en la base de datos.
func (h *DriverController) Positions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPositionsBodyBytes)
	var req positionsBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Cuerpo de la petición inválido o demasiado grande")
		return
	}

	n := len(req.Positions)
	if n == 0 {
		writeError(w, http.StatusBadRequest, "Envía al menos una posición")
		return
	}
	if n > maxPositionsPerBatch {
		writeError(w, http.StatusBadRequest, "Máximo 100 posiciones por lote; divide el envío en varios lotes")
		return
	}

	bus, err := h.fleet.GetByDriverID(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error buscando tu bus asignado")
		return
	}
	if bus == nil {
		writeError(w, http.StatusConflict, "No tienes un bus asignado; solicita al administrador que te asigne uno")
		return
	}

	route, err := h.routes.GetByID(r.Context(), bus.RouteID)
	if err != nil || route == nil {
		writeError(w, http.StatusConflict, "Tu bus no tiene una ruta válida asignada")
		return
	}

	batch := make([]domain.Position, 0, n)
	for _, p := range req.Positions {
		lat, lng := p.coordinates()
		batch = append(batch, domain.Position{
			BusID:     bus.BusID,
			RouteID:   route.ID,
			Latitude:  lat,
			Longitude: lng,
			SpeedKmh:  p.SpeedKmh,
			Heading:   p.Heading,
			Timestamp: p.Timestamp,
		})
	}

	res, err := h.ingestor.Ingest(r.Context(), batch)
	if err != nil {
		// 503 indica a la app móvil que conserve el lote en el buffer y reintente luego.
		writeError(w, http.StatusServiceUnavailable, "No se pudo procesar el lote; reintenta")
		return
	}
	writeJSON(w, http.StatusOK, res)
}
