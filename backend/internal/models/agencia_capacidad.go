package models

import "time"

// AgenciaCapacidad representa la capacidad operativa configurada por la agencia.
// Se usa para validar cuántas salidas puede gestionar simultáneamente.
type AgenciaCapacidad struct {
	ID uint `gorm:"primaryKey" json:"id"`

	AgenciaID uint `gorm:"uniqueIndex;not null" json:"agencia_id"`

	MaxSalidasPorDia     int `gorm:"default:5" json:"max_salidas_por_dia"`
	MaxSalidasPorHorario int `gorm:"default:3" json:"max_salidas_por_horario"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AgenciaCapacidad) TableName() string {
	return "agencia_capacidad"
}

