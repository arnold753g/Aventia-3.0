package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type agenciaVentaPagoRow struct {
	PagoID          uint    `json:"pago_id" gorm:"column:pago_id"`
	CompraID        uint    `json:"compra_id" gorm:"column:compra_id"`
	MetodoPago      string  `json:"metodo_pago" gorm:"column:metodo_pago"`
	Monto           float64 `json:"monto" gorm:"column:monto"`
	Estado          string  `json:"estado" gorm:"column:estado"`
	ComprobanteFoto *string `json:"comprobante_foto,omitempty" gorm:"column:comprobante_foto"`

	ConfirmadoPor     *uint      `json:"confirmado_por,omitempty" gorm:"column:confirmado_por"`
	FechaConfirmacion *time.Time `json:"fecha_confirmacion,omitempty" gorm:"column:fecha_confirmacion"`
	RazonRechazo      *string    `json:"razon_rechazo,omitempty" gorm:"column:razon_rechazo"`
	NotasEncargado    *string    `json:"notas_encargado,omitempty" gorm:"column:notas_encargado"`

	FechaPago time.Time `json:"fecha_pago" gorm:"column:fecha_pago"`

	CompraStatus        string  `json:"compra_status" gorm:"column:compra_status"`
	FechaSeleccionada   string  `json:"fecha_seleccionada" gorm:"column:fecha_seleccionada"`
	HorarioSeleccionado *string `json:"horario_seleccionado,omitempty" gorm:"column:horario_seleccionado"`
	TipoCompra          string  `json:"tipo_compra" gorm:"column:tipo_compra"`
	TotalParticipantes  int     `json:"total_participantes" gorm:"column:total_participantes"`
	PrecioTotal         float64 `json:"precio_total" gorm:"column:precio_total"`

	PaqueteID           uint    `json:"paquete_id" gorm:"column:paquete_id"`
	PaqueteNombre       string  `json:"paquete_nombre" gorm:"column:paquete_nombre"`
	PaqueteFrecuencia   string  `json:"paquete_frecuencia" gorm:"column:paquete_frecuencia"`
	PaqueteDuracionDias *int    `json:"paquete_duracion_dias,omitempty" gorm:"column:paquete_duracion_dias"`
	PaqueteHorario      *string `json:"paquete_horario,omitempty" gorm:"column:paquete_horario"`

	TuristaID              uint   `json:"turista_id" gorm:"column:turista_id"`
	TuristaNombre          string `json:"turista_nombre" gorm:"column:turista_nombre"`
	TuristaApellidoPaterno string `json:"turista_apellido_paterno" gorm:"column:turista_apellido_paterno"`
	TuristaApellidoMaterno string `json:"turista_apellido_materno" gorm:"column:turista_apellido_materno"`
	TuristaPhone           string `json:"turista_phone" gorm:"column:turista_phone"`
	TuristaEmail           string `json:"turista_email" gorm:"column:turista_email"`
}

type agenciaVentaSalidaRow struct {
	SalidaID uint `json:"salida_id" gorm:"column:salida_id"`

	PaqueteID           uint    `json:"paquete_id" gorm:"column:paquete_id"`
	PaqueteNombre       string  `json:"paquete_nombre" gorm:"column:paquete_nombre"`
	PaqueteFrecuencia   string  `json:"paquete_frecuencia" gorm:"column:paquete_frecuencia"`
	PaqueteDuracionDias *int    `json:"paquete_duracion_dias,omitempty" gorm:"column:paquete_duracion_dias"`
	PaqueteHorario      *string `json:"paquete_horario,omitempty" gorm:"column:paquete_horario"`

	FechaSalida      string `json:"fecha_salida" gorm:"column:fecha_salida"`
	TipoSalida       string `json:"tipo_salida" gorm:"column:tipo_salida"`
	CupoMinimo       int    `json:"cupo_minimo" gorm:"column:cupo_minimo"`
	CupoMaximo       int    `json:"cupo_maximo" gorm:"column:cupo_maximo"`
	CuposReservados  int    `json:"cupos_reservados" gorm:"column:cupos_reservados"`
	CuposConfirmados int    `json:"cupos_confirmados" gorm:"column:cupos_confirmados"`
	Estado           string `json:"estado" gorm:"column:estado"`
	UpdatedAt        string `json:"updated_at" gorm:"column:updated_at"`
}

type agenciaVentaSalidaCompraRow struct {
	CompraID            uint      `json:"compra_id" gorm:"column:compra_id"`
	CompraStatus        string    `json:"compra_status" gorm:"column:compra_status"`
	FechaCompra         time.Time `json:"fecha_compra" gorm:"column:fecha_compra"`
	FechaSeleccionada   string    `json:"fecha_seleccionada" gorm:"column:fecha_seleccionada"`
	HorarioSeleccionado *string   `json:"horario_seleccionado,omitempty" gorm:"column:horario_seleccionado"`
	TipoCompra          string    `json:"tipo_compra" gorm:"column:tipo_compra"`
	TotalParticipantes  int       `json:"total_participantes" gorm:"column:total_participantes"`
	PrecioTotal         float64   `json:"precio_total" gorm:"column:precio_total"`
	NotasTurista        *string   `json:"notas_turista,omitempty" gorm:"column:notas_turista"`

	TuristaID              uint   `json:"turista_id" gorm:"column:turista_id"`
	TuristaNombre          string `json:"turista_nombre" gorm:"column:turista_nombre"`
	TuristaApellidoPaterno string `json:"turista_apellido_paterno" gorm:"column:turista_apellido_paterno"`
	TuristaApellidoMaterno string `json:"turista_apellido_materno" gorm:"column:turista_apellido_materno"`
	TuristaPhone           string `json:"turista_phone" gorm:"column:turista_phone"`
	TuristaEmail           string `json:"turista_email" gorm:"column:turista_email"`

	PagoID          *uint      `json:"pago_id,omitempty" gorm:"column:pago_id"`
	MetodoPago      *string    `json:"metodo_pago,omitempty" gorm:"column:metodo_pago"`
	Monto           *float64   `json:"monto,omitempty" gorm:"column:monto"`
	EstadoPago      *string    `json:"estado,omitempty" gorm:"column:estado"`
	ComprobanteFoto *string    `json:"comprobante_foto,omitempty" gorm:"column:comprobante_foto"`
	FechaPago       *time.Time `json:"fecha_pago,omitempty" gorm:"column:fecha_pago"`

	ConfirmadoPor     *uint      `json:"confirmado_por,omitempty" gorm:"column:confirmado_por"`
	FechaConfirmacion *time.Time `json:"fecha_confirmacion,omitempty" gorm:"column:fecha_confirmacion"`
	RazonRechazo      *string    `json:"razon_rechazo,omitempty" gorm:"column:razon_rechazo"`
	NotasEncargado    *string    `json:"notas_encargado,omitempty" gorm:"column:notas_encargado"`
}

// GetAgenciaVentasPagos lista pagos (ventas) de paquetes asociados a una agencia.
func (h *AgenciaHandler) GetAgenciaVentasPagos(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id64).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para ver ventas de esta agencia", nil, http.StatusForbidden)
		return
	}

	estado := r.URL.Query().Get("estado")
	if estado != "" && estado != "pendiente" && estado != "confirmado" && estado != "rechazado" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "estado invalido (use pendiente|confirmado|rechazado)", nil, http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := database.GetDB()
	base := db.Table("pagos_compras pc").
		Joins("JOIN compras_paquetes c ON c.id = pc.compra_id").
		Joins("JOIN paquetes_turisticos p ON p.id = c.paquete_id").
		Joins("JOIN usuarios u ON u.id = c.turista_id").
		Where("p.agencia_id = ?", agencia.ID)

	if estado != "" {
		base = base.Where("pc.estado = ?", estado)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al contar pagos", err.Error(), http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * pageSize
	var rows []agenciaVentaPagoRow
	if err := base.Session(&gorm.Session{}).
		Select(`
			pc.id AS pago_id,
			pc.compra_id,
			pc.metodo_pago,
			pc.monto,
			pc.estado,
			pc.comprobante_foto,
			pc.confirmado_por,
			pc.fecha_confirmacion,
			pc.razon_rechazo,
			pc.notas_encargado,
			pc.created_at AS fecha_pago,
			c.status AS compra_status,
			c.fecha_seleccionada,
			c.horario_seleccionado,
			c.tipo_compra,
			c.total_participantes,
			c.precio_total,
			p.id AS paquete_id,
			p.nombre AS paquete_nombre,
			p.frecuencia AS paquete_frecuencia,
			p.duracion_dias AS paquete_duracion_dias,
			p.horario AS paquete_horario,
			u.id AS turista_id,
			u.nombre AS turista_nombre,
			u.apellido_paterno AS turista_apellido_paterno,
			u.apellido_materno AS turista_apellido_materno,
			u.phone AS turista_phone,
			u.email AS turista_email
		`).
		Order("pc.created_at DESC").
		Order("pc.id DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&rows).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener pagos", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"pagos": rows,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}, "Ventas obtenidas exitosamente", http.StatusOK)
}

// GetAgenciaVentasSalidas lista salidas habilitadas (fechas) de paquetes de una agencia para visualizar cupos.
func (h *AgenciaHandler) GetAgenciaVentasSalidas(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id64).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para ver ventas de esta agencia", nil, http.StatusForbidden)
		return
	}

	db := database.GetDB()
	base := db.Table("paquete_salidas_habilitadas s").
		Joins("JOIN paquetes_turisticos p ON p.id = s.paquete_id").
		Where("p.agencia_id = ?", agencia.ID)

	if value := strings.TrimSpace(r.URL.Query().Get("paquete_id")); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "paquete_id invalido", nil, http.StatusBadRequest)
			return
		}
		base = base.Where("p.id = ?", uint(parsed))
	}

	if value := strings.TrimSpace(r.URL.Query().Get("estado")); value != "" {
		estado := strings.ToLower(value)
		if !allowedSalidaEstado[estado] {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "estado invalido", nil, http.StatusBadRequest)
			return
		}
		base = base.Where("s.estado = ?", estado)
	}

	if value := strings.TrimSpace(r.URL.Query().Get("desde")); value != "" {
		desde, err := time.Parse("2006-01-02", value)
		if err != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "desde invalido (YYYY-MM-DD)", nil, http.StatusBadRequest)
			return
		}
		base = base.Where("s.fecha_salida >= ?", desde)
	}

	if value := strings.TrimSpace(r.URL.Query().Get("hasta")); value != "" {
		hasta, err := time.Parse("2006-01-02", value)
		if err != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "hasta invalido (YYYY-MM-DD)", nil, http.StatusBadRequest)
			return
		}
		base = base.Where("s.fecha_salida < ?", hasta.AddDate(0, 0, 1))
	}

	var rows []agenciaVentaSalidaRow
	if err := base.
		Select(`
			s.id AS salida_id,
			s.paquete_id,
			p.nombre AS paquete_nombre,
			p.frecuencia AS paquete_frecuencia,
			p.duracion_dias AS paquete_duracion_dias,
			p.horario AS paquete_horario,
			s.fecha_salida,
			s.tipo_salida,
			s.cupo_minimo,
			s.cupo_maximo,
			s.cupos_reservados,
			s.cupos_confirmados,
			s.estado,
			s.updated_at
		`).
		Order("p.nombre ASC").
		Order("s.fecha_salida ASC").
		Order("s.tipo_salida ASC").
		Order("s.id ASC").
		Scan(&rows).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener salidas", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"salidas": rows,
	}, "Salidas obtenidas exitosamente", http.StatusOK)
}

// GetAgenciaVentasSalidaCompras lista compras (y su ultimo pago) asociadas a una salida para gestionar confirmaciones.
func (h *AgenciaHandler) GetAgenciaVentasSalidaCompras(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}
	salidaID64, err := strconv.ParseUint(vars["salida_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de salida invalido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID64).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para ver ventas de esta agencia", nil, http.StatusForbidden)
		return
	}

	db := database.GetDB()
	var salida models.PaqueteSalidaHabilitada
	if err := db.Table("paquete_salidas_habilitadas s").
		Joins("JOIN paquetes_turisticos p ON p.id = s.paquete_id").
		Where("s.id = ? AND p.agencia_id = ?", uint(salidaID64), agencia.ID).
		Select("s.*").
		First(&salida).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, "NOT_FOUND", "Salida no encontrada", nil, http.StatusNotFound)
			return
		}
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener salida", err.Error(), http.StatusInternalServerError)
		return
	}

	var compras []agenciaVentaSalidaCompraRow
	if err := db.Raw(`
		SELECT
			c.id AS compra_id,
			c.status AS compra_status,
			c.fecha_compra,
			c.fecha_seleccionada,
			c.horario_seleccionado,
			c.tipo_compra,
			c.total_participantes,
			c.precio_total,
			c.notas_turista,
			u.id AS turista_id,
			u.nombre AS turista_nombre,
			u.apellido_paterno AS turista_apellido_paterno,
			u.apellido_materno AS turista_apellido_materno,
			u.phone AS turista_phone,
			u.email AS turista_email,
			pc.id AS pago_id,
			pc.metodo_pago,
			pc.monto,
			pc.estado,
			pc.comprobante_foto,
			pc.created_at AS fecha_pago,
			pc.confirmado_por,
			pc.fecha_confirmacion,
			pc.razon_rechazo,
			pc.notas_encargado
		FROM compras_paquetes c
		JOIN paquetes_turisticos p ON p.id = c.paquete_id
		JOIN usuarios u ON u.id = c.turista_id
		LEFT JOIN LATERAL (
			SELECT *
			FROM pagos_compras
			WHERE compra_id = c.id
			ORDER BY created_at DESC, id DESC
			LIMIT 1
		) pc ON TRUE
		WHERE c.salida_id = ? AND p.agencia_id = ?
		ORDER BY c.fecha_compra DESC, c.id DESC
	`, uint(salidaID64), agencia.ID).Scan(&compras).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener compras", err.Error(), http.StatusInternalServerError)
		return
	}

	var totalIngresos float64
	var totalParticipantes int
	for _, compra := range compras {
		totalIngresos += compra.PrecioTotal
		totalParticipantes += compra.TotalParticipantes
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"salida":  salida,
		"compras": compras,
		"totales": map[string]interface{}{
			"ingresos":      totalIngresos,
			"participantes": totalParticipantes,
			"compras":       len(compras),
		},
	}, "Detalle de salida obtenido exitosamente", http.StatusOK)
}
