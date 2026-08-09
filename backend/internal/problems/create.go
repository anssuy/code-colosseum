package problems

import (
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

type createProblemRequest struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Difficulty  string `json:"difficulty" binding:"required,oneof=easy medium hard"`
	Description string `json:"description" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid problem data")
		return
	}

	title := strings.TrimSpace(req.Title)
	slug := strings.TrimSpace(req.Slug)
	if title == "" || slug == "" {
		httpx.WriteError(c, http.StatusBadRequest, "title and slug cannot be empty")
		return
	}

	var slugRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

	if !slugRe.MatchString(slug) {
		httpx.WriteError(c, http.StatusBadRequest, "invalid slug format")
		return
	}

	problem, err := h.queries.CreateProblem(
		c.Request.Context(),
		dbgen.CreateProblemParams{
			Title:       title,
			Slug:        slug,
			Difficulty:  dbgen.Difficulty(req.Difficulty),
			Description: req.Description,
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
