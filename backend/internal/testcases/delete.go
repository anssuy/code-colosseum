package testcases

import (
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Delete(c *gin.Context) {
	_, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}
	testCaseID, ok := httpx.ParseUUIDParam(c, "testCaseId")
	if !ok {
		return
	}

	if err := h.queries.DeleteTestCase(c.Request.Context(), testCaseID); err != nil {
		httpx.InternalError(c, "delete test case error", err, "could not delete test case")
		return
	}

	c.Status(http.StatusNoContent)
}
