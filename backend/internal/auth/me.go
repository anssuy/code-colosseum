package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) Me(c *gin.Context) {
	userIDString, exists := getAuthenticatedUserID(c)
	if !exists {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(userIDString); err != nil || !userID.Valid {
		writeError(c, http.StatusUnauthorized, "invalid access token")
		return
	}

	user, err := h.queries.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(c, http.StatusUnauthorized, "user no longer exists")
			return
		}
		internalError(c, "get current user error", err, "could not retrieve user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userJSON(user)})
}
