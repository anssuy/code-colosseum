package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	db "github.com/anssuy/code-colosseum/backend/internal/db/generated"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) createSession(
	ctx context.Context,
	userID pgtype.UUID,
) (string, string, error) {
	accessToken, err := h.tokens.GenerateAccessToken(userID.String())
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = h.queries.CreateRefreshToken(
		ctx,
		db.CreateRefreshTokenParams{
			UserID:    userID,
			TokenHash: HashRefreshToken(refreshToken),
			ExpiresAt: pgtype.Timestamptz{
				Time:  time.Now().Add(refreshTokenTTL),
				Valid: true,
			},
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func setAuthCookies(
	c *gin.Context,
	accessToken string,
	refreshToken string,
) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		accessTokenCookieName,
		accessToken,
		int(accessTokenTTL.Seconds()),
		"/",
		"",
		false, // Change to true in production with HTTPS.
		true,
	)

	c.SetCookie(
		refreshTokenCookieName,
		refreshToken,
		int(refreshTokenTTL.Seconds()),
		"/api/auth",
		"",
		false, // Change to true in production with HTTPS.
		true,
	)
}

func clearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		accessTokenCookieName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		refreshTokenCookieName,
		"",
		-1,
		"/api/auth",
		"",
		false,
		true,
	)
}
