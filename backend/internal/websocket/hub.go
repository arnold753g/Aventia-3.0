package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"andaria-backend/internal/models"
)

// Hub mantiene el conjunto de clientes activos y transmite mensajes
type Hub struct {
	// Clientes registrados
	clients map[*Client]bool

	// Mensajes de broadcast desde los clientes
	broadcast chan []byte

	// Registrar requests de los clientes
	register chan *Client

	// Unregister requests de los clientes
	unregister chan *Client

	// Mapa de usuarios a sus conexiones WebSocket
	usuarios map[uint][]*Client

	// Mutex para acceso concurrente seguro
	mu sync.RWMutex
}

// NewHub crea una nueva instancia de Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		usuarios:   make(map[uint][]*Client),
	}
}

// Run inicia el loop principal del hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			// Agregar cliente al mapa de usuarios
			h.usuarios[client.UsuarioID] = append(h.usuarios[client.UsuarioID], client)
			h.mu.Unlock()
			log.Printf("Cliente registrado: usuario_id=%d, total_conexiones=%d", client.UsuarioID, len(h.usuarios[client.UsuarioID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remover cliente del mapa de usuarios
				if conns, ok := h.usuarios[client.UsuarioID]; ok {
					// Filtrar este cliente
					newConns := make([]*Client, 0)
					for _, c := range conns {
						if c != client {
							newConns = append(newConns, c)
						}
					}
					if len(newConns) > 0 {
						h.usuarios[client.UsuarioID] = newConns
					} else {
						delete(h.usuarios, client.UsuarioID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Cliente desregistrado: usuario_id=%d", client.UsuarioID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// EnviarAUsuario envía una notificación a un usuario específico
func (h *Hub) EnviarAUsuario(usuarioID uint, notif *models.Notificacion) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns, existe := h.usuarios[usuarioID]
	if !existe || len(conns) == 0 {
		log.Printf("Usuario %d no está conectado, notificación guardada en BD", usuarioID)
		return
	}

	// Convertir notificación a JSON
	message, err := json.Marshal(notif.ToDTO())
	if err != nil {
		log.Printf("Error al serializar notificación: %v", err)
		return
	}

	// Enviar a todas las conexiones del usuario
	for _, client := range conns {
		select {
		case client.send <- message:
			log.Printf("Notificación enviada vía WebSocket a usuario %d", usuarioID)
		default:
			// El canal está lleno, cerrar cliente
			close(client.send)
			delete(h.clients, client)
			log.Printf("Cliente desconectado por canal lleno: usuario_id=%d", usuarioID)
		}
	}
}

// ContarClientes retorna el número de clientes conectados
func (h *Hub) ContarClientes() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// ContarUsuarios retorna el número de usuarios únicos conectados
func (h *Hub) ContarUsuarios() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.usuarios)
}
