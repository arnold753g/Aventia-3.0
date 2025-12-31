package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Notificacion representa una notificación en el sistema
type Notificacion struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UsuarioID  uint           `gorm:"not null;index:idx_notificaciones_usuario_leida" json:"usuario_id"`
	Tipo       string         `gorm:"size:50;not null" json:"tipo"` // nuevo_pago_pendiente, pago_confirmado, pago_rechazado, compra_expirada
	Titulo     string         `gorm:"size:255;not null" json:"titulo"`
	Mensaje    string         `gorm:"type:text;not null" json:"mensaje"`
	DatosJSON  NotifDatosJSON `gorm:"type:jsonb" json:"datos_json"`
	Leida      bool           `gorm:"default:false;index:idx_notificaciones_usuario_leida" json:"leida"`
	FechaLeida *time.Time     `json:"fecha_leida,omitempty"`
	CreatedAt  time.Time      `gorm:"index:idx_notificaciones_created_at" json:"created_at"`

	// Relación
	Usuario *Usuario `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
}

// NotifDatosJSON maneja el campo JSONB de datos
type NotifDatosJSON map[string]interface{}

// Implementar interfaz Scanner para leer desde la BD
func (n *NotifDatosJSON) Scan(value interface{}) error {
	if value == nil {
		*n = make(NotifDatosJSON)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, n)
}

// Implementar interfaz Valuer para escribir a la BD
func (n NotifDatosJSON) Value() (driver.Value, error) {
	if n == nil {
		return nil, nil
	}
	return json.Marshal(n)
}

// TableName especifica el nombre de la tabla
func (Notificacion) TableName() string {
	return "notificaciones"
}

// BeforeCreate hook de GORM
func (n *Notificacion) BeforeCreate(tx *gorm.DB) error {
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	return nil
}

// MarcarComoLeida marca la notificación como leída
func (n *Notificacion) MarcarComoLeida(db *gorm.DB) error {
	now := time.Now()
	n.Leida = true
	n.FechaLeida = &now
	return db.Model(n).Updates(map[string]interface{}{
		"leida":       true,
		"fecha_leida": now,
	}).Error
}

// Constantes para tipos de notificaciones
const (
	TipoNuevoPagoPendiente = "nuevo_pago_pendiente"
	TipoPagoConfirmado     = "pago_confirmado"
	TipoPagoRechazado      = "pago_rechazado"
	TipoCompraExpirada     = "compra_expirada"
)
