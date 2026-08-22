package match

import (
	"errors"
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/auth"
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	queries *dbgen.Queries
	manager *Manager
}

func NewWSHandler(queries *dbgen.Queries, manager *Manager) *WSHandler {
	return &WSHandler{queries: queries, manager: manager}
}

func (h *WSHandler) Connect(c *gin.Context) {
	matchID, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	userIDString, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	dbMatch, err := h.queries.GetMatch(c.Request.Context(), matchID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "match not found")
			return
		}
		httpx.InternalError(c, "get match error", err, "could not load match")
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(userIDString); err != nil {
		httpx.WriteError(c, http.StatusUnauthorized, "invalid user")
		return
	}

	if userID != dbMatch.PlayerOneID && userID != dbMatch.PlayerTwoID {
		httpx.WriteError(c, http.StatusForbidden, "not a participant in this match")
		return
	}

	room, ok := h.manager.GetRoom(matchID)
	if !ok {
		room, err = h.manager.CreateRoom(c.Request.Context(), dbMatch)
		if err != nil {
			httpx.InternalError(c, "create room error", err, "could not start match room")
			return
		}
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	conn := NewConn(ws, room, userIDString)
	room.AddPlayer(userIDString, conn)

	go conn.WritePump()
	go conn.ReadPump()
}
