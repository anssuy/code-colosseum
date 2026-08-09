package tags

import (
	"errors"
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) Get(c *gin.Context) {
	slug := c.Param("slug")

	tag, err := h.queries.GetTagBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "tag not found")
			return
		}
		httpx.InternalError(c, "get tag by slug error", err, "could not get tag")
		return
	}

	c.JSON(http.StatusOK, gin.H{"tag": TagResponseFrom(tag)})
}
