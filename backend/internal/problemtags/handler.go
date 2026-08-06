package problemtags

import (
	"net/http"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"
	"github.com/anssuy/code-colosseum/backend/internal/problems"
	"github.com/anssuy/code-colosseum/backend/internal/tags"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	queries *dbgen.Queries
}

func NewHandler(queries *dbgen.Queries) *Handler {
	return &Handler{queries: queries}
}

func (h *Handler) AddTag(c *gin.Context) {
	problemID, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}
	tagID, ok := httpx.ParseUUIDParam(c, "tagId")
	if !ok {
		return
	}

	if err := h.queries.AddTagToProblem(
		c.Request.Context(),
		dbgen.AddTagToProblemParams{
			ProblemID: problemID,
			TagID:     tagID,
		},
	); err != nil {
		httpx.InternalError(c, "add tag to problem error", err, "could not add tag")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) RemoveTag(c *gin.Context) {
	problemID, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}
	tagID, ok := httpx.ParseUUIDParam(c, "tagId")
	if !ok {
		return
	}

	if err := h.queries.RemoveTagFromProblem(
		c.Request.Context(),
		dbgen.RemoveTagFromProblemParams{
			ProblemID: problemID,
			TagID:     tagID,
		},
	); err != nil {
		httpx.InternalError(c, "remove tag from problem error", err, "could not remove tag")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListTags(c *gin.Context) {
	problemID, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	rows, err := h.queries.ListTagsForProblem(c.Request.Context(), problemID)
	if err != nil {
		httpx.InternalError(c, "list tags for problem error", err, "could not list tags")
		return
	}

	items := make([]tags.TagResponse, len(rows))
	for i, t := range rows {
		items[i] = tags.TagResponseFrom(t)
	}

	c.JSON(http.StatusOK, gin.H{"tags": items})
}

func (h *Handler) ListProblems(c *gin.Context) {
	tagID, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	rows, err := h.queries.ListProblemsForTag(c.Request.Context(), tagID)
	if err != nil {
		httpx.InternalError(c, "list problems for tag error", err, "could not list problems")
		return
	}

	items := make([]problems.ProblemResponse, len(rows))
	for i, p := range rows {
		items[i] = problems.ProblemResponseFrom(p)
	}

	c.JSON(http.StatusOK, gin.H{"problems": items})
}
