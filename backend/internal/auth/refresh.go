package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) Refresh(c *gin.Context) {
	oldRefreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		httpx.WriteError(c, http.StatusUnauthorized, "refresh token is required")
		return
	}

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		httpx.InternalError(c, "generate refresh token error", err, "could not refresh session")
		return
	}

	userID, err := h.queries.RotateRefreshToken(
		c.Request.Context(),
		dbgen.RotateRefreshTokenParams{
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
			httpx.WriteError(c, http.StatusUnauthorized, "invalid or expired refresh token")
			return
		}

		httpx.InternalError(c, "rotate refresh token error", err, "could not refresh session")
		return
	}

	accessToken, err := h.tokens.GenerateAccessToken(userID.String())
	if err != nil {
		httpx.InternalError(c, "generate access token error", err, "could not refresh session")
		return
	}

	user, err := h.queries.GetUserByID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		log.Printf("get refreshed user error: %v", err)
		clearAuthCookies(c)
		httpx.WriteError(c, http.StatusUnauthorized, "user no longer exists")
		return
	}

	setAuthCookies(c, accessToken, newRefreshToken)

	c.JSON(http.StatusOK, gin.H{"user": userJSON(user)})
}
