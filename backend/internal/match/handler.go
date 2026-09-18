package match

import (
	"errors"
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/auth"
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"
	"github.com/anssuy/code-colosseum/backend/internal/problems"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type listMatchesQuery struct {
	httpx.PaginationQuery
}

type Handler struct {
	queries *dbgen.Queries
}

func NewHandler(queries *dbgen.Queries) *Handler {
	return &Handler{queries: queries}
}

func (h *Handler) Get(c *gin.Context) {
	matchID, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	userIDString, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(userIDString); err != nil {
		httpx.WriteError(c, http.StatusUnauthorized, "invalid user")
		return
	}

	ctx := c.Request.Context()

	dbMatch, err := h.queries.GetMatch(ctx, matchID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "match not found")
			return
		}
		httpx.InternalError(c, "get match error", err, "could not load match")
		return
	}

	if dbMatch.PlayerOneID != userID && dbMatch.PlayerTwoID != userID {
		httpx.WriteError(c, http.StatusForbidden, "you are not part of this match")
		return
	}

	problem, err := h.queries.GetProblemByID(ctx, dbMatch.ProblemID)
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

func (h *Handler) GetActive(c *gin.Context) {
	userIDString, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(userIDString); err != nil {
		httpx.WriteError(c, http.StatusUnauthorized, "invalid user")
		return
	}

	dbMatch, err := h.queries.GetActiveMatchForUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusOK, gin.H{"match": nil})
			return
		}
		httpx.InternalError(c, "get active match error", err, "could not check active match")
		return
	}

	c.JSON(http.StatusOK, gin.H{"match": MatchResponseFrom(dbMatch)})
}

func (h *Handler) List(c *gin.Context) {
	var query listMatchesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	userIDString, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(userIDString); err != nil {
		httpx.WriteError(c, http.StatusUnauthorized, "invalid user")
		return
	}

	var (
		matches []dbgen.Match
		total   int64
		err     error
	)

	matches, err = h.queries.ListMatchesForUser(
		c.Request.Context(),
		dbgen.ListMatchesForUserParams{
			PlayerOneID: userID,
			Limit:       query.Limit,
			Offset:      query.Offset,
		},
	)
	if err != nil {
		httpx.InternalError(c, "list matches error", err, "could not list matches")
		return
	}

	total, err = h.queries.CountMatchesForUser(c.Request.Context(), userID)

	data := make([]MatchResponse, len(matches))
	for i, m := range matches {
		data[i] = MatchResponseFrom(m)
	}

	c.JSON(http.StatusOK, httpx.PaginationResponse[MatchResponse]{
		Data:       data,
		TotalCount: total,
		Limit:      query.Limit,
		Offset:     query.Offset,
	})
}
