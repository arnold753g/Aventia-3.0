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

// CompraDetalleResponse representa el detalle de una compra para el turista.
type CompraDetalleResponse struct {
	ID                 uint                  `json:"id"`
	FechaCompra        time.Time             `json:"fecha_compra"`
	FechaSeleccionada  time.Time             `json:"fecha_seleccionada"`
	TipoCompra         string                `json:"tipo_compra"`
	TotalParticipantes int                   `json:"total_participantes"`
	PrecioTotal        float64               `json:"precio_total"`
	Status             string                `json:"status"`
	Paquete            PaqueteSimpleResponse `json:"paquete"`
	Salida             *SalidaSimpleResponse `json:"salida"`
	UltimoPago         *PagoSimpleResponse   `json:"ultimo_pago"`
}

type PaqueteSimpleResponse struct {
	ID           uint    `json:"id"`
	Nombre       string  `json:"nombre"`
	Frecuencia   string  `json:"frecuencia"`
	DuracionDias *int    `json:"duracion_dias"`
	Horario      *string `json:"horario"`
}

type SalidaSimpleResponse struct {
	ID               uint   `json:"id"`
	FechaSalida      string `json:"fecha_salida"`
	TipoSalida       string `json:"tipo_salida"`
	Estado           string `json:"estado"`
	CuposDisponibles int    `json:"cupos_disponibles"`
}

type PagoSimpleResponse struct {
	ID                uint       `json:"id"`
	MetodoPago        string     `json:"metodo_pago"`
	Monto             float64    `json:"monto"`
	Estado            string     `json:"estado"`
	FechaConfirmacion *time.Time `json:"fecha_confirmacion"`
}
