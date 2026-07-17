package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	db "github.com/anssuy/code-colosseum/backend/internal/db/generated"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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

	accessToken, refreshToken, err := h.createSession(
		c.Request.Context(),
		user.ID,
	)
	if err != nil {
		log.Printf("create login session error: %v", err)
		writeError(c, http.StatusInternalServerError, "could not log in")
		return
	}

	setAuthCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"user": userJSON(user),
	})
}

func (h *Handler) Me(c *gin.Context) {
	userIDString, exists := getAuthenticatedUserID(c)
	if !exists {
		writeError(
			c,
			http.StatusUnauthorized,
			"authentication required",
		)
		return
	}

	var userID pgtype.UUID

	if err := userID.Scan(userIDString); err != nil || !userID.Valid {
		writeError(
			c,
			http.StatusUnauthorized,
			"invalid access token",
		)
		return
	}

	user, err := h.queries.GetUserByID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(
				c,
				http.StatusUnauthorized,
				"user no longer exists",
			)
			return
		}

		log.Printf("get current user error: %v", err)
		writeError(
			c,
			http.StatusInternalServerError,
			"could not retrieve user",
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userJSON(user),
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	oldRefreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		writeError(
			c,
			http.StatusUnauthorized,
			"refresh token is required",
		)
		return
	}

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		log.Printf("generate refresh token error: %v", err)
		writeError(c, http.StatusInternalServerError, "could not refresh session")
		return
	}

	userID, err := h.queries.RotateRefreshToken(
		c.Request.Context(),
		db.RotateRefreshTokenParams{
			TokenHash:   HashRefreshToken(oldRefreshToken),
			TokenHash_2: HashRefreshToken(newRefreshToken),
			ExpiresAt: pgtype.Timestamptz{
				Time:  time.Now().Add(refreshTokenTTL),
				Valid: true,
			},
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			clearAuthCookies(c)
			writeError(
				c,
				http.StatusUnauthorized,
				"invalid or expired refresh token",
			)
			return
		}

		log.Printf("rotate refresh token error: %v", err)
		writeError(c, http.StatusInternalServerError, "could not refresh session")
		return
	}

	accessToken, err := h.tokens.GenerateAccessToken(userID.String())
	if err != nil {
		log.Printf("generate access token error: %v", err)
		writeError(c, http.StatusInternalServerError, "could not refresh session")
		return
	}

	user, err := h.queries.GetUserByID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		log.Printf("get refreshed user error: %v", err)
		clearAuthCookies(c)
		writeError(c, http.StatusUnauthorized, "user no longer exists")
		return
	}

	setAuthCookies(c, accessToken, newRefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"user": userJSON(user),
	})
}

func (h *Handler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookieName)

	if err == nil {
		tokenHash := HashRefreshToken(refreshToken)

		if err := h.queries.DeleteRefreshTokenByHash(
			c.Request.Context(),
			tokenHash,
		); err != nil {
			log.Printf("delete refresh token error: %v", err)
			writeError(c, http.StatusInternalServerError, "could not log out")
			return
		}
	}

	clearAuthCookies(c)
	c.Status(http.StatusNoContent)
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
