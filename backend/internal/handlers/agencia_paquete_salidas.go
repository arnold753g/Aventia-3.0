package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
)

var allowedSalidaEstado = map[string]bool{
	"pendiente":  true,
	"activa":     true,
	"completada": true,
	"cancelada":  true,
}

type updateSalidaRequest struct {
	PuntoEncuentro        *string `json:"punto_encuentro"`
	HoraEncuentro         *string `json:"hora_encuentro"`
	NotasLogistica        *string `json:"notas_logistica"`
	InstruccionesTuristas *string `json:"instrucciones_turistas"`
	GuiaNombre            *string `json:"guia_nombre"`
	GuiaTelefono          *string `json:"guia_telefono"`
	Estado                *string `json:"estado"`
	RazonCancelacion      *string `json:"razon_cancelacion"`
}

type createSalidaRequest struct {
	FechaSalida string  `json:"fecha_salida"`
	TipoSalida  *string `json:"tipo_salida"`
	Estado      *string `json:"estado"`
}

// CreatePaqueteSalida habilita una salida manualmente para una fecha especifica.
func (h *AgenciaHandler) CreatePaqueteSalida(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, _ := strconv.ParseUint(vars["id"], 10, 32)
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de paquete invalido", nil, http.StatusBadRequest)
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
	if paquete.Status != "activo" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El paquete no esta activo", nil, http.StatusBadRequest)
		return
	}

	var req createSalidaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.FechaSalida) == "" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "fecha_salida es obligatoria", nil, http.StatusBadRequest)
		return
	}

	loc := time.Now().Location()
	fecha, err := time.ParseInLocation("2006-01-02", req.FechaSalida, loc)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "fecha_salida invalida (YYYY-MM-DD)", nil, http.StatusBadRequest)
		return
	}
	fechaStr := fecha.Format("2006-01-02")

	tipo := "compartido"
	if req.TipoSalida != nil && strings.TrimSpace(*req.TipoSalida) != "" {
		tipo = strings.ToLower(strings.TrimSpace(*req.TipoSalida))
	}
	if tipo != "compartido" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "tipo_salida invalido (solo compartido)", nil, http.StatusBadRequest)
		return
	}

	estado := "pendiente"
	if req.Estado != nil && strings.TrimSpace(*req.Estado) != "" {
		estado = strings.ToLower(strings.TrimSpace(*req.Estado))
	}
	if estado != "pendiente" && estado != "activa" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "estado invalido (pendiente|activa)", nil, http.StatusBadRequest)
		return
	}

	if paquete.Frecuencia == "salida_unica" {
		if paquete.FechaSalidaFija == nil || strings.TrimSpace(*paquete.FechaSalidaFija) == "" {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "El paquete no tiene fecha de salida fija", nil, http.StatusBadRequest)
			return
		}
		if fechaStr != *paquete.FechaSalidaFija {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "La fecha debe coincidir con la salida unica del paquete", nil, http.StatusBadRequest)
			return
		}
	}

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, loc)
	minDate := today.AddDate(0, 0, paquete.DiasPreviosCompra)
	if fecha.Before(minDate) {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "La fecha seleccionada no cumple los dias previos de compra", nil, http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	var existentes int64
	if err := db.Model(&models.PaqueteSalidaHabilitada{}).
		Where("paquete_id = ? AND fecha_salida = ? AND tipo_salida = ? AND estado <> 'cancelada'", paquete.ID, fechaStr, tipo).
		Count(&existentes).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al validar salidas existentes", err.Error(), http.StatusInternalServerError)
		return
	}
	if existentes > 0 {
		utils.ErrorResponse(w, "DUPLICATE", "Ya existe una salida para esa fecha", nil, http.StatusBadRequest)
		return
	}

	maxSalidasDia := 5
	maxSalidasHorario := 3
	var capacidad models.AgenciaCapacidad
	if err := db.Where("agencia_id = ?", agencia.ID).First(&capacidad).Error; err == nil {
		if capacidad.MaxSalidasPorDia > 0 {
			maxSalidasDia = capacidad.MaxSalidasPorDia
		}
		if capacidad.MaxSalidasPorHorario > 0 {
			maxSalidasHorario = capacidad.MaxSalidasPorHorario
		}
	}

	horarioCapacidad := "todo_dia"
	if paquete.DuracionDias == nil || *paquete.DuracionDias <= 1 {
		if paquete.Horario != nil && strings.TrimSpace(*paquete.Horario) != "" {
			horarioCapacidad = strings.TrimSpace(*paquete.Horario)
		}
	}

	var salidasDia int64
	if err := db.Table("paquete_salidas_habilitadas s").
		Joins("JOIN paquetes_turisticos p ON p.id = s.paquete_id").
		Where("p.agencia_id = ? AND s.fecha_salida = ? AND s.estado IN ('pendiente', 'activa')", agencia.ID, fechaStr).
		Count(&salidasDia).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al validar capacidad diaria", err.Error(), http.StatusInternalServerError)
		return
	}
	if salidasDia >= int64(maxSalidasDia) {
		utils.ErrorResponse(w, "CAPACITY_LIMIT", "La agencia alcanzo su maximo de salidas para ese dia", nil, http.StatusBadRequest)
		return
	}

	var salidasHorario int64
	if err := db.Table("paquete_salidas_habilitadas s").
		Joins("JOIN paquetes_turisticos p ON p.id = s.paquete_id").
		Where("p.agencia_id = ? AND s.fecha_salida = ? AND s.estado IN ('pendiente', 'activa') AND COALESCE(p.horario, 'todo_dia') = ?", agencia.ID, fechaStr, horarioCapacidad).
		Count(&salidasHorario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al validar capacidad por horario", err.Error(), http.StatusInternalServerError)
		return
	}
	if salidasHorario >= int64(maxSalidasHorario) {
		utils.ErrorResponse(w, "CAPACITY_LIMIT", "La agencia alcanzo su maximo de salidas simultaneas", nil, http.StatusBadRequest)
		return
	}

	salida := models.PaqueteSalidaHabilitada{
		PaqueteID:  paquete.ID,
		FechaSalida: fechaStr,
		TipoSalida:  tipo,
		CupoMinimo:  paquete.CupoMinimo,
		CupoMaximo:  paquete.CupoMaximo,
		Estado:      estado,
	}

	if err := db.Create(&salida).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al habilitar salida", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, salida, "Salida habilitada exitosamente", http.StatusCreated)
}

func (h *AgenciaHandler) GetPaqueteSalidas(w http.ResponseWriter, r *http.Request) {
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

	var salidas []models.PaqueteSalidaHabilitada
	if err := database.GetDB().
		Where("paquete_id = ?", paqueteID).
		Order("fecha_salida asc").
		Order("tipo_salida asc").
		Find(&salidas).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener salidas", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, salidas, "Salidas obtenidas exitosamente", http.StatusOK)
}

func (h *AgenciaHandler) UpdatePaqueteSalida(w http.ResponseWriter, r *http.Request) {
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
	salidaID64, err := strconv.ParseUint(vars["salida_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de salida inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)
	salidaID := uint(salidaID64)

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
	var salida models.PaqueteSalidaHabilitada
	if err := db.Where("id = ? AND paquete_id = ?", salidaID, paqueteID).First(&salida).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Salida no encontrada", nil, http.StatusNotFound)
		return
	}

	var req updateSalidaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	if req.PuntoEncuentro != nil {
		salida.PuntoEncuentro = normalizeStringPtr(req.PuntoEncuentro)
	}
	if req.HoraEncuentro != nil {
		salida.HoraEncuentro = normalizeStringPtr(req.HoraEncuentro)
	}
	if req.NotasLogistica != nil {
		salida.NotasLogistica = normalizeStringPtr(req.NotasLogistica)
	}
	if req.InstruccionesTuristas != nil {
		salida.InstruccionesTuristas = normalizeStringPtr(req.InstruccionesTuristas)
	}
	if req.GuiaNombre != nil {
		salida.GuiaNombre = normalizeStringPtr(req.GuiaNombre)
	}
	if req.GuiaTelefono != nil {
		salida.GuiaTelefono = normalizeStringPtr(req.GuiaTelefono)
	}

	if req.Estado != nil {
		estado := strings.TrimSpace(*req.Estado)
		if !allowedSalidaEstado[estado] {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "Estado inválido", nil, http.StatusBadRequest)
			return
		}
		salida.Estado = estado
	}

	if salida.Estado == "cancelada" {
		rc := normalizeStringPtr(req.RazonCancelacion)
		if rc == nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "Debe especificar razon_cancelacion al cancelar", nil, http.StatusBadRequest)
			return
		}
		salida.RazonCancelacion = rc
	} else if req.RazonCancelacion != nil {
		salida.RazonCancelacion = nil
	}

	if err := db.Save(&salida).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar salida", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, salida, "Salida actualizada exitosamente", http.StatusOK)
}

// ActivarSalida activa manualmente una salida (endpoint específico para encargados)
func (h *AgenciaHandler) ActivarSalida(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, _ := strconv.ParseUint(vars["id"], 10, 32)
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de paquete inválido", nil, http.StatusBadRequest)
		return
	}
	salidaID64, err := strconv.ParseUint(vars["salida_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de salida inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)
	salidaID := uint(salidaID64)

	// Validar permisos
	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}
	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	// Validar paquete
	var paquete models.PaqueteTuristico
	if err := database.GetDB().Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	// Obtener salida
	db := database.GetDB()
	var salida models.PaqueteSalidaHabilitada
	if err := db.Where("id = ? AND paquete_id = ?", salidaID, paqueteID).First(&salida).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Salida no encontrada", nil, http.StatusNotFound)
		return
	}

	// Validar estado actual
	if salida.Estado == "activa" {
		utils.ErrorResponse(w, "ALREADY_ACTIVE", "La salida ya está activa", nil, http.StatusBadRequest)
		return
	}
	if salida.Estado == "cancelada" {
		utils.ErrorResponse(w, "CANCELLED", "No se puede activar una salida cancelada", nil, http.StatusBadRequest)
		return
	}
	if salida.Estado == "completada" {
		utils.ErrorResponse(w, "COMPLETED", "No se puede activar una salida completada", nil, http.StatusBadRequest)
		return
	}

	// Validar cupo mínimo (opcional: el encargado puede forzar la activación)
	var req struct {
		ForzarActivacion bool `json:"forzar_activacion"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	cuposActuales := salida.CuposConfirmados + salida.CuposReservados
	if cuposActuales < salida.CupoMinimo && !req.ForzarActivacion {
		utils.ErrorResponse(w, "CUPO_MINIMO_NO_ALCANZADO",
			"No se ha alcanzado el cupo mínimo. Envíe forzar_activacion=true para activar de todos modos",
			map[string]interface{}{
				"cupo_minimo":      salida.CupoMinimo,
				"cupos_actuales":   cuposActuales,
				"cupos_reservados": salida.CuposReservados,
				"cupos_confirmados": salida.CuposConfirmados,
			}, http.StatusBadRequest)
		return
	}

	// Activar la salida
	salida.Estado = "activa"
	if err := db.Save(&salida).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al activar salida", err.Error(), http.StatusInternalServerError)
		return
	}

	mensaje := "Salida activada exitosamente"
	if cuposActuales < salida.CupoMinimo {
		mensaje = "Salida activada (forzada sin alcanzar cupo mínimo)"
	}

	utils.SuccessResponse(w, salida, mensaje, http.StatusOK)
}
