package leaderboard

import (
	"net/http"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	queries *dbgen.Queries
}

func NewHandler(queries *dbgen.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	var query httpx.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination"})
		return
	}

	users, err := h.queries.ListLeaderboard(c, dbgen.ListLeaderboardParams{
		Limit:  query.Limit,
		Offset: query.Offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get leaderboard"})
		return
	}

	entries := make([]EntryResponse, 0, len(users))
	for i, user := range users {
		entries = append(entries, EntryResponseFrom(user, int(query.Offset)+i+1))
	}

	c.JSON(http.StatusOK, entries)
}
