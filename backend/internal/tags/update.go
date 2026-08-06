package tags

import (
	"errors"
	"net/http"
	"strings"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type updateTagRequest struct {
	Slug string `json:"slug" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req updateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid tag data")
		return
	}

	slug := strings.TrimSpace(req.Slug)
	name := strings.TrimSpace(req.Name)
	if slug == "" || name == "" {
		httpx.WriteError(c, http.StatusBadRequest, "slug and name cannot be empty")
		return
	}

	tag, err := h.queries.UpdateTag(
		c.Request.Context(),
		dbgen.UpdateTagParams{
			ID:   id,
			Slug: slug,
			Name: name,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "tag not found")
			return
		}
		handleTagUniqueViolation(c, err, "update tag error", "could not update tag")
		return
	}

	c.JSON(http.StatusOK, gin.H{"tag": TagResponseFrom(tag)})
}
