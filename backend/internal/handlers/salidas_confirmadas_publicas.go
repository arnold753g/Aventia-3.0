package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/pkg/utils"
)

type salidaConfirmadaPublicaRow struct {
	SalidaID uint `json:"salida_id" gorm:"column:salida_id"`

	PaqueteID           uint    `json:"paquete_id" gorm:"column:paquete_id"`
	PaqueteNombre       string  `json:"paquete_nombre" gorm:"column:paquete_nombre"`
	PaqueteFrecuencia   string  `json:"paquete_frecuencia" gorm:"column:paquete_frecuencia"`
	PaqueteDuracionDias *int    `json:"paquete_duracion_dias,omitempty" gorm:"column:paquete_duracion_dias"`
	PaqueteHorario      *string `json:"paquete_horario,omitempty" gorm:"column:paquete_horario"`
	PaqueteHoraSalida   *string `json:"paquete_hora_salida,omitempty" gorm:"column:paquete_hora_salida"`
	PaqueteFoto         *string `json:"paquete_foto,omitempty" gorm:"column:paquete_foto"`

	AgenciaID     uint   `json:"agencia_id" gorm:"column:agencia_id"`
	AgenciaNombre string `json:"agencia_nombre" gorm:"column:agencia_nombre"`

	FechaSalida      string `json:"fecha_salida" gorm:"column:fecha_salida"`
	TipoSalida       string `json:"tipo_salida" gorm:"column:tipo_salida"`
	Estado           string `json:"estado" gorm:"column:estado"`
	CupoMinimo       int    `json:"cupo_minimo" gorm:"column:cupo_minimo"`
	CupoMaximo       int    `json:"cupo_maximo" gorm:"column:cupo_maximo"`
	CuposReservados  int    `json:"cupos_reservados" gorm:"column:cupos_reservados"`
	CuposConfirmados int    `json:"cupos_confirmados" gorm:"column:cupos_confirmados"`
	CuposDisponibles int    `json:"cupos_disponibles" gorm:"column:cupos_disponibles"`
}

// GetSalidasConfirmadasPublicas lista salidas habilitadas (no canceladas) que ya tienen cupos confirmados.
// Nota: se filtra por paquetes/agencias visibles al público.
func (h *AgenciaHandler) GetSalidasConfirmadasPublicas(w http.ResponseWriter, r *http.Request) {
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
	tipo := strings.TrimSpace(r.URL.Query().Get("tipo"))   // compartido | privado
	desde := strings.TrimSpace(r.URL.Query().Get("desde")) // YYYY-MM-DD
	hasta := strings.TrimSpace(r.URL.Query().Get("hasta")) // YYYY-MM-DD
	minDesdeStr := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	if tipo != "" && tipo != "compartido" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "tipo inválido (solo se permite compartido)", nil, http.StatusBadRequest)
		return
	}

	if desde != "" {
		if _, err := time.Parse("2006-01-02", desde); err != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "desde inválido (use YYYY-MM-DD)", nil, http.StatusBadRequest)
			return
		}
	}
	if hasta != "" {
		if _, err := time.Parse("2006-01-02", hasta); err != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "hasta inválido (use YYYY-MM-DD)", nil, http.StatusBadRequest)
			return
		}
	}

	// No mostrar salidas del día: mínimo desde mañana.
	if desde == "" || desde < minDesdeStr {
		desde = minDesdeStr
	}

	db := database.GetDB()
	base := db.Table("paquete_salidas_habilitadas s").
		Joins("JOIN paquetes_turisticos p ON p.id = s.paquete_id").
		Joins("JOIN agencias_turismo a ON a.id = p.agencia_id").
		Joins(`LEFT JOIN LATERAL (
			SELECT foto
			FROM paquete_fotos pf
			WHERE pf.paquete_id = p.id
			ORDER BY pf.es_principal DESC, pf.orden ASC, pf.id ASC
			LIMIT 1
		) pfoto ON TRUE`).
		Where("p.status = ?", "activo").
		Where("p.visible_publico = TRUE").
		Where("a.status = ?", "activa").
		Where("a.visible_publico = TRUE").
		Where("s.cupos_confirmados > 0").
		Where("s.estado <> 'cancelada'").
		Where("s.tipo_salida = 'compartido'").
		Where("s.fecha_salida >= ?", desde)

	if hasta != "" {
		base = base.Where("s.fecha_salida <= ?", hasta)
	}

	if search != "" {
		pattern := "%" + search + "%"
		base = base.Where("p.nombre ILIKE ?", pattern)
	}

	_ = tipo // Compatibilidad: solo se permite compartido (ya filtrado arriba).

	var total int64
	if err := base.Count(&total).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al contar salidas", err.Error(), http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * limit
	var rows []salidaConfirmadaPublicaRow
	if err := base.
		Select(`
			s.id AS salida_id,
			s.paquete_id,
			p.nombre AS paquete_nombre,
			p.frecuencia AS paquete_frecuencia,
			p.duracion_dias AS paquete_duracion_dias,
			p.horario AS paquete_horario,
			p.hora_salida AS paquete_hora_salida,
			pfoto.foto AS paquete_foto,
			a.id AS agencia_id,
			a.nombre_comercial AS agencia_nombre,
			s.fecha_salida,
			s.tipo_salida,
			s.estado,
			s.cupo_minimo,
			s.cupo_maximo,
			s.cupos_reservados,
			s.cupos_confirmados,
			(s.cupo_maximo - s.cupos_reservados - s.cupos_confirmados) AS cupos_disponibles
		`).
		Order("s.fecha_salida ASC").
		Order("p.nombre ASC").
		Order("s.tipo_salida ASC").
		Order("s.id ASC").
		Limit(limit).
		Offset(offset).
		Scan(&rows).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener salidas", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"salidas": rows,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}, "Salidas obtenidas exitosamente", http.StatusOK)
}
