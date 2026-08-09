package auth

import (
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err == nil {
		tokenHash := HashRefreshToken(refreshToken)
		if err := h.queries.DeleteRefreshTokenByHash(c.Request.Context(), tokenHash); err != nil {
			httpx.InternalError(c, "delete refresh token", err, "could not log out")
			return
		}
	}

	clearAuthCookies(c)
	c.Status(http.StatusNoContent)
}
