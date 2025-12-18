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

var allowedPaqueteFrecuencias = map[string]bool{
	"salida_diaria": true,
	"salida_unica":  true,
}

var allowedPaqueteHorarios = map[string]bool{
	"mañana":   true,
	"tarde":    true,
	"todo_dia": true,
}

var allowedPaqueteDificultades = map[string]bool{
	"facil":   true,
	"medio":   true,
	"dificil": true,
	"extremo": true,
}

var allowedPaqueteStatus = map[string]bool{
	"activo":    true,
	"inactivo":  true,
	"borrador":  true,
	"eliminado": true,
}

type createPaqueteRequest struct {
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion"`

	Frecuencia string `json:"frecuencia"`

	DuracionDias   *int `json:"duracion_dias"`
	DuracionNoches *int `json:"duracion_noches"`

	FechaSalidaFija *string `json:"fecha_salida_fija"`

	Horario       *string `json:"horario"`
	HoraSalida    *string `json:"hora_salida"`
	DuracionHoras *string `json:"duracion_horas"`

	DiasPreviosCompra *int    `json:"dias_previos_compra"`
	NivelDificultad   *string `json:"nivel_dificultad"`

	CupoMinimo     *int  `json:"cupo_minimo"`
	CupoMaximo     *int  `json:"cupo_maximo"`
	PermitePrivado *bool `json:"permite_privado"`

	PrecioBaseNacionales       *float64 `json:"precio_base_nacionales"`
	PrecioAdicionalExtranjeros *float64 `json:"precio_adicional_extranjeros"`

	Incluye   []string `json:"incluye"`
	NoIncluye []string `json:"no_incluye"`
	QueLlevar []string `json:"que_llevar"`

	Status         *string `json:"status"`
	VisiblePublico *bool   `json:"visible_publico"`
}

type updatePaqueteRequest struct {
	Nombre      *string `json:"nombre"`
	Descripcion *string `json:"descripcion"`

	// No se permite cambiar frecuencia; debe definirse al crear.
	Frecuencia *string `json:"frecuencia"`

	DuracionDias   *int `json:"duracion_dias"`
	DuracionNoches *int `json:"duracion_noches"`

	FechaSalidaFija *string `json:"fecha_salida_fija"`

	Horario       *string `json:"horario"`
	HoraSalida    *string `json:"hora_salida"`
	DuracionHoras *string `json:"duracion_horas"`

	DiasPreviosCompra *int    `json:"dias_previos_compra"`
	NivelDificultad   *string `json:"nivel_dificultad"`

	CupoMinimo     *int  `json:"cupo_minimo"`
	CupoMaximo     *int  `json:"cupo_maximo"`
	PermitePrivado *bool `json:"permite_privado"`

	PrecioBaseNacionales       *float64 `json:"precio_base_nacionales"`
	PrecioAdicionalExtranjeros *float64 `json:"precio_adicional_extranjeros"`

	Incluye   *[]string `json:"incluye"`
	NoIncluye *[]string `json:"no_incluye"`
	QueLlevar *[]string `json:"que_llevar"`

	Status         *string `json:"status"`
	VisiblePublico *bool   `json:"visible_publico"`
}

func normalizeDatePtr(raw *string) (*string, error) {
	if raw == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil, nil
	}
	if _, err := time.Parse("2006-01-02", trimmed); err != nil {
		return nil, err
	}
	return &trimmed, nil
}

func normalizeStringPtr(raw *string) *string {
	if raw == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func paqueteDuracionDiasValue(p *models.PaqueteTuristico) int {
	if p == nil || p.DuracionDias == nil || *p.DuracionDias <= 0 {
		return 1
	}
	return *p.DuracionDias
}

func (h *AgenciaHandler) GetAgenciaPaquetes(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	visible := strings.TrimSpace(r.URL.Query().Get("visible_publico"))
	includeEliminado := strings.TrimSpace(r.URL.Query().Get("include_eliminado")) == "true"

	sortBy := strings.TrimSpace(r.URL.Query().Get("sort_by"))
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("sort_order")))
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	db := database.GetDB()
	query := db.Model(&models.PaqueteTuristico{}).Where("agencia_id = ?", agenciaID)

	if !includeEliminado {
		query = query.Where("status <> 'eliminado'")
	}

	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("nombre ILIKE ? OR descripcion ILIKE ?", pattern, pattern)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if visible != "" {
		query = query.Where("visible_publico = ?", visible == "true")
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * limit
	var paquetes []models.PaqueteTuristico
	if err := query.
		Preload("Fotos").
		Order(sortBy + " " + sortOrder).
		Limit(limit).
		Offset(offset).
		Find(&paquetes).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener paquetes", err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"paquetes": paquetes,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utils.SuccessResponse(w, response, "Paquetes obtenidos exitosamente", http.StatusOK)
}

func (h *AgenciaHandler) CreateAgenciaPaquete(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var req createPaqueteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(req.Nombre)
	if nombre == "" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El nombre es obligatorio", nil, http.StatusBadRequest)
		return
	}

	frecuencia := strings.TrimSpace(req.Frecuencia)
	if !allowedPaqueteFrecuencias[frecuencia] {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Frecuencia inválida", nil, http.StatusBadRequest)
		return
	}

	if req.CupoMinimo == nil || req.CupoMaximo == nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Debe especificar cupo_minimo y cupo_maximo", nil, http.StatusBadRequest)
		return
	}
	if *req.CupoMinimo <= 0 || *req.CupoMaximo <= 0 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Los cupos deben ser mayores a 0", nil, http.StatusBadRequest)
		return
	}
	if *req.CupoMaximo < *req.CupoMinimo {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "cupo_maximo debe ser >= cupo_minimo", nil, http.StatusBadRequest)
		return
	}

	if req.PrecioBaseNacionales == nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Debe especificar precio_base_nacionales", nil, http.StatusBadRequest)
		return
	}
	if *req.PrecioBaseNacionales < 0 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El precio base no puede ser negativo", nil, http.StatusBadRequest)
		return
	}

	precioExtra := 0.0
	if req.PrecioAdicionalExtranjeros != nil {
		if *req.PrecioAdicionalExtranjeros < 0 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "El precio adicional no puede ser negativo", nil, http.StatusBadRequest)
			return
		}
		precioExtra = *req.PrecioAdicionalExtranjeros
	}

	diasPrevios := 1
	if req.DiasPreviosCompra != nil {
		diasPrevios = *req.DiasPreviosCompra
	}
	if diasPrevios < 1 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "dias_previos_compra debe ser >= 1", nil, http.StatusBadRequest)
		return
	}

	var duracionDias *int
	var duracionNoches *int
	isMultiDay := false
	if req.DuracionDias != nil && *req.DuracionDias > 1 {
		isMultiDay = true
	}

	if isMultiDay {
		if req.DuracionDias == nil || *req.DuracionDias <= 1 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "duracion_dias debe ser > 1 para paquetes de varios días", nil, http.StatusBadRequest)
			return
		}
		if req.DuracionNoches == nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "duracion_noches es obligatoria para paquetes de varios días", nil, http.StatusBadRequest)
			return
		}
		expectedNoches := *req.DuracionDias - 1
		if *req.DuracionNoches != expectedNoches {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "duracion_noches debe ser duracion_dias - 1", nil, http.StatusBadRequest)
			return
		}
		d := *req.DuracionDias
		n := *req.DuracionNoches
		duracionDias = &d
		duracionNoches = &n
	} else {
		one := 1
		duracionDias = &one
		duracionNoches = nil
	}

	// Validaciones por frecuencia
	fechaSalidaFija, err := normalizeDatePtr(req.FechaSalidaFija)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "fecha_salida_fija inválida (use YYYY-MM-DD)", nil, http.StatusBadRequest)
		return
	}
	if frecuencia == "salida_unica" && fechaSalidaFija == nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "fecha_salida_fija es obligatoria para salida_unica", nil, http.StatusBadRequest)
		return
	}
	if frecuencia == "salida_diaria" {
		fechaSalidaFija = nil
	}

	// Validaciones por duración
	horario := normalizeStringPtr(req.Horario)
	if !isMultiDay && horario != nil && !allowedPaqueteHorarios[*horario] {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Horario inválido", nil, http.StatusBadRequest)
		return
	}
	if isMultiDay {
		horario = nil
	}

	horaSalida := normalizeStringPtr(req.HoraSalida)
	duracionHoras := normalizeStringPtr(req.DuracionHoras)
	if isMultiDay {
		horaSalida = nil
		duracionHoras = nil
	}

	nivelDificultad := normalizeStringPtr(req.NivelDificultad)
	if nivelDificultad != nil && !allowedPaqueteDificultades[*nivelDificultad] {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Nivel de dificultad inválido", nil, http.StatusBadRequest)
		return
	}

	status := "borrador"
	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		s := strings.TrimSpace(*req.Status)
		if !allowedPaqueteStatus[s] || s == "eliminado" {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "Status inválido", nil, http.StatusBadRequest)
			return
		}
		status = s
	}

	visiblePublico := true
	if req.VisiblePublico != nil {
		visiblePublico = *req.VisiblePublico
	}

	permitePrivado := true
	if req.PermitePrivado != nil {
		permitePrivado = *req.PermitePrivado
	}

	p := models.PaqueteTuristico{
		AgenciaID:   agenciaID,
		Nombre:     nombre,
		Descripcion: normalizeStringPtr(req.Descripcion),
		Frecuencia: frecuencia,

		DuracionDias:   duracionDias,
		DuracionNoches: duracionNoches,

		FechaSalidaFija: fechaSalidaFija,

		Horario:       horario,
		HoraSalida:    horaSalida,
		DuracionHoras: duracionHoras,

		DiasPreviosCompra: diasPrevios,
		NivelDificultad:   nivelDificultad,

		CupoMinimo:     *req.CupoMinimo,
		CupoMaximo:     *req.CupoMaximo,
		PermitePrivado: permitePrivado,

		PrecioBaseNacionales:       *req.PrecioBaseNacionales,
		PrecioAdicionalExtranjeros: precioExtra,

		Incluye:   models.StringArray(req.Incluye),
		NoIncluye: models.StringArray(req.NoIncluye),
		QueLlevar: models.StringArray(req.QueLlevar),

		Status:         status,
		VisiblePublico: visiblePublico,
	}

	db := database.GetDB()
	if err := db.Create(&p).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear paquete", err.Error(), http.StatusInternalServerError)
		return
	}

	// Recargar con salidas/fotos (si aplica trigger)
	db.Preload("Fotos").Preload("Salidas").First(&p, p.ID)

	utils.SuccessResponse(w, p, "Paquete creado exitosamente", http.StatusCreated)
}

func (h *AgenciaHandler) GetAgenciaPaquete(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de paquete inválido", nil, http.StatusBadRequest)
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
	if err := database.GetDB().
		Preload("Fotos").
		Preload("Itinerario").
		Preload("Atracciones.Atraccion").
		Preload("Salidas").
		Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).
		First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	utils.SuccessResponse(w, paquete, "Paquete obtenido exitosamente", http.StatusOK)
}

func (h *AgenciaHandler) UpdateAgenciaPaquete(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de paquete inválido", nil, http.StatusBadRequest)
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

	db := database.GetDB()

	var paquete models.PaqueteTuristico
	if err := db.Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	var req updatePaqueteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	if req.Frecuencia != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "No se permite cambiar la frecuencia del paquete", nil, http.StatusBadRequest)
		return
	}

	if req.Nombre != nil {
		n := strings.TrimSpace(*req.Nombre)
		if n == "" {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "El nombre no puede estar vacío", nil, http.StatusBadRequest)
			return
		}
		paquete.Nombre = n
	}

	if req.Descripcion != nil {
		paquete.Descripcion = normalizeStringPtr(req.Descripcion)
	}

	if req.DiasPreviosCompra != nil {
		if *req.DiasPreviosCompra < 1 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "dias_previos_compra debe ser >= 1", nil, http.StatusBadRequest)
			return
		}
		paquete.DiasPreviosCompra = *req.DiasPreviosCompra
	}

	if req.NivelDificultad != nil {
		nd := normalizeStringPtr(req.NivelDificultad)
		if nd != nil && !allowedPaqueteDificultades[*nd] {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "Nivel de dificultad inválido", nil, http.StatusBadRequest)
			return
		}
		paquete.NivelDificultad = nd
	}

	if req.CupoMinimo != nil {
		if *req.CupoMinimo <= 0 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "cupo_minimo debe ser > 0", nil, http.StatusBadRequest)
			return
		}
		paquete.CupoMinimo = *req.CupoMinimo
	}
	if req.CupoMaximo != nil {
		if *req.CupoMaximo <= 0 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "cupo_maximo debe ser > 0", nil, http.StatusBadRequest)
			return
		}
		paquete.CupoMaximo = *req.CupoMaximo
	}
	if paquete.CupoMaximo < paquete.CupoMinimo {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "cupo_maximo debe ser >= cupo_minimo", nil, http.StatusBadRequest)
		return
	}

	if req.PermitePrivado != nil {
		paquete.PermitePrivado = *req.PermitePrivado
	}

	if req.PrecioBaseNacionales != nil {
		if *req.PrecioBaseNacionales < 0 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "El precio base no puede ser negativo", nil, http.StatusBadRequest)
			return
		}
		paquete.PrecioBaseNacionales = *req.PrecioBaseNacionales
	}
	if req.PrecioAdicionalExtranjeros != nil {
		if *req.PrecioAdicionalExtranjeros < 0 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "El precio adicional no puede ser negativo", nil, http.StatusBadRequest)
			return
		}
		paquete.PrecioAdicionalExtranjeros = *req.PrecioAdicionalExtranjeros
	}

	if req.Incluye != nil {
		paquete.Incluye = models.StringArray(*req.Incluye)
	}
	if req.NoIncluye != nil {
		paquete.NoIncluye = models.StringArray(*req.NoIncluye)
	}
	if req.QueLlevar != nil {
		paquete.QueLlevar = models.StringArray(*req.QueLlevar)
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		s := strings.TrimSpace(*req.Status)
		if !allowedPaqueteStatus[s] {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "Status inválido", nil, http.StatusBadRequest)
			return
		}
		paquete.Status = s
	}

	if req.VisiblePublico != nil {
		paquete.VisiblePublico = *req.VisiblePublico
	}

	// Duración (requiere validar consistencia)
	newDuracionDias := paqueteDuracionDiasValue(&paquete)
	if req.DuracionDias != nil {
		if *req.DuracionDias <= 0 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "duracion_dias debe ser >= 1", nil, http.StatusBadRequest)
			return
		}
		newDuracionDias = *req.DuracionDias
	}

	isMultiDay := newDuracionDias > 1

	if isMultiDay {
		paquete.DuracionDias = &newDuracionDias
		if req.DuracionNoches != nil {
			paquete.DuracionNoches = req.DuracionNoches
		}
		if paquete.DuracionNoches == nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "duracion_noches es obligatoria para paquetes de varios días", nil, http.StatusBadRequest)
			return
		}
		expectedNoches := newDuracionDias - 1
		if *paquete.DuracionNoches != expectedNoches {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "duracion_noches debe ser duracion_dias - 1", nil, http.StatusBadRequest)
			return
		}

		// Limpiar campos solo de 1 día
		paquete.Horario = nil
		paquete.HoraSalida = nil
		paquete.DuracionHoras = nil
	} else {
		one := 1
		paquete.DuracionDias = &one
		paquete.DuracionNoches = nil

		// Permitir actualizar campos de 1 día
		if req.Horario != nil {
			hor := normalizeStringPtr(req.Horario)
			if hor != nil && !allowedPaqueteHorarios[*hor] {
				utils.ErrorResponse(w, "VALIDATION_ERROR", "Horario inválido", nil, http.StatusBadRequest)
				return
			}
			paquete.Horario = hor
		}
		if req.HoraSalida != nil {
			paquete.HoraSalida = normalizeStringPtr(req.HoraSalida)
		}
		if req.DuracionHoras != nil {
			paquete.DuracionHoras = normalizeStringPtr(req.DuracionHoras)
		}
	}

	// Fecha salida fija: solo para salida_unica
	if req.FechaSalidaFija != nil {
		if paquete.Frecuencia != "salida_unica" {
			paquete.FechaSalidaFija = nil
		} else {
			newFecha, err := normalizeDatePtr(req.FechaSalidaFija)
			if err != nil || newFecha == nil {
				utils.ErrorResponse(w, "VALIDATION_ERROR", "fecha_salida_fija inválida (use YYYY-MM-DD)", nil, http.StatusBadRequest)
				return
			}

			// Si existe una salida compartida creada automáticamente, intentar mantener consistencia.
			tx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			var privateCount int64
			tx.Model(&models.PaqueteSalidaHabilitada{}).
				Where("paquete_id = ? AND tipo_salida = 'privado'", paquete.ID).
				Count(&privateCount)
			if privateCount > 0 {
				tx.Rollback()
				utils.ErrorResponse(w, "VALIDATION_ERROR", "No se puede cambiar la fecha: existen salidas privadas", nil, http.StatusBadRequest)
				return
			}

			var salida models.PaqueteSalidaHabilitada
			if err := tx.
				Where("paquete_id = ? AND tipo_salida = 'compartido'", paquete.ID).
				Order("id asc").
				First(&salida).Error; err == nil {
				if salida.CuposReservados > 0 || salida.CuposConfirmados > 0 {
					tx.Rollback()
					utils.ErrorResponse(w, "VALIDATION_ERROR", "No se puede cambiar la fecha: ya existen cupos reservados/confirmados", nil, http.StatusBadRequest)
					return
				}
				if err := tx.Model(&models.PaqueteSalidaHabilitada{}).
					Where("id = ?", salida.ID).
					Update("fecha_salida", *newFecha).Error; err != nil {
					tx.Rollback()
					utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar la salida compartida", err.Error(), http.StatusInternalServerError)
					return
				}
			}

			paquete.FechaSalidaFija = newFecha

			if err := tx.Save(&paquete).Error; err != nil {
				tx.Rollback()
				utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar paquete", err.Error(), http.StatusInternalServerError)
				return
			}
			if err := tx.Commit().Error; err != nil {
				utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar paquete", err.Error(), http.StatusInternalServerError)
				return
			}

			// Recargar y responder
			db.Preload("Fotos").Preload("Salidas").First(&paquete, paquete.ID)
			utils.SuccessResponse(w, paquete, "Paquete actualizado exitosamente", http.StatusOK)
			return
		}
	}

	// Guardar cambios
	if err := db.Save(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar paquete", err.Error(), http.StatusInternalServerError)
		return
	}

	db.Preload("Fotos").
		Preload("Itinerario").
		Preload("Atracciones.Atraccion").
		Preload("Salidas").
		First(&paquete, paquete.ID)

	utils.SuccessResponse(w, paquete, "Paquete actualizado exitosamente", http.StatusOK)
}

func (h *AgenciaHandler) DeleteAgenciaPaquete(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de paquete inválido", nil, http.StatusBadRequest)
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

	db := database.GetDB()
	var paquete models.PaqueteTuristico
	if err := db.Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	paquete.Status = "eliminado"
	paquete.VisiblePublico = false

	if err := db.Save(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al eliminar paquete", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Paquete eliminado exitosamente", http.StatusOK)
}
