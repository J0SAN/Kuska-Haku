package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"kuska-haku-backend/internal/adapters/postgres"
	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type SuggestionController struct {
	repo ports.SuggestionRepository
}

func NewSuggestionController(repo ports.SuggestionRepository) *SuggestionController {
	return &SuggestionController{repo: repo}
}

type suggestionRequest struct {
	Message string `json:"message"`
}

func (h *SuggestionController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	var req suggestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}
	if err := domain.ValidateSuggestionMessage(req.Message); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	suggestion := domain.Suggestion{
		UserID:  claims.UserID,
		Message: strings.TrimSpace(req.Message),
	}
	if err := h.repo.Create(r.Context(), &suggestion); err != nil {
		writeError(w, http.StatusInternalServerError, "Error guardando la sugerencia")
		return
	}
	writeJSON(w, http.StatusCreated, suggestion)
}

// AdminCollection maneja GET /api/v1/admin/suggestions[?only=new]
func (h *SuggestionController) AdminCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}
	onlyNew := r.URL.Query().Get("only") == "new"
	list, err := h.repo.GetAll(r.Context(), onlyNew)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error listando sugerencias")
		return
	}
	if list == nil {
		list = []domain.Suggestion{}
	}
	writeJSON(w, http.StatusOK, list)
}

// AdminItem maneja PUT /api/v1/admin/suggestions/{id}/read
func (h *SuggestionController) AdminItem(w http.ResponseWriter, r *http.Request) {
	id, rest := idFromPath(r.URL.Path, "/api/v1/admin/suggestions/")
	if id == "" || rest != "read" {
		writeError(w, http.StatusNotFound, "Recurso no encontrado")
		return
	}
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}
	if err := h.repo.MarkRead(r.Context(), id); err != nil {
		if errors.Is(err, postgres.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "Sugerencia no encontrada")
			return
		}
		writeError(w, http.StatusInternalServerError, "Error actualizando la sugerencia")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
