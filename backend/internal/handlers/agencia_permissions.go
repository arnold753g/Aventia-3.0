package handlers

import (
	"net/http"

	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"
)

func getClaimsOrUnauthorized(w http.ResponseWriter, r *http.Request) (*utils.JWTClaims, bool) {
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return nil, false
	}
	return claims, true
}

func canManageAgencia(claims *utils.JWTClaims, agencia *models.AgenciaTurismo) bool {
	if claims == nil || agencia == nil {
		return false
	}
	if claims.Rol == "admin" {
		return true
	}
	if claims.Rol != "encargado_agencia" {
		return false
	}
	if agencia.EncargadoPrincipalID == nil {
		return false
	}
	return *agencia.EncargadoPrincipalID == claims.UserID
}

