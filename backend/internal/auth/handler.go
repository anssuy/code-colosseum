package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	db "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/jackc/pgx/v5"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolationCode = "23505"

type Handler struct {
	queries *db.Queries
	tokens  *TokenManager
}

func NewHandler(
	queries *db.Queries,
	tokens *TokenManager,
) *Handler {
	return &Handler{
		queries: queries,
		tokens:  tokens,
	}
}

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
		writeError(
			c,
			http.StatusBadRequest,
			"username must be between 3 and 30 characters",
		)
		return
	}

	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.queries.CreateUser(
		c.Request.Context(),
		db.CreateUserParams{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		},
	)
	if err != nil {
		handleCreateUserError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": userJSON(user),
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid login data")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := h.queries.GetUserByEmail(
		c.Request.Context(),
		email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(
				c,
				http.StatusUnauthorized,
				"invalid email or password",
			)
			return
		}

		log.Printf("get user by email error: %v", err)
		writeError(c, http.StatusInternalServerError, "could not log in")
		return
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		writeError(
			c,
			http.StatusUnauthorized,
			"invalid email or password",
		)
		return
	}

	accessToken, err := h.tokens.GenerateAccessToken(user.ID.String())
	if err != nil {
		log.Printf("generate access token error: %v", err)
		writeError(c, http.StatusInternalServerError, "could not log in")
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		"access_token",
		accessToken,
		15*60,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"user": userJSON(user),
	})
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

	log.Printf("create user error: %v", err)
	writeError(c, http.StatusInternalServerError, "could not create user")
}

func userJSON(user db.User) gin.H {
	return gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"rating":     user.Rating,
		"wins":       user.Wins,
		"losses":     user.Losses,
		"created_at": user.CreatedAt,
	}
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": message,
	})
}
