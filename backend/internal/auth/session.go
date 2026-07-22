package auth

import (
	"context"
	"fmt"
	"time"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"

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
		dbgen.CreateRefreshTokenParams{
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
