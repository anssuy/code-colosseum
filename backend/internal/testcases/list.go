package testcases

import (
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
)

func (h *Handler) List(c *gin.Context) {
	problemID, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	testCases, err := h.queries.ListSampleTestCasesForProblem(c.Request.Context(), problemID)
	if err != nil {
		httpx.InternalError(c, "list tags error", err, "could not list tags")
		return
	}

	items := make([]TestCaseResponse, len(testCases))
	for i, t := range testCases {
		items[i] = TestCaseResponseFrom(t, false)
	}

	c.JSON(http.StatusOK, gin.H{"testCases": items})

}
