package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid login data")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := h.queries.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusUnauthorized, "invalid email or password")
			return
		}
		httpx.InternalError(c, "get user by email error", err, "could not log in")
		return
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		httpx.WriteError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	accessToken, refreshToken, err := h.createSession(c.Request.Context(), user.ID)
	if err != nil {
		httpx.InternalError(c, "create login session error", err, "could not log in")
		return
	}

	setAuthCookies(c, accessToken, refreshToken)
	c.JSON(http.StatusOK, gin.H{"user": userJSON(user)})
}
