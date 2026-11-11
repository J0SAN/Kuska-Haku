package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
	"kuska-haku-backend/internal/usecases"
)

// maxAuthBodyBytes limita el cuerpo de los endpoints públicos de autenticación.
const maxAuthBodyBytes = 4 << 10

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
	userRepo  ports.UserRepository
	hasher    ports.PasswordHasher
	tokens    ports.TokenService
	registrar *usecases.PassengerRegistrar
}

func NewAuthController(
	userRepo ports.UserRepository,
	hasher ports.PasswordHasher,
	tokens ports.TokenService,
	registrar *usecases.PassengerRegistrar,
) *AuthController {
	return &AuthController{
		userRepo:  userRepo,
		hasher:    hasher,
		tokens:    tokens,
		registrar: registrar,
	}
}

// RegisterRequest documenta el contrato del registro público. Role se acepta
// solo para poder rechazarlo explícitamente: el registro únicamente crea pasajeros.
type RegisterRequest struct {
	DNI         string `json:"dni"`
	CheckDigit  string `json:"checkDigit"`
	FirstNames  string `json:"firstNames"`
	LastNames   string `json:"lastNames"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AcceptTerms bool   `json:"acceptTerms"`
	Role        string `json:"role"`
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	if req.Role != "" && req.Role != string(domain.RolePassenger) {
		writeError(w, http.StatusForbidden, domain.ErrRoleNotAllowed.Error())
		return
	}

	user, token, err := c.registrar.Register(r.Context(), usecases.RegisterInput{
		DNI:         req.DNI,
		CheckDigit:  req.CheckDigit,
		FirstNames:  req.FirstNames,
		LastNames:   req.LastNames,
		Email:       req.Email,
		Password:    req.Password,
		AcceptTerms: req.AcceptTerms,
	})
	if err != nil {
		writeRegisterError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "Usuario creado exitosamente",
		"token":   token,
		"user":    toUserResponse(user),
	})
}

func writeRegisterError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrTermsNotAccepted),
		errors.Is(err, domain.ErrInvalidDNI),
		errors.Is(err, domain.ErrInvalidCheckDigit),
		errors.Is(err, domain.ErrInvalidName),
		errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrWeakPassword):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrDNIAlreadyRegistered):
		writeError(w, http.StatusConflict, "Ya existe una cuenta con este DNI. Si no la creaste tú, contacta a soporte.")
	case errors.Is(err, domain.ErrEmailAlreadyRegistered):
		writeError(w, http.StatusConflict, "Este correo ya está registrado.")
	case errors.Is(err, domain.ErrRoleNotAllowed):
		writeError(w, http.StatusForbidden, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "Error interno al crear la cuenta")
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Cuerpo de petición inválido")
		return
	}

	user, err := c.userRepo.GetByEmail(r.Context(), domain.NormalizeEmail(req.Email))
	if err != nil || user == nil {
		writeError(w, http.StatusUnauthorized, "Credenciales inválidas")
		return
	}

	if err := c.hasher.Compare(user.PasswordHash, req.Password); err != nil {
		writeError(w, http.StatusUnauthorized, "Credenciales inválidas")
		return
	}

	token, err := c.tokens.Generate(*user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error al generar token de sesión")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  toUserResponse(user),
	})
}

func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}

	user, err := c.userRepo.GetByID(r.Context(), claims.UserID)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}
