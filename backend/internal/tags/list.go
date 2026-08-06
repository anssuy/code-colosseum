package tags

import (
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
)

func (h *Handler) List(c *gin.Context) {
	tags, err := h.queries.ListTags(c.Request.Context())
	if err != nil {
		httpx.InternalError(c, "list tags error", err, "could not list tags")
		return
	}

	items := make([]TagResponse, len(tags))
	for i, t := range tags {
		items[i] = TagResponseFrom(t)
	}

	c.JSON(http.StatusOK, gin.H{"tags": items})
}
