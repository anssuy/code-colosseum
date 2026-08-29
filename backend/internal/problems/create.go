package problems

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolationCode = "23505"

var slugRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type createProblemRequest struct {
	Title        string           `json:"title" binding:"required"`
	Slug         string           `json:"slug" binding:"required"`
	Difficulty   dbgen.Difficulty `json:"difficulty" binding:"required,oneof=easy medium hard"`
	Description  string           `json:"description" binding:"required"`
	FunctionName string           `json:"functionName" binding:"required"`
	Params       []ParamSpec      `json:"params" binding:"required,dive"`
	ReturnType   string           `json:"returnType" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createProblemRequest
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

	if !slugRe.MatchString(slug) {
		httpx.WriteError(c, http.StatusBadRequest, "invalid slug format")
		return
	}

	paramsJSON, err := json.Marshal(req.Params)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid params")
		return
	}

	problem, err := h.queries.CreateProblem(
		c.Request.Context(),
		dbgen.CreateProblemParams{
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
		handleCreateProblemError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"problem": ProblemResponseFrom(problem)})
}

func handleCreateProblemError(c *gin.Context, err error) {
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
	httpx.InternalError(c, "create problem error", err, "could not create problem")
}
