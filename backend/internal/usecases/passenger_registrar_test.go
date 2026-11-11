package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type fakeUsers struct {
	ports.UserRepository
	createErr error
	calls     int
	created   *domain.User
}

func (f *fakeUsers) Create(_ context.Context, user *domain.User) error {
	f.calls++
	if f.createErr != nil {
		return f.createErr
	}
	user.ID = "user-1"
	copied := *user
	f.created = &copied
	return nil
}

type fakeHasher struct {
	ports.PasswordHasher
	hashErr error
}

func (f *fakeHasher) Hash(password string) (string, error) {
	if f.hashErr != nil {
		return "", f.hashErr
	}
	return "hash:" + password, nil
}

type fakeTokens struct {
	ports.TokenService
	generateErr error
}

func (f *fakeTokens) Generate(_ domain.User) (string, error) {
	if f.generateErr != nil {
		return "", f.generateErr
	}
	return "token-1", nil
}

func newRegistrar(users *fakeUsers) (*PassengerRegistrar, time.Time) {
	now := time.Date(2026, 11, 11, 15, 0, 0, 0, time.UTC)
	registrar := NewPassengerRegistrar(users, &fakeHasher{}, &fakeTokens{})
	registrar.now = func() time.Time { return now }
	return registrar, now
}

func validInput() RegisterInput {
	return RegisterInput{
		DNI:         "12345678",
		CheckDigit:  "1",
		FirstNames:  "  Roxana   Karina ",
		LastNames:   "Delgado Huamán",
		Email:       " Roxana@Correo.COM ",
		Password:    "secreto123",
		AcceptTerms: true,
	}
}

func TestPassengerRegistrar_HappyPath(t *testing.T) {
	users := &fakeUsers{}
	registrar, now := newRegistrar(users)

	user, token, err := registrar.Register(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if token != "token-1" {
		t.Errorf("token = %q", token)
	}
	if users.calls != 1 {
		t.Fatalf("Create() llamado %d veces, se esperaba 1", users.calls)
	}
	if user.Role != string(domain.RolePassenger) {
		t.Errorf("rol = %q, se esperaba passenger", user.Role)
	}
	if user.FullName != "Roxana Karina Delgado Huamán" {
		t.Errorf("nombre completo = %q", user.FullName)
	}
	if user.Email != "roxana@correo.com" {
		t.Errorf("correo = %q, se esperaba normalizado", user.Email)
	}
	if user.DNI == nil || *user.DNI != "12345678" {
		t.Errorf("DNI = %v", user.DNI)
	}
	if user.TermsAcceptedAt == nil || !user.TermsAcceptedAt.Equal(now) {
		t.Errorf("TermsAcceptedAt = %v, se esperaba %v", user.TermsAcceptedAt, now)
	}
	if user.PasswordHash != "hash:secreto123" {
		t.Errorf("PasswordHash = %q", user.PasswordHash)
	}
}

func TestPassengerRegistrar_ErroresDeCreate(t *testing.T) {
	tests := []struct {
		name      string
		createErr error
	}{
		{name: "DNI duplicado", createErr: domain.ErrDNIAlreadyRegistered},
		{name: "correo duplicado", createErr: domain.ErrEmailAlreadyRegistered},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &fakeUsers{createErr: tt.createErr}
			registrar, _ := newRegistrar(users)

			_, _, err := registrar.Register(context.Background(), validInput())
			if !errors.Is(err, tt.createErr) {
				t.Errorf("se esperaba %v, obtuve %v", tt.createErr, err)
			}
		})
	}
}

func TestPassengerRegistrar_ValidaAntesDeCrear(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*RegisterInput)
		want   error
	}{
		{
			name:   "sin aceptar términos",
			mutate: func(in *RegisterInput) { in.AcceptTerms = false },
			want:   domain.ErrTermsNotAccepted,
		},
		{
			name:   "DNI inválido",
			mutate: func(in *RegisterInput) { in.DNI = "1234" },
			want:   domain.ErrInvalidDNI,
		},
		{
			name:   "dígito verificador incorrecto",
			mutate: func(in *RegisterInput) { in.CheckDigit = "2" },
			want:   domain.ErrInvalidCheckDigit,
		},
		{
			name:   "nombres inválidos",
			mutate: func(in *RegisterInput) { in.FirstNames = "R2D2" },
			want:   domain.ErrInvalidName,
		},
		{
			name:   "apellidos inválidos",
			mutate: func(in *RegisterInput) { in.LastNames = "😀" },
			want:   domain.ErrInvalidName,
		},
		{
			name:   "correo inválido",
			mutate: func(in *RegisterInput) { in.Email = "sin-arroba" },
			want:   domain.ErrInvalidEmail,
		},
		{
			name:   "contraseña corta",
			mutate: func(in *RegisterInput) { in.Password = "corta" },
			want:   domain.ErrWeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &fakeUsers{}
			registrar, _ := newRegistrar(users)

			input := validInput()
			tt.mutate(&input)

			_, _, err := registrar.Register(context.Background(), input)
			if !errors.Is(err, tt.want) {
				t.Errorf("se esperaba %v, obtuve %v", tt.want, err)
			}
			if users.calls != 0 {
				t.Errorf("Create() no debía llamarse, se llamó %d veces", users.calls)
			}
		})
	}
}
