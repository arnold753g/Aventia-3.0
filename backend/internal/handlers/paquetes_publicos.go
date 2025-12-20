package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetPaquetesPublicos lista paquetes turísticos visibles para el público.
// Nota: estos endpoints se usan para la vista de Turista (no requieren ser admin/encargado).
func (h *AgenciaHandler) GetPaquetesPublicos(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 12
	}
	if limit > 100 {
		limit = 100
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	frecuencia := strings.TrimSpace(r.URL.Query().Get("frecuencia"))
	nivelDificultad := strings.TrimSpace(r.URL.Query().Get("nivel_dificultad"))
	tipoDuracion := strings.TrimSpace(r.URL.Query().Get("tipo_duracion")) // un_dia | varios_dias

	precioMinStr := strings.TrimSpace(r.URL.Query().Get("precio_min"))
	precioMaxStr := strings.TrimSpace(r.URL.Query().Get("precio_max"))

	sortBy := strings.TrimSpace(r.URL.Query().Get("sort_by"))
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("sort_order")))
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	allowedSort := map[string]string{
		"created_at": "paquetes_turisticos.created_at",
		"precio":     "paquetes_turisticos.precio_base_nacionales",
		"nombre":     "paquetes_turisticos.nombre",
	}
	orderCol, ok := allowedSort[sortBy]
	if !ok {
		orderCol = allowedSort["created_at"]
	}

	db := database.GetDB()
	query := db.Model(&models.PaqueteTuristico{}).
		Joins("JOIN agencias_turismo a ON a.id = paquetes_turisticos.agencia_id").
		Where("paquetes_turisticos.status = ?", "activo").
		Where("paquetes_turisticos.visible_publico = TRUE").
		Where("a.status = ?", "activa").
		Where("a.visible_publico = TRUE")

	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("paquetes_turisticos.nombre ILIKE ? OR paquetes_turisticos.descripcion ILIKE ?", pattern, pattern)
	}

	if frecuencia != "" {
		if !allowedPaqueteFrecuencias[frecuencia] {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "Frecuencia inválida", nil, http.StatusBadRequest)
			return
		}
		query = query.Where("paquetes_turisticos.frecuencia = ?", frecuencia)
	}

	if nivelDificultad != "" {
		query = query.Where("paquetes_turisticos.nivel_dificultad = ?", nivelDificultad)
	}

	if tipoDuracion != "" {
		switch tipoDuracion {
		case "un_dia":
			query = query.Where("paquetes_turisticos.duracion_dias = 1")
		case "varios_dias":
			query = query.Where("paquetes_turisticos.duracion_dias > 1")
		default:
			utils.ErrorResponse(w, "VALIDATION_ERROR", "tipo_duracion inválido (use un_dia o varios_dias)", nil, http.StatusBadRequest)
			return
		}
	}

	if precioMinStr != "" {
		min, err := strconv.ParseFloat(precioMinStr, 64)
		if err != nil || min < 0 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "precio_min inválido", nil, http.StatusBadRequest)
			return
		}
		query = query.Where("paquetes_turisticos.precio_base_nacionales >= ?", min)
	}

	if precioMaxStr != "" {
		max, err := strconv.ParseFloat(precioMaxStr, 64)
		if err != nil || max < 0 {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "precio_max inválido", nil, http.StatusBadRequest)
			return
		}
		query = query.Where("paquetes_turisticos.precio_base_nacionales <= ?", max)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * limit
	var paquetes []models.PaqueteTuristico
	if err := query.
		Preload("Fotos", func(db *gorm.DB) *gorm.DB {
			return db.Order("es_principal DESC").Order("orden ASC").Order("id ASC")
		}).
		Preload("Agencia.Departamento").
		Order(orderCol + " " + sortOrder).
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

// GetPaquetePublico obtiene el detalle de un paquete visible para el público.
func (h *AgenciaHandler) GetPaquetePublico(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	db := database.GetDB()

	var paquete models.PaqueteTuristico
	if err := db.Model(&models.PaqueteTuristico{}).
		Joins("JOIN agencias_turismo a ON a.id = paquetes_turisticos.agencia_id").
		Where("paquetes_turisticos.id = ?", uint(id64)).
		Where("paquetes_turisticos.status = ?", "activo").
		Where("paquetes_turisticos.visible_publico = TRUE").
		Where("a.status = ?", "activa").
		Where("a.visible_publico = TRUE").
		Preload("Fotos", func(db *gorm.DB) *gorm.DB {
			return db.Order("es_principal DESC").Order("orden ASC").Order("id ASC")
		}).
		Preload("Itinerario", func(db *gorm.DB) *gorm.DB {
			return db.Order("dia_numero ASC").Order("id ASC")
		}).
		Preload("Atracciones", func(db *gorm.DB) *gorm.DB {
			return db.Order("dia_numero ASC NULLS FIRST, orden_visita ASC, id ASC")
		}).
		Preload("Atracciones.Atraccion.Provincia.Departamento").
		Preload("Atracciones.Atraccion.Fotos", func(db *gorm.DB) *gorm.DB {
			return db.Order("es_principal DESC").Order("orden ASC").Order("id ASC")
		}).
		Preload("Agencia.Departamento").
		Preload("Agencia.Fotos", func(db *gorm.DB) *gorm.DB {
			return db.Order("es_principal DESC").Order("orden ASC").Order("id ASC")
		}).
		First(&paquete).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
			return
		}
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener paquete", err.Error(), http.StatusInternalServerError)
		return
	}

	// Políticas (por agencia)
	var politicas models.PaquetePolitica
	if err := db.Where("agencia_id = ?", paquete.AgenciaID).First(&politicas).Error; err == nil {
		paquete.Politicas = &politicas
	}

	// Datos de pago (por agencia, opcional)
	var datos models.AgenciaDatosPago
	if err := db.Where("agencia_id = ? AND activo = TRUE", paquete.AgenciaID).First(&datos).Error; err == nil {
		paquete.AgenciaDatosPago = &datos
	}

	utils.SuccessResponse(w, paquete, "Paquete obtenido exitosamente", http.StatusOK)
}
