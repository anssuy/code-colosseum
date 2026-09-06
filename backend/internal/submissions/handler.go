package submissions

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/auth"
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"
	"github.com/anssuy/code-colosseum/backend/internal/judge"
	"github.com/anssuy/code-colosseum/backend/internal/language"

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
	MatchID    string `json:"matchId"`
}

func (h *Handler) Submit(c *gin.Context) {
	var req submitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "language and source code are required")
		return
	}

	if !language.IsValid(req.Language) {
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

	var matchID pgtype.UUID
	if req.MatchID != "" {
		if err := matchID.Scan(req.MatchID); err != nil {
			httpx.WriteError(c, http.StatusBadRequest, "invalid match ID")
			return
		}
	}

	ctx := c.Request.Context()

	problem, err := h.queries.GetProblemByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "problem not found")
			return
		}
		httpx.InternalError(c, "get problem by id error", err, "could not load problem")
		return
	}

	var rawParams []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	_ = json.Unmarshal(problem.Params, &rawParams)

	params := make([]language.Param, len(rawParams))
	for i, p := range rawParams {
		params[i] = language.Param{Name: p.Name, Type: p.Type}
	}

	dbTestCases, err := h.queries.ListTestCasesForProblem(ctx, problemID)
	if err != nil {
		httpx.InternalError(c, "list test cases error", err, "could not load test cases")
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

	result := judge.Run(ctx, req.Language, problem.FunctionName, params, problem.ReturnType, req.SourceCode, testCases)

	submission, err := h.queries.CreateSubmission(ctx, dbgen.CreateSubmissionParams{
		UserID:      userID,
		ProblemID:   problemID,
		MatchID:     matchID,
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
		httpx.InternalError(c, "create submission error", err, "could not save submission")
		return
	}

	c.JSON(http.StatusCreated, SubmissionResponseFrom(submission))
}
