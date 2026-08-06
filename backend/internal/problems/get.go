package problems

import (
	"errors"
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"
	"github.com/anssuy/code-colosseum/backend/internal/testcases"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) Get(c *gin.Context) {
	slug := c.Param("slug")

	problem, err := h.queries.GetProblemBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "problem not found")
			return
		}
		httpx.InternalError(c, "get problem by slug error", err, "could not get problem")
		return
	}

	testCases, err := h.queries.ListSampleTestCasesForProblem(c.Request.Context(), problem.ID)
	if err != nil {
		httpx.InternalError(c, "get sample test cases error", err, "could not get test cases")
		return
	}

	items := make([]testcases.TestCaseResponse, len(testCases))
	for i, t := range testCases {
		items[i] = testcases.TestCaseResponseFrom(t, false) // will have a use later
	}

	c.JSON(http.StatusOK, gin.H{
		"problem":   ProblemResponseFrom(problem),
		"testCases": items,
	})
}
