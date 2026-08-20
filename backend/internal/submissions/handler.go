package submissions

import (
	"errors"
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/auth"
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"
	"github.com/anssuy/code-colosseum/backend/internal/judge"
	"github.com/anssuy/code-colosseum/backend/internal/sandbox"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	queries *dbgen.Queries
}

func NewHandler(queries *dbgen.Queries) *Handler {
	return &Handler{queries: queries}
}

type submitRequest struct {
	Language   string `json:"language" binding:"required"`
	SourceCode string `json:"sourceCode" binding:"required"`
}

func (h *Handler) Submit(c *gin.Context) {
	var req submitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "language and source code are required")
		return
	}

	if !sandbox.IsSupported(req.Language) {
		httpx.WriteError(c, http.StatusBadRequest, "unsupported language")
		return
	}

	userIDString, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var userID, problemID pgtype.UUID

	if err := userID.Scan(userIDString); err != nil {
		httpx.WriteError(c, http.StatusUnauthorized, "invalid user")
		return
	}

	if err := problemID.Scan(c.Param("id")); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid problem ID")
		return
	}

	ctx := c.Request.Context()

	if _, err := h.queries.GetProblemByID(ctx, problemID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "problem not found")
			return
		}

		httpx.WriteError(c, http.StatusInternalServerError, "could not load problem")
		return
	}

	dbTestCases, err := h.queries.ListTestCasesForProblem(ctx, problemID)
	if err != nil {
		httpx.WriteError(c, http.StatusInternalServerError, "could not load test cases")
		return
	}

	if len(dbTestCases) == 0 {
		httpx.WriteError(c, http.StatusInternalServerError, "problem has no test cases")
		return
	}

	testCases := make([]judge.TestCase, len(dbTestCases))
	for i, tc := range dbTestCases {
		testCases[i] = judge.TestCase{
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
		}
	}

	result := judge.Run(ctx, req.Language, req.SourceCode, testCases)

	submission, err := h.queries.CreateSubmission(ctx, dbgen.CreateSubmissionParams{
		UserID:      userID,
		ProblemID:   problemID,
		Language:    req.Language,
		SourceCode:  req.SourceCode,
		Status:      result.Status,
		PassedTests: result.PassedTests,
		TotalTests:  result.TotalTests,
		ExecutionTimeMs: pgtype.Int8{
			Int64: result.ExecutionTimeMS,
			Valid: true,
		},
	})
	if err != nil {
		httpx.WriteError(c, http.StatusInternalServerError, "could not save submission")
		return
	}

	c.JSON(http.StatusCreated, SubmissionResponseFrom(submission))
}
