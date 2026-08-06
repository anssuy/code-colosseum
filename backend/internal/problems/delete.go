package problems

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

	rows, err := h.queries.DeleteProblem(c.Request.Context(), id)
	if err != nil {
		httpx.InternalError(c, "delete problem error", err, "could not delete problem")
		return
	}
	if rows == 0 {
		httpx.WriteError(c, http.StatusNotFound, "problem not found")
		return
	}

	c.Status(http.StatusNoContent)
}
