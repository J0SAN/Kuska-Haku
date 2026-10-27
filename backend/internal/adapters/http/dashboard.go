package http

import (
	"net/http"
	"strconv"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

const (
	defaultDashboardDays = 7
	maxDashboardDays     = 90
)

// DashboardHandler atiende GET /api/v1/admin/dashboard?days=7
func DashboardHandler(stats ports.StatsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
			return
		}
		days := defaultDashboardDays
		if raw := r.URL.Query().Get("days"); raw != "" {
			n, err := strconv.Atoi(raw)
			if err != nil || n < 1 || n > maxDashboardDays {
				writeError(w, http.StatusBadRequest, "days debe ser un entero entre 1 y 90")
				return
			}
			days = n
		}
		d, err := stats.GetDashboard(r.Context(), domain.TodayLima(), days)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error calculando el dashboard")
			return
		}
		writeJSON(w, http.StatusOK, d)
	}
}
