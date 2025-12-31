package models

import "time"

// NotificacionDTO representa una notificación en respuestas API
type NotificacionDTO struct {
	ID         uint                   `json:"id"`
	Tipo       string                 `json:"tipo"`
	Titulo     string                 `json:"titulo"`
	Mensaje    string                 `json:"mensaje"`
	DatosJSON  map[string]interface{} `json:"datos_json"`
	Leida      bool                   `json:"leida"`
	FechaLeida *time.Time             `json:"fecha_leida,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// Pagination información de paginación
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// NotificacionesResponse respuesta paginada de notificaciones
type NotificacionesResponse struct {
	Notificaciones []NotificacionDTO `json:"notificaciones"`
	NoLeidas       int64             `json:"no_leidas"`
	Pagination     *Pagination       `json:"pagination,omitempty"`
}

// ContadorNotificacionesResponse contador de notificaciones no leídas
type ContadorNotificacionesResponse struct {
	NoLeidas int64 `json:"no_leidas"`
}

// ToNotificacionDTO convierte Notificacion a DTO
func (n *Notificacion) ToDTO() *NotificacionDTO {
	return &NotificacionDTO{
		ID:         n.ID,
		Tipo:       n.Tipo,
		Titulo:     n.Titulo,
		Mensaje:    n.Mensaje,
		DatosJSON:  n.DatosJSON,
		Leida:      n.Leida,
		FechaLeida: n.FechaLeida,
		CreatedAt:  n.CreatedAt,
	}
}

// ToNotificacionDTOs convierte slice de Notificacion a DTOs
func ToNotificacionDTOs(notificaciones []Notificacion) []NotificacionDTO {
	dtos := make([]NotificacionDTO, len(notificaciones))
	for i, n := range notificaciones {
		dtos[i] = *n.ToDTO()
	}
	return dtos
}
