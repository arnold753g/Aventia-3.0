package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type AgenciaHandler struct {
	validate *validator.Validate
}

func NewAgenciaHandler() *AgenciaHandler {
	return &AgenciaHandler{
		validate: validator.New(),
	}
}

// capitalizeAgenciaName capitaliza el nombre comercial
func capitalizeAgenciaName(name string) string {
	words := strings.Fields(strings.ToLower(name))
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// GetAgencias obtiene lista de agencias con filtros
func (h *AgenciaHandler) GetAgencias(w http.ResponseWriter, r *http.Request) {
	query := database.GetDB().Model(&models.AgenciaTurismo{})

	// Filtros
	if search := r.URL.Query().Get("search"); search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nombre_comercial ILIKE ? OR descripcion ILIKE ? OR direccion ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if deptID := r.URL.Query().Get("departamento_id"); deptID != "" {
		query = query.Where("departamento_id = ?", deptID)
	}

	if status := r.URL.Query().Get("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if licencia := r.URL.Query().Get("licencia_turistica"); licencia != "" {
		query = query.Where("licencia_turistica = ?", licencia == "true")
	}

	if espID := r.URL.Query().Get("especialidad_id"); espID != "" {
		query = query.Joins("JOIN agencia_especialidades ON agencia_especialidades.agencia_id = agencias_turismo.id").
			Where("agencia_especialidades.categoria_id = ?", espID)
	}

	if encID := r.URL.Query().Get("encargado_id"); encID != "" {
		query = query.Where("encargado_principal_id = ?", encID)
	}

	if visible := r.URL.Query().Get("visible_publico"); visible != "" {
		query = query.Where("visible_publico = ?", visible == "true")
	}

	// Paginación
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Ordenamiento
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := strings.ToLower(r.URL.Query().Get("sort_order"))
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	var total int64
	query.Count(&total)

	var agencias []models.AgenciaTurismo
	offset := (page - 1) * limit
	err := query.
		Preload("Departamento").
		Preload("EncargadoPrincipal").
		Preload("Fotos").
		Preload("Especialidades.Categoria").
		Preload("Dias").
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&agencias).Error

	if err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener agencias", err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"agencias": agencias,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utils.SuccessResponse(w, response, "Agencias obtenidas exitosamente", http.StatusOK)
}

// GetAgencia obtiene una agencia por ID
func (h *AgenciaHandler) GetAgencia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	err = database.GetDB().
		Preload("Departamento").
		Preload("EncargadoPrincipal").
		Preload("Fotos").
		Preload("Especialidades.Categoria").
		Preload("Dias").
		First(&agencia, id).Error

	if err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	utils.SuccessResponse(w, agencia, "Agencia obtenida exitosamente", http.StatusOK)
}

// CreateAgenciaRapida crea una agencia con información básica
func (h *AgenciaHandler) CreateAgenciaRapida(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAgenciaRapidaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", err.Error(), http.StatusBadRequest)
		return
	}

	// Validar
	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validación", err.Error(), http.StatusBadRequest)
		return
	}

	// Capitalizar nombre
	req.NombreComercial = capitalizeAgenciaName(req.NombreComercial)

	// Verificar que el encargado exista y tenga el rol correcto
	if req.EncargadoPrincipalID != nil {
		var encargado models.Usuario
		if err := database.GetDB().First(&encargado, *req.EncargadoPrincipalID).Error; err != nil {
			utils.ErrorResponse(w, "ENCARGADO_NOT_FOUND", "Encargado no encontrado", nil, http.StatusBadRequest)
			return
		}

		if encargado.Rol != "encargado_agencia" {
			utils.ErrorResponse(w, "INVALID_ROL", "El encargado debe tener el rol 'encargado_agencia'", nil, http.StatusBadRequest)
			return
		}
	}

	// Obtener usuario del contexto
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	// Solo administradores pueden crear agencias
	if claims.Rol != "admin" {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para crear agencias", nil, http.StatusForbidden)
		return
	}

	// Crear agencia con valores por defecto para campos requeridos
	agencia := models.AgenciaTurismo{
		NombreComercial:      req.NombreComercial,
		Direccion:            "Por completar",
		DepartamentoID:       req.DepartamentoID,
		Telefono:             req.Telefono,
		Email:                "pendiente@agencia.temp",
		EncargadoPrincipalID: req.EncargadoPrincipalID,
		Status:               "en_revision",
		VisiblePublico:       false,
		AceptaQR:             true,
		AceptaTransferencia:  true,
		AceptaEfectivo:       true,
		LicenciaTuristica:    false,
		CreatedBy:            claims.UserID,
	}

	tx := database.GetDB().Begin()

	if err := tx.Create(&agencia).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear agencia", err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := ensurePaquetePoliticasRow(tx, agencia.ID); err != nil {
		tx.Rollback()
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear politicas de paquetes", err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear agencia", err.Error(), http.StatusInternalServerError)
		return
	}

	// Recargar con relaciones
	database.GetDB().
		Preload("Departamento").
		Preload("EncargadoPrincipal").
		First(&agencia, agencia.ID)

	utils.SuccessResponse(w, agencia, "Agencia creada exitosamente. Complete los datos restantes.", http.StatusCreated)
}

// CreateAgenciaCompleta crea una agencia con toda la información
func (h *AgenciaHandler) CreateAgenciaCompleta(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAgenciaCompletaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validación", err.Error(), http.StatusBadRequest)
		return
	}

	req.NombreComercial = capitalizeAgenciaName(req.NombreComercial)
	req.Email = strings.ToLower(req.Email)

	// Verificar encargado
	if req.EncargadoPrincipalID != nil {
		var encargado models.Usuario
		if err := database.GetDB().First(&encargado, *req.EncargadoPrincipalID).Error; err != nil {
			utils.ErrorResponse(w, "ENCARGADO_NOT_FOUND", "Encargado no encontrado", nil, http.StatusBadRequest)
			return
		}

		if encargado.Rol != "encargado_agencia" {
			utils.ErrorResponse(w, "INVALID_ROL", "El encargado debe tener el rol 'encargado_agencia'", nil, http.StatusBadRequest)
			return
		}
	}

	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	// Solo administradores pueden crear agencias
	if claims.Rol != "admin" {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para crear agencias", nil, http.StatusForbidden)
		return
	}

	// Iniciar transacción
	tx := database.GetDB().Begin()

	agencia := models.AgenciaTurismo{
		NombreComercial:      req.NombreComercial,
		Descripcion:          req.Descripcion,
		Direccion:            req.Direccion,
		DepartamentoID:       req.DepartamentoID,
		Latitud:              req.Latitud,
		Longitud:             req.Longitud,
		Telefono:             req.Telefono,
		Email:                req.Email,
		SitioWeb:             req.SitioWeb,
		Facebook:             req.Facebook,
		Instagram:            req.Instagram,
		LicenciaTuristica:    req.LicenciaTuristica,
		AceptaQR:             req.AceptaQR,
		AceptaTransferencia:  req.AceptaTransferencia,
		AceptaEfectivo:       req.AceptaEfectivo,
		EncargadoPrincipalID: req.EncargadoPrincipalID,
		Status:               req.Status,
		VisiblePublico:       req.VisiblePublico,
		CreatedBy:            claims.UserID,
	}

	// Horarios (opcionales en BD)
	if strings.TrimSpace(req.HorarioApertura) != "" {
		agencia.HorarioApertura = &req.HorarioApertura
	}
	if strings.TrimSpace(req.HorarioCierre) != "" {
		agencia.HorarioCierre = &req.HorarioCierre
	}

	if err := tx.Create(&agencia).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear agencia", err.Error(), http.StatusInternalServerError)
		return
	}

	// Asociar dias (tabla agencia_dias: agencia_id, dia_id)
	if len(req.DiasIDs) > 0 {
		for _, diaID := range req.DiasIDs {
			if err := tx.Exec(
				"INSERT INTO agencia_dias (agencia_id, dia_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				agencia.ID,
				diaID,
			).Error; err != nil {
				tx.Rollback()
				utils.ErrorResponse(w, "DB_ERROR", "Error al asociar dias", err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if _, err := ensurePaquetePoliticasRow(tx, agencia.ID); err != nil {
		tx.Rollback()
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear politicas de paquetes", err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear agencia", err.Error(), http.StatusInternalServerError)
		return
	}

	// Recargar
	database.GetDB().
		Preload("Departamento").
		Preload("EncargadoPrincipal").
		Preload("Dias").
		First(&agencia, agencia.ID)

	utils.SuccessResponse(w, agencia, "Agencia creada exitosamente", http.StatusCreated)
}

// UpdateAgencia actualiza una agencia
func (h *AgenciaHandler) UpdateAgencia(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var req models.UpdateAgenciaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", err.Error(), http.StatusBadRequest)
		return
	}

	// Encargado: bloquear campos restringidos
	if claims.Rol == "encargado_agencia" {
		if req.VisiblePublico != nil || req.Status != nil || req.LicenciaTuristica != nil || req.EncargadoPrincipalID != nil {
			utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para modificar campos restringidos", nil, http.StatusForbidden)
			return
		}
	}

	// Aplicar actualizaciones
	if req.NombreComercial != nil {
		agencia.NombreComercial = capitalizeAgenciaName(*req.NombreComercial)
	}

	if req.Descripcion != nil {
		agencia.Descripcion = *req.Descripcion
	}

	if req.Direccion != nil {
		agencia.Direccion = *req.Direccion
	}

	if req.DepartamentoID != nil {
		agencia.DepartamentoID = *req.DepartamentoID
	}

	if req.Latitud != nil {
		agencia.Latitud = req.Latitud
	}

	if req.Longitud != nil {
		agencia.Longitud = req.Longitud
	}

	if req.Telefono != nil {
		agencia.Telefono = *req.Telefono
	}

	if req.Email != nil {
		agencia.Email = strings.ToLower(*req.Email)
	}

	if req.SitioWeb != nil {
		agencia.SitioWeb = *req.SitioWeb
	}

	if req.Facebook != nil {
		agencia.Facebook = *req.Facebook
	}

	if req.Instagram != nil {
		agencia.Instagram = *req.Instagram
	}

	if req.LicenciaTuristica != nil {
		agencia.LicenciaTuristica = *req.LicenciaTuristica
	}

	if req.HorarioApertura != nil {
		if strings.TrimSpace(*req.HorarioApertura) == "" {
			agencia.HorarioApertura = nil
		} else {
			agencia.HorarioApertura = req.HorarioApertura
		}
	}

	if req.HorarioCierre != nil {
		if strings.TrimSpace(*req.HorarioCierre) == "" {
			agencia.HorarioCierre = nil
		} else {
			agencia.HorarioCierre = req.HorarioCierre
		}
	}

	if req.AceptaQR != nil {
		agencia.AceptaQR = *req.AceptaQR
	}

	if req.AceptaTransferencia != nil {
		agencia.AceptaTransferencia = *req.AceptaTransferencia
	}

	if req.AceptaEfectivo != nil {
		agencia.AceptaEfectivo = *req.AceptaEfectivo
	}

	if req.EncargadoPrincipalID != nil {
		var encargado models.Usuario
		if err := database.GetDB().First(&encargado, *req.EncargadoPrincipalID).Error; err != nil {
			utils.ErrorResponse(w, "ENCARGADO_NOT_FOUND", "Encargado no encontrado", nil, http.StatusBadRequest)
			return
		}

		if encargado.Rol != "encargado_agencia" {
			utils.ErrorResponse(w, "INVALID_ROL", "El encargado debe tener el rol 'encargado_agencia'", nil, http.StatusBadRequest)
			return
		}

		agencia.EncargadoPrincipalID = req.EncargadoPrincipalID
	}

	if req.Status != nil {
		agencia.Status = *req.Status
	}

	if req.VisiblePublico != nil {
		agencia.VisiblePublico = *req.VisiblePublico
	}

	tx := database.GetDB().Begin()

	if err := tx.Save(&agencia).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar agencia", err.Error(), http.StatusInternalServerError)
		return
	}

	// Actualizar dias solo si viene el campo dias_ids en el request
	if req.DiasIDs != nil {
		if err := tx.Exec("DELETE FROM agencia_dias WHERE agencia_id = ?", agencia.ID).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar dias", err.Error(), http.StatusInternalServerError)
			return
		}

		for _, diaID := range req.DiasIDs {
			if err := tx.Exec(
				"INSERT INTO agencia_dias (agencia_id, dia_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				agencia.ID,
				diaID,
			).Error; err != nil {
				tx.Rollback()
				utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar dias", err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar agencia", err.Error(), http.StatusInternalServerError)
		return
	}

	database.GetDB().
		Preload("Departamento").
		Preload("EncargadoPrincipal").
		Preload("Fotos").
		Preload("Especialidades.Categoria").
		Preload("Dias").
		First(&agencia, agencia.ID)

	utils.SuccessResponse(w, agencia, "Agencia actualizada exitosamente", http.StatusOK)
}

// DeleteAgencia soft delete de agencia
func (h *AgenciaHandler) DeleteAgencia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	agencia.Status = "inactiva"
	agencia.VisiblePublico = false

	if err := database.GetDB().Save(&agencia).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al eliminar agencia", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Agencia eliminada exitosamente", http.StatusOK)
}

// UpdateAgenciaStatus actualiza solo el status
func (h *AgenciaHandler) UpdateAgenciaStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	var req models.UpdateAgenciaStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validación", err.Error(), http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	agencia.Status = req.Status
	if err := database.GetDB().Save(&agencia).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar status", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, agencia, "Status actualizado exitosamente", http.StatusOK)
}

// AddEspecialidad agrega una especialidad a la agencia
func (h *AgenciaHandler) AddEspecialidad(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invǭlido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var req struct {
		CategoriaID uint `json:"categoria_id" validate:"required"`
		EsPrincipal bool `json:"es_principal"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validación", err.Error(), http.StatusBadRequest)
		return
	}

	// Verificar que no exista
	var existing models.AgenciaEspecialidad
	if err := database.GetDB().Where("agencia_id = ? AND categoria_id = ?", agenciaID, req.CategoriaID).First(&existing).Error; err == nil {
		utils.ErrorResponse(w, "DUPLICATE", "Esta especialidad ya está asignada", nil, http.StatusBadRequest)
		return
	}

	// Si es principal, quitar flag de otras
	if req.EsPrincipal {
		database.GetDB().Model(&models.AgenciaEspecialidad{}).
			Where("agencia_id = ?", agenciaID).
			Update("es_principal", false)
	}

	especialidad := models.AgenciaEspecialidad{
		AgenciaID:   uint(agenciaID),
		CategoriaID: req.CategoriaID,
		EsPrincipal: req.EsPrincipal,
	}

	if err := database.GetDB().Create(&especialidad).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al agregar especialidad", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, especialidad, "Especialidad agregada exitosamente", http.StatusCreated)
}

// RemoveEspecialidad elimina una especialidad
func (h *AgenciaHandler) RemoveEspecialidad(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}
	especialidadID, err := strconv.ParseUint(vars["especialidad_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	if err := database.GetDB().Where("id = ? AND agencia_id = ?", especialidadID, agenciaID).Delete(&models.AgenciaEspecialidad{}).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al eliminar especialidad", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Especialidad eliminada exitosamente", http.StatusOK)
}

// GetStats obtiene estadísticas de agencias
func (h *AgenciaHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		Total       int64 `json:"total"`
		Activas     int64 `json:"activas"`
		Inactivas   int64 `json:"inactivas"`
		Suspendidas int64 `json:"suspendidas"`
		EnRevision  int64 `json:"en_revision"`
		ConLicencia int64 `json:"con_licencia"`
	}

	database.GetDB().Model(&models.AgenciaTurismo{}).Count(&stats.Total)
	database.GetDB().Model(&models.AgenciaTurismo{}).Where("status = ?", "activa").Count(&stats.Activas)
	database.GetDB().Model(&models.AgenciaTurismo{}).Where("status = ?", "inactiva").Count(&stats.Inactivas)
	database.GetDB().Model(&models.AgenciaTurismo{}).Where("status = ?", "suspendida").Count(&stats.Suspendidas)
	database.GetDB().Model(&models.AgenciaTurismo{}).Where("status = ?", "en_revision").Count(&stats.EnRevision)
	database.GetDB().Model(&models.AgenciaTurismo{}).Where("licencia_turistica = ?", true).Count(&stats.ConLicencia)

	utils.SuccessResponse(w, stats, "Estadísticas obtenidas exitosamente", http.StatusOK)
}

// GetDepartamentos lista departamentos
func (h *AgenciaHandler) GetDepartamentos(w http.ResponseWriter, r *http.Request) {
	var departamentos []models.Departamento
	if err := database.GetDB().Find(&departamentos).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener departamentos", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, departamentos, "", http.StatusOK)
}

// GetCategorias lista categorías
func (h *AgenciaHandler) GetCategorias(w http.ResponseWriter, r *http.Request) {
	var categorias []models.CategoriaAtraccion
	if err := database.GetDB().Find(&categorias).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener categorías", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, categorias, "", http.StatusOK)
}

// GetDias lista días
func (h *AgenciaHandler) GetDias(w http.ResponseWriter, r *http.Request) {
	var dias []models.Dia
	if err := database.GetDB().Order("id_dia").Find(&dias).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener días", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, dias, "", http.StatusOK)
}

// GetEncargados lista encargados
func (h *AgenciaHandler) GetEncargados(w http.ResponseWriter, r *http.Request) {
	var encargados []models.Usuario

	onlyUnassigned := strings.ToLower(r.URL.Query().Get("only_unassigned"))
	agenciaIDStr := r.URL.Query().Get("agencia_id")

	db := database.GetDB().Model(&models.Usuario{}).
		Where("rol = ? AND status = ?", "encargado_agencia", "active")

	if onlyUnassigned == "true" || onlyUnassigned == "1" {
		subquery := database.GetDB().
			Model(&models.AgenciaTurismo{}).
			Select("encargado_principal_id").
			Where("encargado_principal_id IS NOT NULL")

		// Si se envía agencia_id, no excluir el encargado asignado a esa agencia
		if agenciaIDStr != "" {
			if agenciaID, err := strconv.ParseUint(agenciaIDStr, 10, 64); err == nil && agenciaID > 0 {
				subquery = subquery.Where("id <> ?", agenciaID)
			}
		}

		db = db.Where("id NOT IN (?)", subquery)
	}

	if err := db.Order("nombre").Order("apellido_paterno").Order("apellido_materno").Find(&encargados).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener encargados", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, encargados, "", http.StatusOK)
}

// GetMiAgencia obtiene la agencia asignada al encargado autenticado
func (h *AgenciaHandler) GetMiAgencia(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	if claims.Rol != "encargado_agencia" {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para ver esta agencia", nil, http.StatusForbidden)
		return
	}

	var agencia models.AgenciaTurismo
	err := database.GetDB().
		Preload("Departamento").
		Preload("EncargadoPrincipal").
		Preload("Fotos").
		Preload("Especialidades.Categoria").
		Preload("Dias").
		Where("encargado_principal_id = ?", claims.UserID).
		First(&agencia).Error

	if err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "No tiene una agencia asignada", nil, http.StatusNotFound)
		return
	}

	utils.SuccessResponse(w, agencia, "Agencia obtenida exitosamente", http.StatusOK)
}
