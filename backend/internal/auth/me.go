package auth

import (
	"errors"
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) Me(c *gin.Context) {
	userIDString, exists := getAuthenticatedUserID(c)
	if !exists {
		httpx.WriteError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(userIDString); err != nil || !userID.Valid {
		httpx.WriteError(c, http.StatusUnauthorized, "invalid access token")
		return
	}

	user, err := h.queries.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(c, http.StatusUnauthorized, "user no longer exists")
			return
		}
		httpx.InternalError(c, "get current user error", err, "could not retrieve user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": UserResponseFrom(user)})
}
