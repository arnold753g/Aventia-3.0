package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"andaria-backend/internal/models"
	"andaria-backend/internal/websocket"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NotificationListener escucha notificaciones de PostgreSQL
type NotificationListener struct {
	connPool *pgxpool.Pool
	hub      *websocket.Hub
	ctx      context.Context
	cancel   context.CancelFunc
}

// NotificationPayload estructura del payload de NOTIFY
type NotificationPayload struct {
	UsuarioID      uint `json:"usuario_id"`
	NotificacionID uint `json:"notificacion_id"`
}

// NewNotificationListener crea una nueva instancia del listener
func NewNotificationListener(connPool *pgxpool.Pool, hub *websocket.Hub) *NotificationListener {
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationListener{
		connPool: connPool,
		hub:      hub,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start inicia el listener de notificaciones
func (nl *NotificationListener) Start() error {
	// Obtener una conexión del pool para el listener
	conn, err := nl.connPool.Acquire(nl.ctx)
	if err != nil {
		return fmt.Errorf("error al obtener conexión del pool: %w", err)
	}

	// Liberar la conexión cuando terminemos
	defer conn.Release()

	// LISTEN al canal de notificaciones
	_, err = conn.Exec(nl.ctx, "LISTEN notificaciones")
	if err != nil {
		return fmt.Errorf("error al ejecutar LISTEN: %w", err)
	}

	log.Println("✅ PostgreSQL Listener iniciado - escuchando canal 'notificaciones'")

	// Loop infinito para escuchar notificaciones
	for {
		select {
		case <-nl.ctx.Done():
			log.Println("PostgreSQL Listener detenido")
			return nil

		default:
			// Esperar notificación con timeout
			notification, err := conn.Conn().WaitForNotification(nl.ctx)
			if err != nil {
				// Si es contexto cancelado, salir gracefully
				if nl.ctx.Err() != nil {
					return nil
				}
				log.Printf("Error esperando notificación: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Procesar la notificación
			nl.handleNotification(notification.Payload)
		}
	}
}

// handleNotification procesa una notificación recibida
func (nl *NotificationListener) handleNotification(payload string) {
	var notifPayload NotificationPayload

	// Parsear el payload JSON
	if err := json.Unmarshal([]byte(payload), &notifPayload); err != nil {
		log.Printf("Error al parsear payload de notificación: %v", err)
		return
	}

	log.Printf("📬 Notificación recibida: usuario_id=%d, notificacion_id=%d",
		notifPayload.UsuarioID, notifPayload.NotificacionID)

	// Obtener la notificación completa de la base de datos
	notif, err := nl.getNotificacion(notifPayload.NotificacionID)
	if err != nil {
		log.Printf("Error al obtener notificación %d: %v", notifPayload.NotificacionID, err)
		return
	}

	// Enviar vía WebSocket al usuario
	nl.hub.EnviarAUsuario(notifPayload.UsuarioID, notif)
}

// getNotificacion obtiene una notificación de la base de datos
func (nl *NotificationListener) getNotificacion(notificacionID uint) (*models.Notificacion, error) {
	var notif models.Notificacion

	query := `
		SELECT id, usuario_id, tipo, titulo, mensaje, datos_json, leida, fecha_leida, created_at
		FROM notificaciones
		WHERE id = $1
	`

	row := nl.connPool.QueryRow(nl.ctx, query, notificacionID)

	err := row.Scan(
		&notif.ID,
		&notif.UsuarioID,
		&notif.Tipo,
		&notif.Titulo,
		&notif.Mensaje,
		&notif.DatosJSON,
		&notif.Leida,
		&notif.FechaLeida,
		&notif.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error al escanear notificación: %w", err)
	}

	return &notif, nil
}

// Stop detiene el listener
func (nl *NotificationListener) Stop() {
	nl.cancel()
}
