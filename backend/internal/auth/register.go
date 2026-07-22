package auth

import (
	"errors"
	"net/http"
	"strings"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolationCode = "23505"

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid registration data")
		return
	}

	username := strings.TrimSpace(req.Username)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if len(username) < 3 || len(username) > 30 {
		writeError(c, http.StatusBadRequest, "username must be between 3 and 30 characters")
		return
	}

	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.queries.CreateUser(
		c.Request.Context(),
		dbgen.CreateUserParams{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		},
	)
	if err != nil {
		handleCreateUserError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": userJSON(user)})
}

func handleCreateUserError(c *gin.Context, err error) {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		switch pgErr.ConstraintName {
		case "users_username_key":
			writeError(c, http.StatusConflict, "username is already taken")

		case "users_email_key":
			writeError(c, http.StatusConflict, "email is already registered")

		default:
			writeError(c, http.StatusConflict, "user already exists")
		}

		return
	}

	internalError(c, "create user error", err, "could not create user")
}
