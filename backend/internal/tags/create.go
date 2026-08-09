package tags

import (
	"errors"
	"net/http"
	"strings"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolationCode = "23505"

type createTagRequest struct {
	Slug string `json:"slug" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createTagRequest
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

	tag, err := h.queries.CreateTag(
		c.Request.Context(),
		dbgen.CreateTagParams{
			Slug: slug,
			Name: name,
		},
	)
	if err != nil {
		handleTagUniqueViolation(c, err, "create tag error", "could not create tag")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tag": TagResponseFrom(tag)})
}

func handleTagUniqueViolation(c *gin.Context, err error, logMsg, fallbackMsg string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		switch pgErr.ConstraintName {
		case "tags_slug_key":
			httpx.WriteError(c, http.StatusConflict, "slug is already taken")
		case "tags_name_key":
			httpx.WriteError(c, http.StatusConflict, "name is already taken")
		default:
			httpx.WriteError(c, http.StatusConflict, "tag already exists")
		}
		return
	}
	httpx.InternalError(c, logMsg, err, fallbackMsg)
}
