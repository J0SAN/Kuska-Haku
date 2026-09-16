package http

import (
	"context"
	"net/http"
	"strings"

	// Ajusta el nombre de tu módulo
	"kuska-haku-backend/internal/ports"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	authSvc ports.AuthService
}

func NewAuthMiddleware(authSvc ports.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authSvc: authSvc}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Se requiere cabecera de autorización", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de token inválido. Use: Bearer <token>", http.StatusUnauthorized)
			return
		}

		user, err := m.authSvc.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
