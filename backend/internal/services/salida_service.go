package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"andaria-backend/internal/models"
)

type SalidaService struct {
	db *gorm.DB
}

func NewSalidaService(db *gorm.DB) *SalidaService {
	return &SalidaService{db: db}
}

// CrearSalidaManual crea una nueva salida habilitada manualmente por la agencia
func (s *SalidaService) CrearSalidaManual(agenciaID, usuarioID, paqueteID uint, req CrearSalidaManualRequest) (*models.PaqueteSalidaHabilitada, error) {
	// Validar que el paquete pertenece a la agencia
	var paquete models.PaqueteTuristico
	if err := s.db.First(&paquete, "id = ? AND agencia_id = ?", paqueteID, agenciaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("paquete no encontrado o no pertenece a esta agencia")
		}
		return nil, err
	}

	// Validar fecha
	fechaSalida, err := time.Parse("2006-01-02", req.FechaSalida)
	if err != nil {
		return nil, errors.New("formato de fecha inválido (use YYYY-MM-DD)")
	}

	// La fecha de salida debe ser futura
	if fechaSalida.Before(time.Now().AddDate(0, 0, 1)) {
		return nil, errors.New("la fecha de salida debe ser al menos mañana")
	}

	// Validar cupos
	if req.CupoMaximo <= 0 {
		return nil, errors.New("el cupo máximo debe ser mayor a 0")
	}

	cupoMinimo := req.CupoMinimo
	if cupoMinimo <= 0 {
		cupoMinimo = paquete.CupoMinimo // Usar el del paquete por defecto
	}

	if cupoMinimo > req.CupoMaximo {
		return nil, errors.New("el cupo mínimo no puede ser mayor al cupo máximo")
	}

	// Validar que no exista ya una salida compartida para esa fecha
	var existente models.PaqueteSalidaHabilitada
	err = s.db.First(&existente,
		"paquete_id = ? AND fecha_salida = ? AND tipo_salida = 'compartido' AND estado IN ('pendiente', 'activa')",
		paqueteID, req.FechaSalida).Error

	if err == nil {
		return nil, errors.New("ya existe una salida compartida habilitada para esta fecha")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Crear la salida
	salida := &models.PaqueteSalidaHabilitada{
		PaqueteID:              paqueteID,
		FechaSalida:            req.FechaSalida,
		TipoSalida:             "compartido",
		CupoMinimo:             cupoMinimo,
		CupoMaximo:             req.CupoMaximo,
		CuposReservados:        0,
		CuposConfirmados:       0,
		Estado:                 "pendiente",
		CreadaManualmente:      true,
		CreadaPorUsuarioID:     &usuarioID,
		FechaLimiteInscripcion: req.FechaLimiteInscripcion,
		DescripcionSalida:      req.Descripcion,
		PuntoEncuentro:         req.PuntoEncuentro,
		HoraEncuentro:          req.HoraEncuentro,
		NotasInternas:          req.NotasInternas,
		InstruccionesTuristas:  req.InstruccionesTuristas,
		GuiaNombre:             req.GuiaNombre,
		GuiaTelefono:           req.GuiaTelefono,
	}

	if err := s.db.Create(salida).Error; err != nil {
		return nil, fmt.Errorf("error al crear salida: %w", err)
	}

	return salida, nil
}

// ObtenerSalidasPorPaquete obtiene todas las salidas de un paquete
func (s *SalidaService) ObtenerSalidasPorPaquete(agenciaID, paqueteID uint, filtros SalidaFiltros) ([]models.PaqueteSalidaHabilitada, error) {
	// Validar que el paquete pertenece a la agencia
	var paquete models.PaqueteTuristico
	if err := s.db.First(&paquete, "id = ? AND agencia_id = ?", paqueteID, agenciaID).Error; err != nil {
		return nil, errors.New("paquete no encontrado")
	}

	query := s.db.Where("paquete_id = ?", paqueteID)

	// Aplicar filtros
	if filtros.Estado != "" {
		query = query.Where("estado = ?", filtros.Estado)
	}

	if filtros.FechaDesde != "" {
		query = query.Where("fecha_salida >= ?", filtros.FechaDesde)
	}

	if filtros.FechaHasta != "" {
		query = query.Where("fecha_salida <= ?", filtros.FechaHasta)
	}

	if filtros.SoloManuales {
		query = query.Where("creada_manualmente = ?", true)
	}

	var salidas []models.PaqueteSalidaHabilitada
	if err := query.Order("fecha_salida ASC").Find(&salidas).Error; err != nil {
		return nil, err
	}

	return salidas, nil
}

// ActualizarSalida actualiza una salida existente
func (s *SalidaService) ActualizarSalida(agenciaID, salidaID uint, req ActualizarSalidaRequest) (*models.PaqueteSalidaHabilitada, error) {
	var salida models.PaqueteSalidaHabilitada

	// Buscar la salida y verificar que pertenece a un paquete de la agencia
	err := s.db.Joins("JOIN paquetes_turisticos ON paquete_salidas_habilitadas.paquete_id = paquetes_turisticos.id").
		Where("paquete_salidas_habilitadas.id = ? AND paquetes_turisticos.agencia_id = ?", salidaID, agenciaID).
		First(&salida).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("salida no encontrada")
		}
		return nil, err
	}

	// No permitir modificar salidas con compras confirmadas
	if salida.CuposConfirmados > 0 && req.CupoMaximo != nil && *req.CupoMaximo < salida.CuposConfirmados {
		return nil, errors.New("no se puede reducir el cupo máximo por debajo de los cupos ya confirmados")
	}

	// Actualizar campos permitidos
	updates := make(map[string]interface{})

	if req.CupoMaximo != nil {
		updates["cupo_maximo"] = *req.CupoMaximo
	}

	if req.FechaLimiteInscripcion != nil {
		updates["fecha_limite_inscripcion"] = req.FechaLimiteInscripcion
	}

	if req.Descripcion != nil {
		updates["descripcion_salida"] = req.Descripcion
	}

	if req.PuntoEncuentro != nil {
		updates["punto_encuentro"] = req.PuntoEncuentro
	}

	if req.HoraEncuentro != nil {
		updates["hora_encuentro"] = req.HoraEncuentro
	}

	if req.NotasInternas != nil {
		updates["notas_internas"] = req.NotasInternas
	}

	if req.InstruccionesTuristas != nil {
		updates["instrucciones_turistas"] = req.InstruccionesTuristas
	}

	if req.GuiaNombre != nil {
		updates["guia_nombre"] = req.GuiaNombre
	}

	if req.GuiaTelefono != nil {
		updates["guia_telefono"] = req.GuiaTelefono
	}

	if req.Estado != nil {
		// Validar transiciones de estado
		if err := s.validarTransicionEstado(salida.Estado, *req.Estado); err != nil {
			return nil, err
		}
		updates["estado"] = *req.Estado
	}

	if len(updates) > 0 {
		if err := s.db.Model(&salida).Updates(updates).Error; err != nil {
			return nil, err
		}

		// Recargar para obtener valores actualizados
		s.db.First(&salida, salidaID)
	}

	return &salida, nil
}

// CancelarSalida cancela una salida
func (s *SalidaService) CancelarSalida(agenciaID, salidaID uint, razon string) error {
	var salida models.PaqueteSalidaHabilitada

	err := s.db.Joins("JOIN paquetes_turisticos ON paquete_salidas_habilitadas.paquete_id = paquetes_turisticos.id").
		Where("paquete_salidas_habilitadas.id = ? AND paquetes_turisticos.agencia_id = ?", salidaID, agenciaID).
		First(&salida).Error

	if err != nil {
		return errors.New("salida no encontrada")
	}

	if salida.Estado == "cancelada" {
		return errors.New("la salida ya está cancelada")
	}

	if salida.Estado == "completada" {
		return errors.New("no se puede cancelar una salida completada")
	}

	// TODO: Aquí se debería notificar a los turistas afectados
	// y procesar devoluciones si hay compras confirmadas

	return s.db.Model(&salida).Updates(map[string]interface{}{
		"estado":            "cancelada",
		"razon_cancelacion": razon,
	}).Error
}

// validarTransicionEstado valida que la transición de estado sea válida
func (s *SalidaService) validarTransicionEstado(estadoActual, estadoNuevo string) error {
	transicionesValidas := map[string][]string{
		"pendiente":  {"activa", "cancelada"},
		"activa":     {"completada", "cancelada"},
		"completada": {}, // No se puede cambiar desde completada
		"cancelada":  {}, // No se puede cambiar desde cancelada
	}

	permitidos, existe := transicionesValidas[estadoActual]
	if !existe {
		return errors.New("estado actual inválido")
	}

	for _, permitido := range permitidos {
		if estadoNuevo == permitido {
			return nil
		}
	}

	return fmt.Errorf("no se puede cambiar de '%s' a '%s'", estadoActual, estadoNuevo)
}

// DTOs
type CrearSalidaManualRequest struct {
	FechaSalida            string     `json:"fecha_salida" binding:"required"` // YYYY-MM-DD
	CupoMinimo             int        `json:"cupo_minimo"`
	CupoMaximo             int        `json:"cupo_maximo" binding:"required,min=1"`
	FechaLimiteInscripcion *time.Time `json:"fecha_limite_inscripcion,omitempty"`
	Descripcion            *string    `json:"descripcion,omitempty"`
	PuntoEncuentro         *string    `json:"punto_encuentro,omitempty"`
	HoraEncuentro          *string    `json:"hora_encuentro,omitempty"` // HH:MM
	NotasInternas          *string    `json:"notas_internas,omitempty"`
	InstruccionesTuristas  *string    `json:"instrucciones_turistas,omitempty"`
	GuiaNombre             *string    `json:"guia_nombre,omitempty"`
	GuiaTelefono           *string    `json:"guia_telefono,omitempty"`
}

type ActualizarSalidaRequest struct {
	CupoMaximo             *int       `json:"cupo_maximo,omitempty"`
	FechaLimiteInscripcion *time.Time `json:"fecha_limite_inscripcion,omitempty"`
	Descripcion            *string    `json:"descripcion,omitempty"`
	PuntoEncuentro         *string    `json:"punto_encuentro,omitempty"`
	HoraEncuentro          *string    `json:"hora_encuentro,omitempty"`
	NotasInternas          *string    `json:"notas_internas,omitempty"`
	InstruccionesTuristas  *string    `json:"instrucciones_turistas,omitempty"`
	GuiaNombre             *string    `json:"guia_nombre,omitempty"`
	GuiaTelefono           *string    `json:"guia_telefono,omitempty"`
	Estado                 *string    `json:"estado,omitempty"`
}

type SalidaFiltros struct {
	Estado       string `form:"estado"`
	FechaDesde   string `form:"fecha_desde"` // YYYY-MM-DD
	FechaHasta   string `form:"fecha_hasta"` // YYYY-MM-DD
	SoloManuales bool   `form:"solo_manuales"`
}
