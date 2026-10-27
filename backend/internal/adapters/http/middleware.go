package http

import (
	"context"
	"net/http"
	"strings"

	"kuska-haku-backend/internal/domain"
	"kuska-haku-backend/internal/ports"
)

type contextKey string

const claimsContextKey contextKey = "claims"

// ClaimsFromContext extrae la identidad verificada del usuario desde el contexto del request.
func ClaimsFromContext(ctx context.Context) (*ports.TokenClaims, bool) {
	c, ok := ctx.Value(claimsContextKey).(*ports.TokenClaims)
	return c, ok
}

type AuthMiddleware struct {
	tokens ports.TokenService
}

func NewAuthMiddleware(tokens ports.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokens: tokens}
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

		claims, err := m.tokens.Parse(parts[1])
		if err != nil {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RequireRole verifica que el usuario autenticado posea uno de los roles permitidos (RBAC).
func (m *AuthMiddleware) RequireRole(next http.HandlerFunc, allowed ...domain.Role) http.HandlerFunc {
	return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok {
			http.Error(w, "No autenticado", http.StatusUnauthorized)
			return
		}

		for _, role := range allowed {
			if claims.Role == role {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "No tienes permiso para acceder a este recurso", http.StatusForbidden)
	})
}
