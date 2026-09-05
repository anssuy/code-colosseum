package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu            sync.RWMutex
	conns         map[string]*Conn
	incoming      chan InboundEvent
	handlers      []Handler
	connectFns    []func(userID string)
	disconnectFns []func(userID string)
}

func NewHub() *Hub {
	h := &Hub{
		conns:    make(map[string]*Conn),
		incoming: make(chan InboundEvent, 64),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for event := range h.incoming {
		for _, handler := range h.handlers {
			handler(event.UserID, event.Data)
		}
	}
}

func (h *Hub) remove(c *Conn) {
	h.mu.Lock()
	existed := false
	if existing, ok := h.conns[c.userID]; ok && existing == c {
		delete(h.conns, c.userID)
		existed = true
	}
	h.mu.Unlock()

	if existed {
		for _, fn := range h.disconnectFns {
			fn(c.userID)
		}
	}
}

func (h *Hub) OnConnect(fn func(userID string)) {
	h.connectFns = append(h.connectFns, fn)
}

func (h *Hub) OnDisconnect(fn func(userID string)) {
	h.disconnectFns = append(h.disconnectFns, fn)
}

func (h *Hub) OnMessage(fn Handler) {
	h.handlers = append(h.handlers, fn)
}

func (h *Hub) Register(userID string, wsConn *websocket.Conn) *Conn {
	c := newConn(wsConn, userID, h)

	h.mu.Lock()
	if old, ok := h.conns[userID]; ok {
		old.Send([]byte(`{"type":"replaced"}`))
		close(old.send)
	}
	h.conns[userID] = c
	h.mu.Unlock()

	go c.writePump()
	go c.readPump()

	for _, fn := range h.connectFns {
		fn(userID)
	}

	return c
}

func (h *Hub) SendTo(userID string, msg []byte) {
	h.mu.RLock()
	c, ok := h.conns[userID]
	h.mu.RUnlock()

	if ok {
		c.Send(msg)
	}
}
