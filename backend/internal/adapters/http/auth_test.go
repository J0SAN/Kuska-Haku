package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
	"kuska-haku-backend/internal/usecases"
)

type stubUserRepo struct {
	ports.UserRepository
	createErr error
	byEmail   *domain.User
	created   *domain.User
}

func (s *stubUserRepo) Create(_ context.Context, user *domain.User) error {
	if s.createErr != nil {
		return s.createErr
	}
	user.ID = "user-1"
	copied := *user
	s.created = &copied
	return nil
}

func (s *stubUserRepo) GetByEmail(_ context.Context, _ string) (*domain.User, error) {
	if s.byEmail == nil {
		return nil, errors.New("usuario no encontrado")
	}
	return s.byEmail, nil
}

type stubHasher struct {
	ports.PasswordHasher
}

func (s *stubHasher) Hash(password string) (string, error) {
	return "hash:" + password, nil
}

func (s *stubHasher) Compare(hash, password string) error {
	if hash != "hash:"+password {
		return errors.New("contraseña incorrecta")
	}
	return nil
}

type stubTokens struct {
	ports.TokenService
}

func (s *stubTokens) Generate(_ domain.User) (string, error) {
	return "token-1", nil
}

func newTestAuthController(repo ports.UserRepository) *AuthController {
	hasher := &stubHasher{}
	tokens := &stubTokens{}
	return NewAuthController(repo, hasher, tokens, usecases.NewPassengerRegistrar(repo, hasher, tokens))
}

func registerBody(overrides map[string]any) *strings.Reader {
	body := map[string]any{
		"dni":         "12345678",
		"checkDigit":  "1",
		"firstNames":  "Roxana Karina",
		"lastNames":   "Delgado Huamán",
		"email":       "roxana@correo.com",
		"password":    "secreto123",
		"acceptTerms": true,
	}
	for key, value := range overrides {
		body[key] = value
	}
	raw, _ := json.Marshal(body)
	return strings.NewReader(string(raw))
}

func TestRegister_RechazaRolDistintoDePasajero(t *testing.T) {
	repo := &stubUserRepo{}
	controller := newTestAuthController(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", registerBody(map[string]any{"role": "admin"}))
	rec := httptest.NewRecorder()
	controller.Register(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("esperaba 403, obtuve %d: %s", rec.Code, rec.Body.String())
	}
	if repo.created != nil {
		t.Error("no debía crear ningún usuario")
	}
}

func TestRegister_RechazaCuerpoMayorA4KB(t *testing.T) {
	controller := newTestAuthController(&stubUserRepo{})

	body := `{"email":"` + strings.Repeat("a", maxAuthBodyBytes+100) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	rec := httptest.NewRecorder()
	controller.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuve %d", rec.Code)
	}
}

func TestRegister_MapeaErrores(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]any
		createErr error
		wantCode  int
	}{
		{name: "términos no aceptados", overrides: map[string]any{"acceptTerms": false}, wantCode: 400},
		{name: "DNI inválido", overrides: map[string]any{"dni": "1234"}, wantCode: 400},
		{name: "dígito incorrecto", overrides: map[string]any{"checkDigit": "2"}, wantCode: 400},
		{name: "nombres inválidos", overrides: map[string]any{"firstNames": "R2D2"}, wantCode: 400},
		{name: "correo inválido", overrides: map[string]any{"email": "sin-arroba"}, wantCode: 400},
		{name: "contraseña corta", overrides: map[string]any{"password": "corta"}, wantCode: 400},
		{name: "DNI duplicado", createErr: domain.ErrDNIAlreadyRegistered, wantCode: 409},
		{name: "correo duplicado", createErr: domain.ErrEmailAlreadyRegistered, wantCode: 409},
		{name: "error de base", createErr: errors.New("conexión caída"), wantCode: 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubUserRepo{createErr: tt.createErr}
			controller := newTestAuthController(repo)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", registerBody(tt.overrides))
			rec := httptest.NewRecorder()
			controller.Register(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("esperaba %d, obtuve %d: %s", tt.wantCode, rec.Code, rec.Body.String())
			}

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("la respuesta no es JSON: %v", err)
			}
			if body["error"] == nil || body["message"] == nil {
				t.Errorf("la respuesta debe tener {error, message}: %v", body)
			}
		})
	}
}

func TestRegister_CreaPasajeroSinExponerDNI(t *testing.T) {
	repo := &stubUserRepo{}
	controller := newTestAuthController(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", registerBody(nil))
	rec := httptest.NewRecorder()
	controller.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuve %d: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Token string         `json:"token"`
		User  map[string]any `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("respuesta inválida: %v", err)
	}
	if body.Token != "token-1" {
		t.Errorf("token = %q", body.Token)
	}
	if body.User["role"] != "passenger" {
		t.Errorf("rol = %v, se esperaba passenger", body.User["role"])
	}
	if _, exists := body.User["dni"]; exists {
		t.Error("la respuesta no debe incluir el DNI")
	}
	if repo.created == nil || repo.created.DNI == nil || *repo.created.DNI != "12345678" {
		t.Errorf("el DNI no se guardó: %+v", repo.created)
	}
}

func TestLogin_ErrorEnJSON(t *testing.T) {
	controller := newTestAuthController(&stubUserRepo{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":"nadie@correo.com","password":"x"}`))
	rec := httptest.NewRecorder()
	controller.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("esperaba 401, obtuve %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("la respuesta no es JSON: %v", err)
	}
	if body["error"] != "Credenciales inválidas" {
		t.Errorf("error = %v", body["error"])
	}
}
