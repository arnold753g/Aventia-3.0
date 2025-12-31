package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
)

type createItinerarioRequest struct {
	DiaNumero     int      `json:"dia_numero"`
	Titulo        string   `json:"titulo"`
	Descripcion   *string  `json:"descripcion"`
	Actividades   []string `json:"actividades"`
	HospedajeInfo *string  `json:"hospedaje_info"`
}

type updateItinerarioRequest struct {
	DiaNumero     *int      `json:"dia_numero"`
	Titulo        *string   `json:"titulo"`
	Descripcion   *string   `json:"descripcion"`
	Actividades   *[]string `json:"actividades"`
	HospedajeInfo *string   `json:"hospedaje_info"`
}

func (h *AgenciaHandler) GetPaqueteItinerario(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID, _ := strconv.ParseUint(vars["id"], 10, 32)
	paqueteID, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, uint(agenciaID)).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}
	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var paquete models.PaqueteTuristico
	if err := database.GetDB().Where("id = ? AND agencia_id = ?", uint(paqueteID), uint(agenciaID)).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	var items []models.PaqueteItinerario
	if err := database.GetDB().
		Where("paquete_id = ?", uint(paqueteID)).
		Order("dia_numero asc").
		Find(&items).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener itinerario", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, items, "Itinerario obtenido exitosamente", http.StatusOK)
}

func (h *AgenciaHandler) CreatePaqueteItinerario(w http.ResponseWriter, r *http.Request) {
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

	dias := paqueteDuracionDiasValue(&paquete)
	if dias <= 1 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "No se puede agregar itinerario a paquetes de un día", nil, http.StatusBadRequest)
		return
	}

	var req createItinerarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	if req.DiaNumero < 1 || req.DiaNumero > dias {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "dia_numero fuera de rango", nil, http.StatusBadRequest)
		return
	}

	titulo := strings.TrimSpace(req.Titulo)
	if titulo == "" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El título es obligatorio", nil, http.StatusBadRequest)
		return
	}

	item := models.PaqueteItinerario{
		PaqueteID:     paqueteID,
		DiaNumero:     req.DiaNumero,
		Titulo:        titulo,
		Descripcion:   normalizeStringPtr(req.Descripcion),
		Actividades:   models.StringArray(req.Actividades),
		HospedajeInfo: normalizeStringPtr(req.HospedajeInfo),
	}

	if err := database.GetDB().Create(&item).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear itinerario", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, item, "Itinerario creado exitosamente", http.StatusCreated)
}

func (h *AgenciaHandler) UpdatePaqueteItinerario(w http.ResponseWriter, r *http.Request) {
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
	itID64, err := strconv.ParseUint(vars["itinerario_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de itinerario inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)
	itID := uint(itID64)

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

	dias := paqueteDuracionDiasValue(&paquete)
	if dias <= 1 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "No se puede editar itinerario en paquetes de un día", nil, http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	var item models.PaqueteItinerario
	if err := db.Where("id = ? AND paquete_id = ?", itID, paqueteID).First(&item).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Itinerario no encontrado", nil, http.StatusNotFound)
		return
	}

	var req updateItinerarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	if req.DiaNumero != nil {
		if *req.DiaNumero < 1 || *req.DiaNumero > dias {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "dia_numero fuera de rango", nil, http.StatusBadRequest)
			return
		}
		item.DiaNumero = *req.DiaNumero
	}
	if req.Titulo != nil {
		t := strings.TrimSpace(*req.Titulo)
		if t == "" {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "El título no puede estar vacío", nil, http.StatusBadRequest)
			return
		}
		item.Titulo = t
	}
	if req.Descripcion != nil {
		item.Descripcion = normalizeStringPtr(req.Descripcion)
	}
	if req.Actividades != nil {
		item.Actividades = models.StringArray(*req.Actividades)
	}
	if req.HospedajeInfo != nil {
		item.HospedajeInfo = normalizeStringPtr(req.HospedajeInfo)
	}

	if err := db.Save(&item).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar itinerario", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, item, "Itinerario actualizado exitosamente", http.StatusOK)
}

func (h *AgenciaHandler) DeletePaqueteItinerario(w http.ResponseWriter, r *http.Request) {
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
	itID64, err := strconv.ParseUint(vars["itinerario_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de itinerario inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)
	itID := uint(itID64)

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
	if err := db.Where("id = ? AND paquete_id = ?", itID, paqueteID).Delete(&models.PaqueteItinerario{}).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al eliminar itinerario", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Itinerario eliminado exitosamente", http.StatusOK)
}
