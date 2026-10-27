package http

import (
	"encoding/json"
	"net/http"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
}

func toUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		FullName: u.FullName,
		Role:     u.Role,
	}
}

type AuthController struct {
	userRepo ports.UserRepository
	hasher   ports.PasswordHasher
	tokens   ports.TokenService
}

func NewAuthController(userRepo ports.UserRepository, hasher ports.PasswordHasher, tokens ports.TokenService) *AuthController {
	return &AuthController{
		userRepo: userRepo,
		hasher:   hasher,
		tokens:   tokens,
	}
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

	if len(req.Password) < 6 {
		http.Error(w, "La contraseña debe tener al menos 6 caracteres", http.StatusBadRequest)
		return
	}

	hashedPassword, err := c.hasher.Hash(req.Password)
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

	token, err := c.tokens.Generate(*user)
	if err != nil {
		http.Error(w, "Error al generar token de sesión", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Usuario creado exitosamente",
		"token":   token,
		"user":    toUserResponse(user),
	})
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
	if err != nil || user == nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	if err := c.hasher.Compare(user.PasswordHash, req.Password); err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	token, err := c.tokens.Generate(*user)
	if err != nil {
		http.Error(w, "Error al generar token de sesión", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token": token,
		"user":  toUserResponse(user),
	})
}

func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	user, err := c.userRepo.GetByID(r.Context(), claims.UserID)
	if err != nil || user == nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toUserResponse(user))
}
