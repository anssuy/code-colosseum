package problems

import (
	"net/http"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
)

type listProblemsQuery struct {
	httpx.PaginationQuery
	Difficulty string `form:"difficulty" binding:"omitempty,oneof=easy medium hard"`
}

func (h *Handler) List(c *gin.Context) {
	var query listProblemsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	var (
		problems []dbgen.Problem
		total    int64
		err      error
	)

	if query.Difficulty == "" {
		problems, err = h.queries.ListProblems(
			c.Request.Context(),
			dbgen.ListProblemsParams{
				Limit:  query.Limit,
				Offset: query.Offset,
			},
		)
		if err == nil {
			total, err = h.queries.CountProblems(c.Request.Context())
		}
	} else {
		diff := dbgen.Difficulty(query.Difficulty)

		problems, err = h.queries.ListProblemsByDifficulty(
			c.Request.Context(),
			dbgen.ListProblemsByDifficultyParams{
				Difficulty: diff,
				Limit:      query.Limit,
				Offset:     query.Offset,
			},
		)
		if err == nil {
			total, err = h.queries.CountProblemsByDifficulty(c.Request.Context(), diff)
		}
	}

	if err != nil {
		httpx.InternalError(c, "list problems error", err, "could not list problems")
		return
	}

	data := make([]ProblemResponse, len(problems))
	for i, p := range problems {
		data[i] = ProblemResponseFrom(p)
	}

	c.JSON(http.StatusOK, httpx.PaginationResponse[ProblemResponse]{
		Data:       data,
		TotalCount: total,
		Limit:      query.Limit,
		Offset:     query.Offset,
	})
}
