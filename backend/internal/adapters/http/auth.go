package http

import (
	"encoding/json"
	"net/http"

	// Ajusta el nombre de tu módulo si es distinto
	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type AuthController struct {
	userRepo ports.UserRepository
	authSvc  ports.AuthService
}

func NewAuthController(userRepo ports.UserRepository, authSvc ports.AuthService) *AuthController {
	return &AuthController{userRepo: userRepo, authSvc: authSvc}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Cuerpo de petición inválido", http.StatusBadRequest)
		return
	}

	hashedPassword, err := c.authSvc.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error interno al procesar seguridad", http.StatusInternalServerError)
		return
	}

	user, err := domain.NewUser(req.Email, hashedPassword, req.FullName, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.userRepo.Create(r.Context(), user); err != nil {
		http.Error(w, "Error al crear usuario o el correo ya está registrado", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario creado exitosamente"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Cuerpo de petición inválido", http.StatusBadRequest)
		return
	}

	user, err := c.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	if !c.authSvc.CheckPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	token, err := c.authSvc.GenerateToken(user)
	if err != nil {
		http.Error(w, "Error al generar token de sesión", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
