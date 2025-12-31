package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Tiempo permitido para escribir un mensaje al peer
	writeWait = 10 * time.Second

	// Tiempo permitido para leer el siguiente pong del peer
	pongWait = 60 * time.Second

	// Enviar pings al peer con este período (debe ser menor que pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Tamaño máximo del mensaje permitido del peer
	maxMessageSize = 512
)

// Client es un intermediario entre la conexión websocket y el hub
type Client struct {
	hub *Hub

	// La conexión websocket
	conn *websocket.Conn

	// Canal buffereado de mensajes salientes
	send chan []byte

	// ID del usuario conectado
	UsuarioID uint
}

// readPump bombea mensajes de la conexión websocket al hub
//
// La aplicación ejecuta readPump en una goroutine por conexión. La aplicación
// asegura que hay como máximo un lector en una conexión ejecutando todas las
// lecturas desde esta goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Por ahora, solo registramos los mensajes recibidos
		// En el futuro aquí podríamos manejar acciones del cliente
		log.Printf("Mensaje recibido del cliente usuario_id=%d: %s", c.UsuarioID, message)
	}
}

// writePump bombea mensajes del hub a la conexión websocket
//
// Una goroutine ejecutando writePump se inicia para cada conexión. La
// aplicación asegura que hay como máximo un escritor en una conexión ejecutando
// todas las escrituras desde esta goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// El hub cerró el canal
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Agregar mensajes en cola al mensaje actual
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs maneja las solicitudes websocket de los peers
func ServeWs(hub *Hub, conn *websocket.Conn, usuarioID uint) {
	client := &Client{
		hub:       hub,
		conn:      conn,
		send:      make(chan []byte, 256),
		UsuarioID: usuarioID,
	}

	client.hub.register <- client

	// Permitir la colección de memoria referenciada por el caller haciendo todo el trabajo en
	// nuevas goroutines
	go client.writePump()
	go client.readPump()
}
