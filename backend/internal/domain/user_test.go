package domain

import (
	"errors"
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		passwordHash  string
		fullName      string
		role          string
		wantErr       bool
		expectedRole  string
		expectedEmail string
	}{
		{
			name:          "Usuario válido con rol explícito",
			email:         "admin@kuskahaku.pe",
			passwordHash:  "hash_seguro_123",
			fullName:      "Ana Admin",
			role:          "admin",
			expectedRole:  "admin",
			expectedEmail: "admin@kuskahaku.pe",
		},
		{
			name:          "Usuario válido con rol por defecto",
			email:         "pasajero@kuskahaku.pe",
			passwordHash:  "hash_seguro_123",
			fullName:      "Juan Pasajero",
			role:          "",
			expectedRole:  "passenger",
			expectedEmail: "pasajero@kuskahaku.pe",
		},
		{
			name:          "Normaliza mayúsculas y espacios del correo",
			email:         "  CARLOS@Kuskahaku.PE ",
			passwordHash:  "hash_seguro_123",
			fullName:      "Carlos Rojas",
			role:          "passenger",
			expectedRole:  "passenger",
			expectedEmail: "carlos@kuskahaku.pe",
		},
		{
			name:         "Fallo por correo sin arroba",
			email:        "correoinvalido.pe",
			passwordHash: "hash_seguro_123",
			fullName:     "Carlos Rojas",
			role:         "passenger",
			wantErr:      true,
		},
		{
			name:         "Fallo por rol inválido",
			email:        "otro@kuskahaku.pe",
			passwordHash: "hash_seguro_123",
			fullName:     "Otro Usuario",
			role:         "superadmin",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.passwordHash, tt.fullName, tt.role)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() arrojó error = %v, se esperaba error = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user.Role != tt.expectedRole {
					t.Errorf("NewUser() asignó rol = %v, se esperaba rol = %v", user.Role, tt.expectedRole)
				}
				if user.Email != tt.expectedEmail {
					t.Errorf("NewUser() guardó correo = %v, se esperaba = %v", user.Email, tt.expectedEmail)
				}
			}
		})
	}
}

func TestNormalizeEmail(t *testing.T) {
	if got := NormalizeEmail(" Jose@Gmail.COM "); got != "jose@gmail.com" {
		t.Errorf("NormalizeEmail() = %q, se esperaba jose@gmail.com", got)
	}
}

func TestValidateEmail(t *testing.T) {
	if err := ValidateEmail("a@b.online"); err != nil {
		t.Errorf("se esperaba correo válido, obtuve %v", err)
	}
	if err := ValidateEmail("a@b"); !errors.Is(err, ErrInvalidEmail) {
		t.Errorf("se esperaba ErrInvalidEmail para a@b, obtuve %v", err)
	}
}
