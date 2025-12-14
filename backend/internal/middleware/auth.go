package middleware

import (
	"context"
	"net/http"
	"strings"

	"andaria-backend/pkg/utils"
)

// AuthMiddleware verifica el token JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener token del header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, "UNAUTHORIZED", "Token no proporcionado", nil, http.StatusUnauthorized)
			return
		}

		// Verificar formato "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(w, "INVALID_TOKEN_FORMAT", "Formato de token inválido", nil, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validar token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(w, "INVALID_TOKEN", "Token inválido o expirado", nil, http.StatusUnauthorized)
			return
		}

		// Agregar claims al contexto
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware verifica que el usuario tenga el rol requerido
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
			if !ok {
				utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
				return
			}

			// Verificar si el rol del usuario está en la lista de roles permitidos
			roleAllowed := false
			for _, role := range allowedRoles {
				if claims.Rol == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para realizar esta acción", nil, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
