package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/internal/services"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
)

type SalidaHandler struct {
	salidaService *services.SalidaService
}

func NewSalidaHandler() *SalidaHandler {
	return &SalidaHandler{
		salidaService: services.NewSalidaService(database.GetDB()),
	}
}

func getAgenciaIDForEncargado(w http.ResponseWriter, claims *utils.JWTClaims) (uint, bool) {
	if claims.Rol != "encargado_agencia" {
		utils.ErrorResponse(w, "FORBIDDEN", "No autorizado", nil, http.StatusForbidden)
		return 0, false
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().Where("encargado_principal_id = ?", claims.UserID).First(&agencia).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "No tiene una agencia asignada", nil, http.StatusNotFound)
		return 0, false
	}

	return agencia.ID, true
}

// CrearSalidaManual crea una nueva salida habilitada manualmente por la agencia
func (h *SalidaHandler) CrearSalidaManual(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	agenciaID, ok := getAgenciaIDForEncargado(w, claims)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	paqueteID, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de paquete invalido", nil, http.StatusBadRequest)
		return
	}

	var req services.CrearSalidaManualRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "Datos invalidos", err.Error(), http.StatusBadRequest)
		return
	}

	salida, err := h.salidaService.CrearSalidaManual(agenciaID, claims.UserID, uint(paqueteID), req)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, salida, "Salida creada exitosamente", http.StatusCreated)
}

// ObtenerSalidasPorPaquete obtiene todas las salidas habilitadas de un paquete
func (h *SalidaHandler) ObtenerSalidasPorPaquete(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	agenciaID, ok := getAgenciaIDForEncargado(w, claims)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	paqueteID, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de paquete invalido", nil, http.StatusBadRequest)
		return
	}

	filtros := services.SalidaFiltros{
		Estado:       r.URL.Query().Get("estado"),
		FechaDesde:   r.URL.Query().Get("fecha_desde"),
		FechaHasta:   r.URL.Query().Get("fecha_hasta"),
		SoloManuales: r.URL.Query().Get("solo_manuales") == "true",
	}

	salidas, err := h.salidaService.ObtenerSalidasPorPaquete(agenciaID, uint(paqueteID), filtros)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, salidas, "Salidas obtenidas", http.StatusOK)
}

// ActualizarSalida actualiza los datos de una salida existente
func (h *SalidaHandler) ActualizarSalida(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	agenciaID, ok := getAgenciaIDForEncargado(w, claims)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	salidaID, err := strconv.ParseUint(vars["salida_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de salida invalido", nil, http.StatusBadRequest)
		return
	}

	var req services.ActualizarSalidaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "Datos invalidos", err.Error(), http.StatusBadRequest)
		return
	}

	salida, err := h.salidaService.ActualizarSalida(agenciaID, uint(salidaID), req)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, salida, "Salida actualizada", http.StatusOK)
}

// CancelarSalida cancela una salida existente
func (h *SalidaHandler) CancelarSalida(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	agenciaID, ok := getAgenciaIDForEncargado(w, claims)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	salidaID, err := strconv.ParseUint(vars["salida_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de salida invalido", nil, http.StatusBadRequest)
		return
	}

	var body struct {
		Razon string `json:"razon"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Razon == "" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Debe proporcionar una razon de cancelacion", nil, http.StatusBadRequest)
		return
	}

	if err := h.salidaService.CancelarSalida(agenciaID, uint(salidaID), body.Razon); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, nil, "Salida cancelada exitosamente", http.StatusOK)
}

// ObtenerSalidasPublicas obtiene salidas disponibles para turistas (públicas)
func (h *SalidaHandler) ObtenerSalidasPublicas(w http.ResponseWriter, r *http.Request) {
	var salidas []models.PaqueteSalidaHabilitada

	query := database.GetDB().Model(&models.PaqueteSalidaHabilitada{}).
		Where("estado IN ('pendiente', 'activa')").
		Where("creada_manualmente = ?", true). // Solo salidas creadas manualmente
		Where("fecha_salida >= CURRENT_DATE"). // Solo salidas futuras
		Order("fecha_salida ASC")

	// Filtros opcionales
	if paqueteID := r.URL.Query().Get("paquete_id"); paqueteID != "" {
		query = query.Where("paquete_id = ?", paqueteID)
	}

	if fechaDesde := r.URL.Query().Get("fecha_desde"); fechaDesde != "" {
		query = query.Where("fecha_salida >= ?", fechaDesde)
	}

	if fechaHasta := r.URL.Query().Get("fecha_hasta"); fechaHasta != "" {
		query = query.Where("fecha_salida <= ?", fechaHasta)
	}

	if err := query.Find(&salidas).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener salidas", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, salidas, "Salidas disponibles", http.StatusOK)
}
