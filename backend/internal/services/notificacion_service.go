package services

import (
	"fmt"
	"time"

	"andaria-backend/internal/models"

	"gorm.io/gorm"
)

// NotificacionService maneja la lógica de negocio de notificaciones
type NotificacionService struct {
	db *gorm.DB
}

// NewNotificacionService crea una nueva instancia del servicio
func NewNotificacionService(db *gorm.DB) *NotificacionService {
	return &NotificacionService{db: db}
}

// ObtenerNotificaciones obtiene las notificaciones de un usuario con paginación
func (s *NotificacionService) ObtenerNotificaciones(usuarioID uint, page, limit int) ([]models.Notificacion, int64, error) {
	var notificaciones []models.Notificacion
	var total int64

	// Contar total
	if err := s.db.Model(&models.Notificacion{}).
		Where("usuario_id = ?", usuarioID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener notificaciones paginadas
	offset := (page - 1) * limit
	if err := s.db.
		Where("usuario_id = ?", usuarioID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notificaciones).Error; err != nil {
		return nil, 0, err
	}

	return notificaciones, total, nil
}

// ContarNoLeidas cuenta las notificaciones no leídas de un usuario
func (s *NotificacionService) ContarNoLeidas(usuarioID uint) (int64, error) {
	var count int64
	err := s.db.Model(&models.Notificacion{}).
		Where("usuario_id = ? AND leida = ?", usuarioID, false).
		Count(&count).Error
	return count, err
}

// MarcarComoLeida marca una notificación como leída
func (s *NotificacionService) MarcarComoLeida(notificacionID, usuarioID uint) error {
	now := time.Now()
	result := s.db.Model(&models.Notificacion{}).
		Where("id = ? AND usuario_id = ?", notificacionID, usuarioID).
		Updates(map[string]interface{}{
			"leida":       true,
			"fecha_leida": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notificación no encontrada o no pertenece al usuario")
	}

	return nil
}

// MarcarTodasLeidas marca todas las notificaciones de un usuario como leídas
func (s *NotificacionService) MarcarTodasLeidas(usuarioID uint) (int64, error) {
	now := time.Now()
	result := s.db.Model(&models.Notificacion{}).
		Where("usuario_id = ? AND leida = ?", usuarioID, false).
		Updates(map[string]interface{}{
			"leida":       true,
			"fecha_leida": now,
		})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// EliminarNotificacion elimina una notificación
func (s *NotificacionService) EliminarNotificacion(notificacionID, usuarioID uint) error {
	result := s.db.
		Where("id = ? AND usuario_id = ?", notificacionID, usuarioID).
		Delete(&models.Notificacion{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notificación no encontrada o no pertenece al usuario")
	}

	return nil
}

// CrearNotificacion crea una nueva notificación (usado por el worker de expiración)
func (s *NotificacionService) CrearNotificacion(notif *models.Notificacion) error {
	return s.db.Create(notif).Error
}

// EliminarAntiguas elimina notificaciones leídas antiguas (limpieza periódica)
func (s *NotificacionService) EliminarAntiguas(diasAntiguedad int) (int64, error) {
	fecha := time.Now().AddDate(0, 0, -diasAntiguedad)
	result := s.db.
		Where("leida = ? AND fecha_leida < ?", true, fecha).
		Delete(&models.Notificacion{})

	return result.RowsAffected, result.Error
}
