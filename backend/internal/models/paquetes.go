package models

import "time"

// PaqueteTuristico representa la plantilla/oferta de un paquete turístico de una agencia.
// Las salidas reales con fecha se registran en paquete_salidas_habilitadas.
type PaqueteTuristico struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	AgenciaID uint `gorm:"not null;index" json:"agencia_id"`

	Agencia *AgenciaTurismo `gorm:"foreignKey:AgenciaID" json:"agencia,omitempty"`

	// Información básica
	Nombre      string  `gorm:"size:255;not null" json:"nombre"`
	Descripcion *string `gorm:"type:text" json:"descripcion"`

	// Tipo de paquete
	Frecuencia string `gorm:"size:20;not null" json:"frecuencia"`

	// Duración
	DuracionDias   *int `json:"duracion_dias"`
	DuracionNoches *int `json:"duracion_noches"`

	// Fecha de salida fija (solo salida_unica)
	FechaSalidaFija *string `gorm:"type:date" json:"fecha_salida_fija"`

	// Horarios (solo paquetes de 1 día)
	Horario       *string `gorm:"size:20" json:"horario"`
	HoraSalida    *string `gorm:"type:time" json:"hora_salida"`
	DuracionHoras *string `gorm:"type:interval" json:"duracion_horas"`

	// Logística
	DiasPreviosCompra int `gorm:"default:1" json:"dias_previos_compra"`

	// Características
	NivelDificultad *string `gorm:"size:20" json:"nivel_dificultad"`

	// Capacidad
	CupoMinimo     int  `gorm:"not null" json:"cupo_minimo"`
	CupoMaximo     int  `gorm:"not null" json:"cupo_maximo"`
	PermitePrivado bool `gorm:"default:true" json:"permite_privado"`

	// Precios
	PrecioBaseNacionales       float64 `gorm:"type:decimal(10,2);not null" json:"precio_base_nacionales"`
	PrecioAdicionalExtranjeros float64 `gorm:"type:decimal(10,2);default:0" json:"precio_adicional_extranjeros"`

	// Incluye / No incluye / Qué llevar
	Incluye   StringArray `gorm:"type:text[]" json:"incluye"`
	NoIncluye StringArray `gorm:"type:text[]" json:"no_incluye"`
	QueLlevar StringArray `gorm:"type:text[]" json:"que_llevar"`

	// Estado
	Status         string `gorm:"size:20;default:'borrador';index" json:"status"`
	VisiblePublico bool   `gorm:"default:true;index" json:"visible_publico"`

	// Relaciones (opcionales en responses)
	Fotos       []PaqueteFoto             `gorm:"foreignKey:PaqueteID" json:"fotos,omitempty"`
	Itinerario  []PaqueteItinerario       `gorm:"foreignKey:PaqueteID" json:"itinerario,omitempty"`
	Atracciones []PaqueteAtraccion        `gorm:"foreignKey:PaqueteID" json:"atracciones,omitempty"`
	Salidas     []PaqueteSalidaHabilitada `gorm:"foreignKey:PaqueteID" json:"salidas,omitempty"`

	// Campos calculados (no persistidos)
	Politicas        *PaquetePolitica  `gorm:"-" json:"politicas,omitempty"`
	AgenciaDatosPago *AgenciaDatosPago `gorm:"-" json:"agencia_datos_pago,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PaqueteTuristico) TableName() string {
	return "paquetes_turisticos"
}

// PaqueteSalidaHabilitada representa una salida real de un paquete con fecha específica.
type PaqueteSalidaHabilitada struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	PaqueteID uint `gorm:"not null;index" json:"paquete_id"`

	FechaSalida string `gorm:"type:date;not null" json:"fecha_salida"`
	TipoSalida  string `gorm:"size:20;not null" json:"tipo_salida"`

	CupoMinimo       int `gorm:"not null" json:"cupo_minimo"`
	CupoMaximo       int `gorm:"not null" json:"cupo_maximo"`
	CuposReservados  int `gorm:"default:0" json:"cupos_reservados"`
	CuposConfirmados int `gorm:"default:0" json:"cupos_confirmados"`

	PuntoEncuentro        *string `gorm:"type:text" json:"punto_encuentro"`
	HoraEncuentro         *string `gorm:"type:time" json:"hora_encuentro"`
	NotasLogistica        *string `gorm:"type:text" json:"notas_logistica"`
	InstruccionesTuristas *string `gorm:"type:text" json:"instrucciones_turistas"`

	GuiaNombre   *string `gorm:"size:255" json:"guia_nombre"`
	GuiaTelefono *string `gorm:"size:20" json:"guia_telefono"`

	Estado           string  `gorm:"size:20;default:'pendiente'" json:"estado"`
	RazonCancelacion *string `gorm:"type:text" json:"razon_cancelacion"`

	// Campos para salidas pre-creadas manualmente
	CreadaManualmente        bool       `gorm:"default:false" json:"creada_manualmente"`
	CreadaPorUsuarioID       *uint      `gorm:"index" json:"creada_por_usuario_id,omitempty"`
	FechaLimiteInscripcion   *time.Time `json:"fecha_limite_inscripcion,omitempty"`
	DescripcionSalida        *string    `gorm:"type:text" json:"descripcion_salida,omitempty"`
	NotasInternas            *string    `gorm:"type:text" json:"notas_internas,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PaqueteSalidaHabilitada) TableName() string {
	return "paquete_salidas_habilitadas"
}

// PaqueteItinerario representa el itinerario día por día para paquetes multi-día.
type PaqueteItinerario struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	PaqueteID uint `gorm:"not null;index;uniqueIndex:uniq_paquete_itinerario,composite:paquete_itinerario" json:"paquete_id"`
	DiaNumero int  `gorm:"not null;uniqueIndex:uniq_paquete_itinerario,composite:paquete_itinerario" json:"dia_numero"`

	Titulo        string   `gorm:"size:255;not null" json:"titulo"`
	Descripcion   *string  `gorm:"type:text" json:"descripcion"`
	Actividades   StringArray `gorm:"type:text[]" json:"actividades"`
	HospedajeInfo *string  `gorm:"type:text" json:"hospedaje_info"`

	CreatedAt time.Time `json:"created_at"`
}

func (PaqueteItinerario) TableName() string {
	return "paquete_itinerario"
}

// PaqueteFoto representa fotos asociadas a un paquete (máximo 6).
type PaqueteFoto struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	PaqueteID uint `gorm:"not null;index" json:"paquete_id"`
	Foto      string `gorm:"type:text;not null" json:"foto"`

	EsPrincipal bool `gorm:"default:false" json:"es_principal"`
	Orden       int  `gorm:"default:0" json:"orden"`

	CreatedAt time.Time `json:"created_at"`
}

func (PaqueteFoto) TableName() string {
	return "paquete_fotos"
}

// PaqueteAtraccion representa una atracción incluida dentro del paquete.
type PaqueteAtraccion struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	PaqueteID   uint `gorm:"not null;index" json:"paquete_id"`
	AtraccionID uint `gorm:"not null;index" json:"atraccion_id"`

	// Relación opcional para mostrar datos de la atracción
	Atraccion *AtraccionTuristica `gorm:"foreignKey:AtraccionID" json:"atraccion,omitempty"`

	DiaNumero             *int `json:"dia_numero"`
	OrdenVisita           int  `gorm:"not null" json:"orden_visita"`
	DuracionEstimadaHoras *int `json:"duracion_estimada_horas"`

	CreatedAt time.Time `json:"created_at"`
}

func (PaqueteAtraccion) TableName() string {
	return "paquete_atracciones"
}
