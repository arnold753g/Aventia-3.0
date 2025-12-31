package models

import "time"

// CompraPaquete representa una compra de un paquete turÍstico realizada por un turista.
// Tabla: compras_paquetes
type CompraPaquete struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// Codigo unico generado al confirmar el pago.
	CodigoConfirmacion *string `gorm:"size:30;uniqueIndex" json:"codigo_confirmacion,omitempty"`

	// Relaciones
	TuristaID uint     `gorm:"not null;index" json:"turista_id"`
	Turista   *Usuario `gorm:"foreignKey:TuristaID" json:"turista,omitempty"`

	PaqueteID uint              `gorm:"not null;index" json:"paquete_id"`
	Paquete   *PaqueteTuristico `gorm:"foreignKey:PaqueteID" json:"paquete,omitempty"`

	SalidaID *uint                    `gorm:"index" json:"salida_id,omitempty"`
	Salida   *PaqueteSalidaHabilitada `gorm:"foreignKey:SalidaID" json:"salida,omitempty"`

	PromocionID *uint `gorm:"index" json:"promocion_id,omitempty"`

	// Información de la compra
	FechaCompra         time.Time `json:"fecha_compra"`
	FechaSeleccionada   time.Time `gorm:"type:date;not null" json:"fecha_seleccionada"`
	HorarioSeleccionado *string   `gorm:"size:20" json:"horario_seleccionado,omitempty"`

	TipoCompra string `gorm:"size:20;not null" json:"tipo_compra"`

	// Participantes
	Extranjero          bool `gorm:"default:false" json:"extranjero"`
	CantidadAdultos     int  `gorm:"not null;default:1" json:"cantidad_adultos"`
	CantidadNinosPagan  int  `gorm:"default:0" json:"cantidad_ninos_pagan"`
	CantidadNinosGratis int  `gorm:"default:0" json:"cantidad_ninos_gratis"`
	TotalParticipantes  int  `gorm:"not null" json:"total_participantes"`

	// Precios
	PrecioUnitario           float64 `gorm:"type:decimal(10,2);not null" json:"precio_unitario"`
	RecargoPrivadoPorcentaje float64 `gorm:"type:decimal(5,2);default:0" json:"recargo_privado_porcentaje"`
	RecargoExtranjero        float64 `gorm:"type:decimal(10,2);default:0" json:"recargo_extranjero"`
	Subtotal                 float64 `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	TotalRecargo             float64 `gorm:"type:decimal(10,2);default:0" json:"total_recargo"`
	PrecioTotal              float64 `gorm:"type:decimal(10,2);not null" json:"precio_total"`

	// Campos de promoción (para futuro)
	PrecioSinDescuento          *float64 `gorm:"type:decimal(10,2)" json:"precio_sin_descuento,omitempty"`
	DescuentoAplicado           float64  `gorm:"type:decimal(10,2);default:0" json:"descuento_aplicado"`
	PorcentajeDescuentoAplicado float64  `gorm:"type:decimal(5,2);default:0" json:"porcentaje_descuento_aplicado"`

	// Información adicional
	TieneDiscapacidad       bool    `gorm:"default:false" json:"tiene_discapacidad"`
	DescripcionDiscapacidad *string `gorm:"type:text" json:"descripcion_discapacidad,omitempty"`
	NotasTurista            *string `gorm:"type:text" json:"notas_turista,omitempty"`

	// Estado
	Status            string     `gorm:"size:30;default:'pendiente_confirmacion';index" json:"status"`
	FechaConfirmacion *time.Time `json:"fecha_confirmacion,omitempty"`
	FechaRechazo      *time.Time `json:"fecha_rechazo,omitempty"`
	RazonRechazo      *string    `gorm:"type:text" json:"razon_rechazo,omitempty"`

	// Relaciones
	Pagos []PagoCompra `gorm:"foreignKey:CompraID" json:"pagos,omitempty"`

	// Auditoría
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (CompraPaquete) TableName() string {
	return "compras_paquetes"
}
