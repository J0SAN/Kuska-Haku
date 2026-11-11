package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"kuska-haku-backend/internal/adapters/postgres"
	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
	"kuska-haku-backend/internal/usecases"
)

const maxPublicFeedbackBodyBytes = 4 * 1024

// VehicleQRController administra los QR por vehículo (solo admin).
type VehicleQRController struct {
	manager *usecases.VehicleQRManager
	repo    ports.VehicleQRRepository
}

func NewVehicleQRController(manager *usecases.VehicleQRManager, repo ports.VehicleQRRepository) *VehicleQRController {
	return &VehicleQRController{manager: manager, repo: repo}
}

// List maneja GET /api/v1/admin/vehicle-qr.
func (c *VehicleQRController) List(w http.ResponseWriter, r *http.Request) {
	list, err := c.repo.ListWithFleet(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error listando los códigos QR")
		return
	}
	if list == nil {
		list = []domain.VehicleQRSummary{}
	}
	writeJSON(w, http.StatusOK, list)
}

// Generate maneja POST /api/v1/admin/vehicle-qr/{fleetId}.
func (c *VehicleQRController) Generate(w http.ResponseWriter, r *http.Request) {
	fleetID := r.PathValue("fleetId")
	if fleetID == "" {
		writeError(w, http.StatusBadRequest, "ID de unidad requerido")
		return
	}
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	qr, err := c.manager.Generate(r.Context(), fleetID, claims.UserID)
	if err != nil {
		writeQRManagerError(w, err, "Error generando el código QR")
		return
	}
	writeJSON(w, http.StatusCreated, qr)
}

// Rotate maneja POST /api/v1/admin/vehicle-qr/{fleetId}/rotate.
func (c *VehicleQRController) Rotate(w http.ResponseWriter, r *http.Request) {
	fleetID := r.PathValue("fleetId")
	if fleetID == "" {
		writeError(w, http.StatusBadRequest, "ID de unidad requerido")
		return
	}
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	qr, err := c.manager.Rotate(r.Context(), fleetID, claims.UserID)
	if err != nil {
		writeQRManagerError(w, err, "Error regenerando el código QR")
		return
	}
	writeJSON(w, http.StatusCreated, qr)
}

// Revoke maneja DELETE /api/v1/admin/vehicle-qr/{fleetId}.
func (c *VehicleQRController) Revoke(w http.ResponseWriter, r *http.Request) {
	fleetID := r.PathValue("fleetId")
	if fleetID == "" {
		writeError(w, http.StatusBadRequest, "ID de unidad requerido")
		return
	}

	if err := c.manager.Revoke(r.Context(), fleetID); err != nil {
		writeQRManagerError(w, err, "Error revocando el código QR")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeQRManagerError(w http.ResponseWriter, err error, fallback string) {
	switch {
	case errors.Is(err, domain.ErrQRAlreadyExists):
		writeError(w, http.StatusConflict, "El vehículo ya tiene un código QR activo")
	case errors.Is(err, domain.ErrQRNotAvailable), errors.Is(err, postgres.ErrRecordNotFound):
		writeError(w, http.StatusNotFound, "Unidad no encontrada")
	default:
		writeError(w, http.StatusInternalServerError, fallback)
	}
}

// VehicleFeedbackAdminController lista y modera las calificaciones.
type VehicleFeedbackAdminController struct {
	repo ports.VehicleFeedbackRepository
	now  func() time.Time
}

func NewVehicleFeedbackAdminController(repo ports.VehicleFeedbackRepository) *VehicleFeedbackAdminController {
	return &VehicleFeedbackAdminController{repo: repo, now: time.Now}
}

// List maneja GET /api/v1/admin/vehicle-feedback?fleetId=&status=
func (c *VehicleFeedbackAdminController) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status != "" && domain.ValidateFeedbackStatus(status) != nil {
		writeError(w, http.StatusBadRequest, "Estado de calificación inválido")
		return
	}

	list, err := c.repo.List(r.Context(), domain.VehicleFeedbackFilter{
		FleetID: r.URL.Query().Get("fleetId"),
		Status:  status,
		Limit:   200,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error listando las calificaciones")
		return
	}
	if list == nil {
		list = []domain.VehicleFeedback{}
	}
	writeJSON(w, http.StatusOK, list)
}

// Summary maneja GET /api/v1/admin/vehicle-feedback/summary?days=30
func (c *VehicleFeedbackAdminController) Summary(w http.ResponseWriter, r *http.Request) {
	days := 30
	if raw := r.URL.Query().Get("days"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 365 {
			writeError(w, http.StatusBadRequest, "El periodo debe estar entre 1 y 365 días")
			return
		}
		days = parsed
	}

	fromDate := domain.DateLima(c.now().AddDate(0, 0, -(days - 1)))
	summaries, err := c.repo.Summaries(r.Context(), fromDate)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error calculando los promedios")
		return
	}
	if summaries == nil {
		summaries = []domain.VehicleRatingSummary{}
	}
	writeJSON(w, http.StatusOK, summaries)
}

type feedbackStatusRequest struct {
	Status string `json:"status"`
}

// UpdateStatus maneja PUT /api/v1/admin/vehicle-feedback/{id}/status
func (c *VehicleFeedbackAdminController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID de calificación requerido")
		return
	}

	var req feedbackStatusRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxPublicFeedbackBodyBytes)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}
	if err := domain.ValidateFeedbackStatus(req.Status); err != nil {
		writeError(w, http.StatusBadRequest, "Estado de calificación inválido")
		return
	}

	if err := c.repo.UpdateStatus(r.Context(), id, req.Status); err != nil {
		if errors.Is(err, postgres.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "Calificación no encontrada")
			return
		}
		writeError(w, http.StatusInternalServerError, "Error actualizando la calificación")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PublicFeedbackController atiende la página pública del QR. El token es la
// credencial de acceso y el limitador por IP es la protección.
type PublicFeedbackController struct {
	submitter *usecases.VehicleFeedbackSubmitter
}

func NewPublicFeedbackController(submitter *usecases.VehicleFeedbackSubmitter) *PublicFeedbackController {
	return &PublicFeedbackController{submitter: submitter}
}

type publicFeedbackRequest struct {
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
	DeviceID string `json:"deviceId"`
}

// Resolve maneja GET /api/v1/public/vehicle-feedback/{token}
func (c *PublicFeedbackController) Resolve(w http.ResponseWriter, r *http.Request) {
	target, err := c.submitter.Resolve(r.Context(), r.PathValue("token"))
	if err != nil {
		writePublicFeedbackError(w, err)
		return
	}
	// Solo datos comprobables por el pasajero: nunca el conductor.
	writeJSON(w, http.StatusOK, struct {
		BusCode   string `json:"busCode"`
		RouteName string `json:"routeName"`
	}{
		BusCode:   target.BusCode,
		RouteName: target.RouteName,
	})
}

// Submit maneja POST /api/v1/public/vehicle-feedback/{token}
func (c *PublicFeedbackController) Submit(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxPublicFeedbackBodyBytes)

	var req publicFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	_, err := c.submitter.Submit(r.Context(), r.PathValue("token"), usecases.VehicleFeedbackInput{
		Rating:   req.Rating,
		Comment:  req.Comment,
		DeviceID: req.DeviceID,
	})
	if err != nil {
		writePublicFeedbackError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, struct {
		Status string `json:"status"`
	}{Status: "received"})
}

func writePublicFeedbackError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrQRNotAvailable):
		writeError(w, http.StatusNotFound, "Este código QR ya no está disponible.")
	case errors.Is(err, domain.ErrInvalidRating):
		writeError(w, http.StatusBadRequest, "Elige una calificación de 1 a 5 estrellas.")
	case errors.Is(err, domain.ErrInvalidFeedbackComment):
		writeError(w, http.StatusBadRequest, "El comentario no puede superar los 500 caracteres.")
	case errors.Is(err, domain.ErrInvalidDeviceID):
		writeError(w, http.StatusBadRequest, "No pudimos identificar tu dispositivo. Recarga la página.")
	case errors.Is(err, domain.ErrFeedbackAlreadySubmitted):
		writeError(w, http.StatusConflict, "Ya calificaste este vehículo hoy. ¡Gracias!")
	default:
		writeError(w, http.StatusInternalServerError, "No pudimos guardar tu calificación. Inténtalo otra vez.")
	}
}
