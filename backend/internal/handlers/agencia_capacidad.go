package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type agenciaCapacidadUpdateRequest struct {
	MaxSalidasPorDia     *int `json:"max_salidas_por_dia"`
	MaxSalidasPorHorario *int `json:"max_salidas_por_horario"`
}

func ensureAgenciaCapacidadRow(db *gorm.DB, agenciaID uint) (*models.AgenciaCapacidad, error) {
	var capacidad models.AgenciaCapacidad
	err := db.Where("agencia_id = ?", agenciaID).First(&capacidad).Error
	if err == nil {
		return &capacidad, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	capacidad = models.AgenciaCapacidad{
		AgenciaID:            agenciaID,
		MaxSalidasPorDia:     5,
		MaxSalidasPorHorario: 3,
	}

	if err := db.Create(&capacidad).Error; err != nil {
		return nil, err
	}

	return &capacidad, nil
}

// GetAgenciaCapacidad obtiene (y crea si no existe) la capacidad operativa configurada por la agencia.
func (h *AgenciaHandler) GetAgenciaCapacidad(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	capacidad, err := ensureAgenciaCapacidadRow(database.GetDB(), uint(id))
	if err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener capacidad", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, capacidad, "Capacidad obtenida exitosamente", http.StatusOK)
}

// UpdateAgenciaCapacidad actualiza la capacidad operativa configurada por la agencia (admin o encargado).
func (h *AgenciaHandler) UpdateAgenciaCapacidad(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var req agenciaCapacidadUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if req.MaxSalidasPorDia != nil && *req.MaxSalidasPorDia < 0 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El maximo de salidas por dia no puede ser negativo", nil, http.StatusBadRequest)
		return
	}

	if req.MaxSalidasPorHorario != nil && *req.MaxSalidasPorHorario < 0 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El maximo de salidas por horario no puede ser negativo", nil, http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	capacidad, err := ensureAgenciaCapacidadRow(db, uint(id))
	if err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al preparar capacidad", err.Error(), http.StatusInternalServerError)
		return
	}

	if req.MaxSalidasPorDia != nil {
		capacidad.MaxSalidasPorDia = *req.MaxSalidasPorDia
	}

	if req.MaxSalidasPorHorario != nil {
		capacidad.MaxSalidasPorHorario = *req.MaxSalidasPorHorario
	}

	if capacidad.MaxSalidasPorHorario > capacidad.MaxSalidasPorDia {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El maximo por horario no puede ser mayor al maximo por dia", nil, http.StatusBadRequest)
		return
	}

	if err := db.Save(capacidad).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar capacidad", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, capacidad, "Capacidad actualizada exitosamente", http.StatusOK)
}

