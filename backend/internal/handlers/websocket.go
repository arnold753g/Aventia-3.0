package handlers

import (
	"log"
	"net/http"

	"andaria-backend/internal/websocket"
	"andaria-backend/pkg/utils"

	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Permitir todas las conexiones (CORS ya se maneja en el middleware)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler maneja las conexiones WebSocket
type WebSocketHandler struct {
	hub *websocket.Hub
}

// NewWebSocketHandler crea un nuevo handler de WebSocket
func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HandleWebSocket maneja las conexiones WebSocket
// WS /api/v1/ws?token=<jwt_token>
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Obtener token del query parameter
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		log.Println("WebSocket: token no proporcionado")
		http.Error(w, "Token requerido", http.StatusUnauthorized)
		return
	}

	// Validar el token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		log.Printf("WebSocket: token inválido - %v", err)
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP connection a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al hacer upgrade a WebSocket: %v", err)
		return
	}

	// Servir la conexión WebSocket
	log.Printf("✅ Nueva conexión WebSocket: usuario_id=%d, email=%s", claims.UserID, claims.Email)
	websocket.ServeWs(h.hub, conn, claims.UserID)
}
