package tags

import (
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Delete(c *gin.Context) {
	id, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	rows, err := h.queries.DeleteTag(c.Request.Context(), id)
	if err != nil {
		httpx.InternalError(c, "delete tag error", err, "could not delete tag")
		return
	}
	if rows == 0 {
		httpx.WriteError(c, http.StatusNotFound, "tag not found")
		return
	}

	c.Status(http.StatusNoContent)
}
