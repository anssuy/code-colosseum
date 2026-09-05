package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Conn struct {
	ws     *websocket.Conn
	send   chan []byte
	userID string
	hub    *Hub
}

func newConn(wsConn *websocket.Conn, userID string, hub *Hub) *Conn {
	return &Conn{
		ws:     wsConn,
		send:   make(chan []byte, 16),
		userID: userID,
		hub:    hub,
	}
}

func (c *Conn) Send(msg []byte) {
	select {
	case c.send <- msg:
	default:
		c.hub.remove(c)
	}
}

func (c *Conn) readPump() {
	defer c.hub.remove(c)
	defer c.ws.Close()

	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.hub.incoming <- InboundEvent{UserID: c.userID, Data: msg}
	}
}

func (c *Conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
