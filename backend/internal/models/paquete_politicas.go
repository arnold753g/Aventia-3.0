package models

import "time"

// PaquetePolitica representa la configuración de políticas de paquetes por agencia.
// Aplica a todos los paquetes de la agencia.
type PaquetePolitica struct {
	ID uint `gorm:"primaryKey" json:"id"`

	AgenciaID uint `gorm:"uniqueIndex;not null" json:"agencia_id"`

	// Política de niños: menores a esta edad no pagan
	EdadMinimaPago int `gorm:"default:6" json:"edad_minima_pago"`

	// Porcentaje de recargo para paquetes privados
	RecargoPrivadoPorcentaje float64 `gorm:"type:decimal(5,2);default:0" json:"recargo_privado_porcentaje"`

	// Política de cancelación (texto libre)
	PoliticaCancelacion *string `gorm:"type:text" json:"politica_cancelacion"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PaquetePolitica) TableName() string {
	return "paquete_politicas"
}

