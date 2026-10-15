package auth

import (
	"kuska-haku-backend/internal/domain"
	"testing"
)

func TestJWTService(t *testing.T) {
	jwtSvc := NewJWTService("clave-secreta-test")
	user := &domain.User{
		ID:    "user-1",
		Email: "test@kuskahaku.pe",
		Role:  "admin",
	}

	token, err := jwtSvc.GenerateToken(user)
	if err != nil || token == "" {
		t.Fatalf("Fallo al generar el token: %v", err)
	}

	validUser, err := jwtSvc.ValidateToken(token)
	if err != nil {
		t.Fatalf("Fallo al validar el token JWT: %v", err)
	}

	// Corregimos la aserción: Evaluamos el ID y el Rol, que son los datos que el JWT realmente almacena
	if validUser.ID != user.ID {
		t.Errorf("El ID decodificado del token (%s) no coincide con el original (%s)", validUser.ID, user.ID)
	}

	// Opcional: También podemos validar que conserve el rol
	if validUser.Role != user.Role {
		t.Errorf("El rol decodificado del token no coincide con el original")
	}
}
