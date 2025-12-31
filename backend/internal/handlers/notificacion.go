package handlers

import (
	"net/http"
	"strconv"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/internal/services"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
)

// NotificacionHandler maneja las peticiones HTTP relacionadas con notificaciones
type NotificacionHandler struct {
	service *services.NotificacionService
}

// NewNotificacionHandler crea un nuevo handler de notificaciones
func NewNotificacionHandler() *NotificacionHandler {
	return &NotificacionHandler{
		service: services.NewNotificacionService(database.GetDB()),
	}
}

// GetNotificaciones obtiene las notificaciones del usuario autenticado
// GET /api/v1/notificaciones?page=1&limit=10
func (h *NotificacionHandler) GetNotificaciones(w http.ResponseWriter, r *http.Request) {
	// Obtener claims del contexto
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	// Parámetros de paginación
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Obtener notificaciones
	notificaciones, total, err := h.service.ObtenerNotificaciones(claims.UserID, page, limit)
	if err != nil {
		utils.ErrorResponse(w, "DATABASE_ERROR", "Error al obtener notificaciones", nil, http.StatusInternalServerError)
		return
	}

	// Contar no leídas
	noLeidas, _ := h.service.ContarNoLeidas(claims.UserID)

	// Preparar respuesta
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	response := models.NotificacionesResponse{
		Notificaciones: models.ToNotificacionDTOs(notificaciones),
		NoLeidas:       noLeidas,
		Pagination: &models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}

	utils.SuccessResponse(w, response, "Notificaciones obtenidas exitosamente", http.StatusOK)
}

// GetContadorNoLeidas obtiene el contador de notificaciones no leídas
// GET /api/v1/notificaciones/no-leidas/count
func (h *NotificacionHandler) GetContadorNoLeidas(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	count, err := h.service.ContarNoLeidas(claims.UserID)
	if err != nil {
		utils.ErrorResponse(w, "DATABASE_ERROR", "Error al contar notificaciones", nil, http.StatusInternalServerError)
		return
	}

	response := models.ContadorNotificacionesResponse{
		NoLeidas: count,
	}

	utils.SuccessResponse(w, response, "", http.StatusOK)
}

// MarcarComoLeida marca una notificación como leída
// PUT /api/v1/notificaciones/{id}/marcar-leida
func (h *NotificacionHandler) MarcarComoLeida(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	if err := h.service.MarcarComoLeida(uint(id), claims.UserID); err != nil {
		utils.ErrorResponse(w, "UPDATE_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, nil, "Notificación marcada como leída", http.StatusOK)
}

// MarcarTodasLeidas marca todas las notificaciones del usuario como leídas
// PUT /api/v1/notificaciones/marcar-todas-leidas
func (h *NotificacionHandler) MarcarTodasLeidas(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	count, err := h.service.MarcarTodasLeidas(claims.UserID)
	if err != nil {
		utils.ErrorResponse(w, "UPDATE_ERROR", "Error al marcar notificaciones", nil, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"marcadas": count,
	}

	utils.SuccessResponse(w, response, "Todas las notificaciones marcadas como leídas", http.StatusOK)
}

// EliminarNotificacion elimina una notificación
// DELETE /api/v1/notificaciones/{id}
func (h *NotificacionHandler) EliminarNotificacion(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	if err := h.service.EliminarNotificacion(uint(id), claims.UserID); err != nil {
		utils.ErrorResponse(w, "DELETE_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, nil, "Notificación eliminada", http.StatusOK)
}
