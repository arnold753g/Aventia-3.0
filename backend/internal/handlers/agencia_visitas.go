package handlers

import (
	"net/http"
	"strconv"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
)

type AgenciaVisitasHandler struct{}

func NewAgenciaVisitasHandler() *AgenciaVisitasHandler {
	return &AgenciaVisitasHandler{}
}

// RegistrarVisita registra una visita a la página pública de una agencia
func (h *AgenciaVisitasHandler) RegistrarVisita(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrSlug := vars["id"]

	// Buscar agencia por ID o slug
	agencia, err := resolveAgenciaByIDOrSlug(database.GetDB(), idOrSlug)
	if err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	// Obtener IP del visitante
	ipAddress := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ipAddress = forwarded
	}

	// Obtener User Agent
	userAgent := r.UserAgent()

	// Obtener Referer
	referer := r.Referer()

	// Crear registro de visita
	visita := models.AgenciaVisita{
		AgenciaID:   agencia.ID,
		FechaVisita: time.Now(),
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Referer:     referer,
	}

	if err := database.GetDB().Create(&visita).Error; err != nil {
		// No fallar la request si no se puede registrar la visita
		// Solo log el error silenciosamente
		utils.SuccessResponse(w, map[string]string{"status": "ok"}, "Visita registrada", http.StatusOK)
		return
	}

	utils.SuccessResponse(w, map[string]string{"status": "ok"}, "Visita registrada exitosamente", http.StatusCreated)
}

// GetEstadisticasVisitas obtiene las estadísticas de visitas de una agencia
func (h *AgenciaVisitasHandler) GetEstadisticasVisitas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	// Verificar que la agencia existe
	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	// Verificar permisos (solo el encargado de la agencia o admin puede ver las estadísticas)
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	if claims.Rol != "admin" {
		if claims.Rol != "encargado_agencia" || agencia.EncargadoPrincipalID == nil || *agencia.EncargadoPrincipalID != claims.UserID {
			utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para ver estas estadísticas", nil, http.StatusForbidden)
			return
		}
	}

	// Obtener estadísticas desde la vista
	var estadisticas models.AgenciaEstadisticasVisitas
	err = database.GetDB().
		Table("agencia_estadisticas_visitas").
		Where("agencia_id = ?", id).
		First(&estadisticas).Error

	if err != nil {
		// Si no hay datos en la vista, devolver estadísticas en cero
		estadisticas = models.AgenciaEstadisticasVisitas{
			AgenciaID:           agencia.ID,
			NombreComercial:     agencia.NombreComercial,
			Slug:                agencia.Slug,
			TotalVisitas:        0,
			DiasConVisitas:      0,
			UltimaVisita:        nil,
			VisitasUltimaSemana: 0,
			VisitasUltimoMes:    0,
		}
	}

	utils.SuccessResponse(w, estadisticas, "Estadísticas obtenidas exitosamente", http.StatusOK)
}

// GetVisitasDetalle obtiene el detalle de visitas de una agencia con paginación
func (h *AgenciaVisitasHandler) GetVisitasDetalle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	// Verificar que la agencia existe
	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	// Verificar permisos
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	if claims.Rol != "admin" {
		if claims.Rol != "encargado_agencia" || agencia.EncargadoPrincipalID == nil || *agencia.EncargadoPrincipalID != claims.UserID {
			utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para ver estas visitas", nil, http.StatusForbidden)
			return
		}
	}

	// Paginación
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Contar total
	var total int64
	database.GetDB().Model(&models.AgenciaVisita{}).Where("agencia_id = ?", id).Count(&total)

	// Obtener visitas
	var visitas []models.AgenciaVisita
	offset := (page - 1) * limit
	err = database.GetDB().
		Where("agencia_id = ?", id).
		Order("fecha_visita DESC").
		Limit(limit).
		Offset(offset).
		Find(&visitas).Error

	if err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener visitas", err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"visitas": visitas,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utils.SuccessResponse(w, response, "Visitas obtenidas exitosamente", http.StatusOK)
}
