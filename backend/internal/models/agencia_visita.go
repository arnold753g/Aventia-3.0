package models

import "time"

// AgenciaVisita representa una visita a la página pública de una agencia
type AgenciaVisita struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AgenciaID   uint      `gorm:"not null;index" json:"agencia_id"`
	FechaVisita time.Time `gorm:"not null;index;default:CURRENT_TIMESTAMP" json:"fecha_visita"`
	IPAddress   string    `gorm:"size:50" json:"ip_address,omitempty"`
	UserAgent   string    `gorm:"type:text" json:"user_agent,omitempty"`
	Referer     string    `gorm:"size:500" json:"referer,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// AgenciaEstadisticasVisitas representa estadísticas de visitas de una agencia
type AgenciaEstadisticasVisitas struct {
	AgenciaID             uint      `json:"agencia_id"`
	NombreComercial       string    `json:"nombre_comercial"`
	Slug                  string    `json:"slug"`
	TotalVisitas          int64     `json:"total_visitas"`
	DiasConVisitas        int64     `json:"dias_con_visitas"`
	UltimaVisita          *time.Time `json:"ultima_visita"`
	VisitasUltimaSemana   int64     `json:"visitas_ultima_semana"`
	VisitasUltimoMes      int64     `json:"visitas_ultimo_mes"`
}

func (AgenciaVisita) TableName() string {
	return "agencia_visitas"
}
