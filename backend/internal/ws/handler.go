package ws

import (
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/auth"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsHandler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *WsHandler {
	return &WsHandler{hub: hub}
}

func (h *WsHandler) Connect(c *gin.Context) {
	userID, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	h.hub.Register(userID, wsConn)
}
