package models

import "time"

// AgenciaDatosPago representa los datos de pago (transferencia/QR) configurados por una agencia.
type AgenciaDatosPago struct {
	ID uint `gorm:"primaryKey" json:"id"`

	AgenciaID uint `gorm:"uniqueIndex;not null" json:"agencia_id"`

	NombreBanco   *string `gorm:"size:100" json:"nombre_banco"`
	NumeroCuenta  *string `gorm:"size:50" json:"numero_cuenta"`
	NombreTitular *string `gorm:"size:255" json:"nombre_titular"`
	QrPagoFoto    *string `gorm:"type:text" json:"qr_pago_foto"`
	Activo        bool    `gorm:"default:true" json:"activo"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AgenciaDatosPago) TableName() string {
	return "agencia_datos_pago"
}

