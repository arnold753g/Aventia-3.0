package models

import "time"

// CrearCompraRequest representa el body JSON para crear una compra de paquete.
// Nota: turista_id se obtiene del JWT (no se acepta desde el cliente).
type CrearCompraRequest struct {
	PaqueteID         uint   `json:"paquete_id" validate:"required"`
	FechaSeleccionada string `json:"fecha_seleccionada" validate:"required"` // YYYY-MM-DD

	TipoCompra string `json:"tipo_compra" validate:"required,oneof=compartido privado"`

	Extranjero          bool `json:"extranjero"`
	CantidadAdultos     int  `json:"cantidad_adultos" validate:"required,min=1"`
	CantidadNinosPagan  int  `json:"cantidad_ninos_pagan" validate:"min=0"`
	CantidadNinosGratis int  `json:"cantidad_ninos_gratis" validate:"min=0"`

	TieneDiscapacidad       bool    `json:"tiene_discapacidad"`
	DescripcionDiscapacidad *string `json:"descripcion_discapacidad"`
	NotasTurista            *string `json:"notas_turista"`
}

// ProcesarCompraPaqueteResult representa el resultado retornado por la función SQL procesar_compra_paquete().
type ProcesarCompraPaqueteResult struct {
	CompraID    uint    `json:"compra_id" gorm:"column:compra_id"`
	SalidaID    uint    `json:"salida_id" gorm:"column:salida_id"`
	PrecioTotal float64 `json:"precio_total" gorm:"column:precio_total"`
	Mensaje     string  `json:"mensaje" gorm:"column:mensaje"`
	Success     bool    `json:"success" gorm:"column:success"`
}

// ParticipanteDetalle representa un participante individual (opcionalmente detallado).
type ParticipanteDetalle struct {
	Nombre *string `json:"nombre,omitempty"`
	Edad   *int    `json:"edad,omitempty"`
	Tipo   string  `json:"tipo"`
	Notas  *string `json:"notas,omitempty"`
}

// ParticipantesDetalle agrupa participantes por tipo.
type ParticipantesDetalle struct {
	Adultos     []ParticipanteDetalle `json:"adultos,omitempty"`
	NinosPagan  []ParticipanteDetalle `json:"ninos_pagan,omitempty"`
	NinosGratis []ParticipanteDetalle `json:"ninos_gratis,omitempty"`
}

// DepartamentoSimpleResponse expone solo el nombre del departamento.
type DepartamentoSimpleResponse struct {
	Nombre string `json:"nombre"`
}

// AgenciaCompraResponse expone datos de contacto de la agencia para la compra.
type AgenciaCompraResponse struct {
	ID              uint                       `json:"id"`
	NombreComercial string                     `json:"nombre_comercial"`
	Direccion       string                     `json:"direccion"`
	Departamento    *DepartamentoSimpleResponse `json:"departamento,omitempty"`
	Telefono        string                     `json:"telefono"`
	Email           string                     `json:"email"`
	Whatsapp        *string                    `json:"whatsapp,omitempty"`
}

// PaqueteDetalleResponse representa el paquete con datos completos para la compra.
type PaqueteDetalleResponse struct {
	ID               uint               `json:"id"`
	Nombre           string             `json:"nombre"`
	Descripcion      *string            `json:"descripcion,omitempty"`
	Frecuencia       string             `json:"frecuencia"`
	DuracionDias     *int               `json:"duracion_dias,omitempty"`
	Horario          *string            `json:"horario,omitempty"`
	HoraSalida       *string            `json:"hora_salida,omitempty"`
	NivelDificultad  *string            `json:"nivel_dificultad,omitempty"`
	PermitePrivado   bool               `json:"permite_privado"`
	Incluye          StringArray        `json:"incluye,omitempty"`
	NoIncluye        StringArray        `json:"no_incluye,omitempty"`
	QueLlevar        StringArray        `json:"que_llevar,omitempty"`
	Fotos            []PaqueteFoto      `json:"fotos,omitempty"`
	Itinerario       []PaqueteItinerario `json:"itinerario,omitempty"`
	Atracciones      []PaqueteAtraccion `json:"atracciones,omitempty"`
	Agencia          *AgenciaCompraResponse `json:"agencia,omitempty"`
	Politicas        *PaquetePolitica       `json:"politicas,omitempty"`
	AgenciaDatosPago *AgenciaDatosPago      `json:"agencia_datos_pago,omitempty"`
}

// CompraDetalleResponse representa el detalle de una compra para el turista.
type CompraDetalleResponse struct {
	ID                     uint                 `json:"id"`
	CodigoConfirmacion     *string              `json:"codigo_confirmacion,omitempty"`
	FechaCompra            time.Time            `json:"fecha_compra"`
	FechaSeleccionada      time.Time            `json:"fecha_seleccionada"`
	FechaConfirmacion      *time.Time           `json:"fecha_confirmacion,omitempty"`
	TipoCompra             string               `json:"tipo_compra"`
	Extranjero             bool                 `json:"extranjero"`
	CantidadAdultos        int                  `json:"cantidad_adultos"`
	CantidadNinosPagan     int                  `json:"cantidad_ninos_pagan"`
	CantidadNinosGratis    int                  `json:"cantidad_ninos_gratis"`
	Participantes          *ParticipantesDetalle `json:"participantes,omitempty"`
	TotalParticipantes     int                  `json:"total_participantes"`
	PrecioTotal            float64              `json:"precio_total"`
	Status                 string               `json:"status"`
	TieneDiscapacidad      bool                 `json:"tiene_discapacidad"`
	DescripcionDiscapacidad *string             `json:"descripcion_discapacidad,omitempty"`
	NotasTurista           *string              `json:"notas_turista,omitempty"`
	Paquete                PaqueteDetalleResponse `json:"paquete"`
	Salida                 *SalidaSimpleResponse `json:"salida"`
	UltimoPago             *PagoSimpleResponse   `json:"ultimo_pago"`
}

type SalidaSimpleResponse struct {
	ID                    uint    `json:"id"`
	FechaSalida           string  `json:"fecha_salida"`
	TipoSalida            string  `json:"tipo_salida"`
	Estado                string  `json:"estado"`
	CuposDisponibles      int     `json:"cupos_disponibles"`
	PuntoEncuentro        *string `json:"punto_encuentro,omitempty"`
	HoraEncuentro         *string `json:"hora_encuentro,omitempty"`
	InstruccionesTuristas *string `json:"instrucciones_turistas,omitempty"`
}

type PagoSimpleResponse struct {
	ID                uint       `json:"id"`
	MetodoPago        string     `json:"metodo_pago"`
	Monto             float64    `json:"monto"`
	Estado            string     `json:"estado"`
	ComprobanteFoto   *string    `json:"comprobante_foto"`
	FechaConfirmacion *time.Time `json:"fecha_confirmacion"`
	FechaRegistro     time.Time  `json:"fecha_registro"`
}
