package match

import (
	"errors"
	"net/http"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"
	"github.com/anssuy/code-colosseum/backend/internal/problems"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	queries *dbgen.Queries
}

type MatchResponse struct {
	ID          pgtype.UUID       `json:"id"`
	ProblemID   pgtype.UUID       `json:"problemId"`
	PlayerOneID pgtype.UUID       `json:"playerOneId"`
	PlayerTwoID pgtype.UUID       `json:"playerTwoId"`
	Status      dbgen.MatchStatus `json:"status"`
}

func NewHandler(queries *dbgen.Queries) *Handler {
	return &Handler{queries: queries}
}

func (h *Handler) Get(c *gin.Context) {
	matchID, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	dbMatch, err := h.queries.GetMatch(c.Request.Context(), matchID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "match not found")
			return
		}
		httpx.InternalError(c, "get match error", err, "could not load match")
		return
	}

	problem, err := h.queries.GetProblemByID(c.Request.Context(), dbMatch.ProblemID)
	if err != nil {
		httpx.InternalError(c, "get problem for match error", err, "could not load problem")
		return
	}

	matchResponse := MatchResponse{
		ID:          dbMatch.ID,
		ProblemID:   dbMatch.ProblemID,
		PlayerOneID: dbMatch.PlayerOneID,
		PlayerTwoID: dbMatch.PlayerTwoID,
		Status:      dbMatch.Status,
	}

	c.JSON(http.StatusOK, gin.H{
		"match":   matchResponse,
		"problem": problems.ProblemResponseFrom(problem),
	})
}
