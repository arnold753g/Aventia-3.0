package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type paquetePoliticasUpdateRequest struct {
	EdadMinimaPago          *int     `json:"edad_minima_pago"`
	RecargoPrivadoPorcentaje *float64 `json:"recargo_privado_porcentaje"`
	PoliticaCancelacion     *string  `json:"politica_cancelacion"`
}

func ensurePaquetePoliticasRow(db *gorm.DB, agenciaID uint) (*models.PaquetePolitica, error) {
	var politica models.PaquetePolitica
	err := db.Where("agencia_id = ?", agenciaID).First(&politica).Error
	if err == nil {
		return &politica, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	politica = models.PaquetePolitica{
		AgenciaID:                agenciaID,
		EdadMinimaPago:           6,
		RecargoPrivadoPorcentaje: 0,
		PoliticaCancelacion:      nil,
	}

	if err := db.Create(&politica).Error; err != nil {
		return nil, err
	}

	return &politica, nil
}

// GetPaquetePoliticas obtiene (y crea si no existe) la configuración de políticas de paquetes por agencia.
func (h *AgenciaHandler) GetPaquetePoliticas(w http.ResponseWriter, r *http.Request) {
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

	politica, err := ensurePaquetePoliticasRow(database.GetDB(), uint(id))
	if err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener politicas", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, politica, "Politicas obtenidas exitosamente", http.StatusOK)
}

// UpdatePaquetePoliticas actualiza la configuración de políticas de paquetes por agencia (admin o encargado).
func (h *AgenciaHandler) UpdatePaquetePoliticas(w http.ResponseWriter, r *http.Request) {
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

	var req paquetePoliticasUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if req.EdadMinimaPago != nil && *req.EdadMinimaPago < 0 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "La edad minima no puede ser negativa", nil, http.StatusBadRequest)
		return
	}

	if req.RecargoPrivadoPorcentaje != nil && (*req.RecargoPrivadoPorcentaje < 0 || *req.RecargoPrivadoPorcentaje > 100) {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El recargo debe estar entre 0 y 100", nil, http.StatusBadRequest)
		return
	}

	if req.PoliticaCancelacion != nil && len(*req.PoliticaCancelacion) > 5000 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "La politica de cancelacion es demasiado larga", nil, http.StatusBadRequest)
		return
	}

	db := database.GetDB()

	politica, err := ensurePaquetePoliticasRow(db, uint(id))
	if err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al preparar politicas", err.Error(), http.StatusInternalServerError)
		return
	}

	if req.EdadMinimaPago != nil {
		politica.EdadMinimaPago = *req.EdadMinimaPago
	}

	if req.RecargoPrivadoPorcentaje != nil {
		politica.RecargoPrivadoPorcentaje = *req.RecargoPrivadoPorcentaje
	}

	if req.PoliticaCancelacion != nil {
		trimmed := strings.TrimSpace(*req.PoliticaCancelacion)
		if trimmed == "" {
			politica.PoliticaCancelacion = nil
		} else {
			politica.PoliticaCancelacion = &trimmed
		}
	}

	if err := db.Save(politica).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar politicas", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, politica, "Politicas actualizadas exitosamente", http.StatusOK)
}

