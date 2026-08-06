package testcases

import (
	"net/http"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
)

type createTestCaseRequest struct {
	Input          string `json:"input" binding:"required"`
	ExpectedOutput string `json:"expectedOutput" binding:"required"`
	IsSample       bool   `json:"isSample"`
	Ord            int32  `json:"ord"`
}

func (h *Handler) Create(c *gin.Context) {
	problemID, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req createTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid test case data")
		return
	}

	tc, err := h.queries.CreateProblemTestCase(
		c.Request.Context(),
		dbgen.CreateProblemTestCaseParams{
			ProblemID:      problemID,
			Input:          req.Input,
			ExpectedOutput: req.ExpectedOutput,
			IsSample:       req.IsSample,
			Ord:            req.Ord,
		},
	)
	if err != nil {
		httpx.InternalError(c, "create test case error", err, "could not create test case")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"testCase": TestCaseResponseFrom(tc, true)})
}
