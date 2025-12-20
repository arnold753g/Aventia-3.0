package models

import "time"

// PagoCompra representa un pago asociado a una compra de paquete.
// Tabla: pagos_compras
type PagoCompra struct {
	ID uint `gorm:"primaryKey" json:"id"`

	CompraID uint           `gorm:"not null;index" json:"compra_id"`
	Compra   *CompraPaquete `gorm:"foreignKey:CompraID" json:"compra,omitempty"`

	MetodoPago string  `gorm:"size:20;not null" json:"metodo_pago"`
	Monto      float64 `gorm:"type:decimal(10,2);not null" json:"monto"`

	ComprobanteFoto *string `gorm:"type:text" json:"comprobante_foto,omitempty"`

	Estado string `gorm:"size:20;default:'pendiente';index" json:"estado"`

	ConfirmadoPor     *uint      `gorm:"index" json:"confirmado_por,omitempty"`
	FechaConfirmacion *time.Time `json:"fecha_confirmacion,omitempty"`
	RazonRechazo      *string    `gorm:"type:text" json:"razon_rechazo,omitempty"`
	NotasEncargado    *string    `gorm:"type:text" json:"notas_encargado,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PagoCompra) TableName() string {
	return "pagos_compras"
}
