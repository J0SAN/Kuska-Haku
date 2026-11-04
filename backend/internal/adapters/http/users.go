package http

import (
	"net/http"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

// userListItemResponse define únicamente los datos públicos que viajan al panel web,
// con etiquetas en minúscula (camelCase) y sin el campo PasswordHash.
type userListItemResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullName"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// UsersHandler atiende GET /api/v1/admin/users?role=driver
// Permite al panel administrativo listar los conductores disponibles para asignarles una unidad.
func UsersHandler(users ports.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "Método no permitido")
			return
		}
		roleStr := r.URL.Query().Get("role")
		if roleStr == "" {
			roleStr = r.URL.Query().Get("rol")
		}
		role := domain.Role(roleStr)
		if !role.Valid() {
			writeError(w, http.StatusBadRequest, "Rol inválido: usa ?role=driver, passenger o admin")
			return
		}
		list, err := users.ListByRole(r.Context(), role)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error listando usuarios")
			return
		}

		out := make([]userListItemResponse, 0, len(list))
		for _, u := range list {
			out = append(out, userListItemResponse{
				ID:        u.ID,
				Email:     u.Email,
				FullName:  u.FullName,
				Role:      u.Role,
				CreatedAt: u.CreatedAt,
			})
		}
		writeJSON(w, http.StatusOK, out)
	}
}
