package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
)

type createPaqueteAtraccionRequest struct {
	AtraccionID            uint `json:"atraccion_id"`
	DiaNumero              *int `json:"dia_numero"`
	OrdenVisita            int  `json:"orden_visita"`
	DuracionEstimadaHoras  *int `json:"duracion_estimada_horas"`
}

type updatePaqueteAtraccionRequest struct {
	DiaNumero             *int `json:"dia_numero"`
	OrdenVisita           *int `json:"orden_visita"`
	DuracionEstimadaHoras *int `json:"duracion_estimada_horas"`
}

func (h *AgenciaHandler) GetPaqueteAtracciones(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, _ := strconv.ParseUint(vars["id"], 10, 32)
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}
	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var paquete models.PaqueteTuristico
	if err := database.GetDB().Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	var items []models.PaqueteAtraccion
	if err := database.GetDB().
		Preload("Atraccion").
		Where("paquete_id = ?", paqueteID).
		Order("dia_numero asc nulls first").
		Order("orden_visita asc").
		Find(&items).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener atracciones del paquete", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, items, "Atracciones del paquete obtenidas exitosamente", http.StatusOK)
}

func (h *AgenciaHandler) AddPaqueteAtraccion(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, _ := strconv.ParseUint(vars["id"], 10, 32)
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}
	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var paquete models.PaqueteTuristico
	if err := database.GetDB().Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	var req createPaqueteAtraccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	if req.AtraccionID == 0 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "atraccion_id es obligatorio", nil, http.StatusBadRequest)
		return
	}
	if req.OrdenVisita < 1 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "orden_visita debe ser >= 1", nil, http.StatusBadRequest)
		return
	}
	if req.DuracionEstimadaHoras != nil && *req.DuracionEstimadaHoras < 1 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "duracion_estimada_horas debe ser >= 1", nil, http.StatusBadRequest)
		return
	}

	// Verificar atracción existe
	var atraccion models.AtraccionTuristica
	if err := database.GetDB().First(&atraccion, req.AtraccionID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Atracción no encontrada", nil, http.StatusNotFound)
		return
	}

	dias := paqueteDuracionDiasValue(&paquete)
	isMultiDay := dias > 1
	if isMultiDay {
		if req.DiaNumero == nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "dia_numero es obligatorio para paquetes de varios días", nil, http.StatusBadRequest)
			return
		}
		if *req.DiaNumero < 1 || *req.DiaNumero > dias {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "dia_numero fuera de rango", nil, http.StatusBadRequest)
			return
		}
	} else {
		if req.DiaNumero != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "Paquetes de un día no deben tener dia_numero", nil, http.StatusBadRequest)
			return
		}
	}

	// Evitar duplicados (especialmente cuando dia_numero es NULL)
	var existing int64
	q := database.GetDB().Model(&models.PaqueteAtraccion{}).
		Where("paquete_id = ? AND atraccion_id = ?", paqueteID, req.AtraccionID)
	if req.DiaNumero == nil {
		q = q.Where("dia_numero IS NULL")
	} else {
		q = q.Where("dia_numero = ?", *req.DiaNumero)
	}
	q.Count(&existing)
	if existing > 0 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "La atracción ya está agregada para ese día", nil, http.StatusBadRequest)
		return
	}

	item := models.PaqueteAtraccion{
		PaqueteID:             paqueteID,
		AtraccionID:           req.AtraccionID,
		DiaNumero:             req.DiaNumero,
		OrdenVisita:           req.OrdenVisita,
		DuracionEstimadaHoras: req.DuracionEstimadaHoras,
	}

	if err := database.GetDB().Create(&item).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al agregar atracción al paquete", err.Error(), http.StatusInternalServerError)
		return
	}

	database.GetDB().Preload("Atraccion").First(&item, item.ID)
	utils.SuccessResponse(w, item, "Atracción agregada exitosamente", http.StatusCreated)
}

func (h *AgenciaHandler) UpdatePaqueteAtraccion(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, _ := strconv.ParseUint(vars["id"], 10, 32)
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	itemID64, err := strconv.ParseUint(vars["paquete_atraccion_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)
	itemID := uint(itemID64)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}
	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var paquete models.PaqueteTuristico
	if err := database.GetDB().Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	db := database.GetDB()
	var item models.PaqueteAtraccion
	if err := db.Where("id = ? AND paquete_id = ?", itemID, paqueteID).First(&item).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Registro no encontrado", nil, http.StatusNotFound)
		return
	}

	var req updatePaqueteAtraccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	dias := paqueteDuracionDiasValue(&paquete)
	isMultiDay := dias > 1
	if req.DiaNumero != nil {
		if isMultiDay {
			if *req.DiaNumero < 1 || *req.DiaNumero > dias {
				utils.ErrorResponse(w, "VALIDATION_ERROR", "dia_numero fuera de rango", nil, http.StatusBadRequest)
				return
			}
			item.DiaNumero = req.DiaNumero
		} else {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "Paquetes de un día no deben tener dia_numero", nil, http.StatusBadRequest)
			return
		}
	}

	if req.OrdenVisita != nil {
		if *req.OrdenVisita < 1 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "orden_visita debe ser >= 1", nil, http.StatusBadRequest)
			return
		}
		item.OrdenVisita = *req.OrdenVisita
	}

	if req.DuracionEstimadaHoras != nil {
		if *req.DuracionEstimadaHoras < 1 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "duracion_estimada_horas debe ser >= 1", nil, http.StatusBadRequest)
			return
		}
		item.DuracionEstimadaHoras = req.DuracionEstimadaHoras
	}

	if err := db.Save(&item).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar registro", err.Error(), http.StatusInternalServerError)
		return
	}

	db.Preload("Atraccion").First(&item, item.ID)
	utils.SuccessResponse(w, item, "Actualizado exitosamente", http.StatusOK)
}

func (h *AgenciaHandler) RemovePaqueteAtraccion(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, _ := strconv.ParseUint(vars["id"], 10, 32)
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	itemID64, err := strconv.ParseUint(vars["paquete_atraccion_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)
	itemID := uint(itemID64)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}
	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var paquete models.PaqueteTuristico
	if err := database.GetDB().Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	if err := database.GetDB().Where("id = ? AND paquete_id = ?", itemID, paqueteID).Delete(&models.PaqueteAtraccion{}).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al eliminar", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Eliminado exitosamente", http.StatusOK)
}

