package problems

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type updateProblemRequest struct {
	Title        string           `json:"title" binding:"required"`
	Slug         string           `json:"slug" binding:"required"`
	Difficulty   dbgen.Difficulty `json:"difficulty" binding:"required,oneof=easy medium hard"`
	Description  string           `json:"description" binding:"required"`
	FunctionName string           `json:"functionName" binding:"required"`
	Params       []ParamSpec      `json:"params" binding:"required,dive"`
	ReturnType   string           `json:"returnType" binding:"required"`
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := httpx.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req updateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid problem data")
		return
	}

	title := strings.TrimSpace(req.Title)
	slug := strings.TrimSpace(req.Slug)
	functionName := strings.TrimSpace(req.FunctionName)
	if title == "" || slug == "" || functionName == "" {
		httpx.WriteError(c, http.StatusBadRequest, "title, slug, and function name cannot be empty")
		return
	}

	paramsJSON, err := json.Marshal(req.Params)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid params")
		return
	}

	problem, err := h.queries.UpdateProblem(
		c.Request.Context(),
		dbgen.UpdateProblemParams{
			ID:           id,
			Title:        title,
			Slug:         slug,
			Difficulty:   req.Difficulty,
			Description:  req.Description,
			FunctionName: functionName,
			Params:       paramsJSON,
			ReturnType:   req.ReturnType,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "problem not found")
			return
		}
		handleUpdateProblemError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"problem": ProblemResponseFrom(problem)})
}

func handleUpdateProblemError(c *gin.Context, err error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		switch pgErr.ConstraintName {
		case "problems_slug_key":
			httpx.WriteError(c, http.StatusConflict, "slug is already taken")
		default:
			httpx.WriteError(c, http.StatusConflict, "problem already exists")
		}
		return
	}
	httpx.InternalError(c, "update problem error", err, "could not update problem")
}
