package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		passwordHash string
		fullName     string
		role         string
		wantErr      bool
		expectedRole string
	}{
		{
			name:         "Usuario válido con rol explícito",
			email:        "admin@kuskahaku.pe",
			passwordHash: "hash_seguro_123",
			fullName:     "Ana Admin",
			role:         "admin",
			wantErr:      false,
			expectedRole: "admin",
		},
		{
			name:         "Usuario válido con rol por defecto",
			email:        "pasajero@kuskahaku.pe",
			passwordHash: "hash_seguro_123",
			fullName:     "Juan Pasajero",
			role:         "",
			wantErr:      false,
			expectedRole: "passenger",
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
			name:         "Fallo por mayúsculas en el correo (fuera de la regex)",
			email:        "CARLOS@kuskahaku.pe",
			passwordHash: "hash_seguro_123",
			fullName:     "Carlos Rojas",
			role:         "passenger",
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

			if !tt.wantErr && user.Role != tt.expectedRole {
				t.Errorf("NewUser() asignó rol = %v, se esperaba rol = %v", user.Role, tt.expectedRole)
			}
		})
	}
}
