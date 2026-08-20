package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const userIDContextKey = "authenticatedUserID"

func Middleware(tokens *TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(accessTokenCookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		userID, err := tokens.ValidateAccessToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired access token"})
			return
		}

		c.Set(userIDContextKey, userID)
		c.Next()
	}
}

func GetAuthenticatedUserID(c *gin.Context) (string, bool) {
	value, exists := c.Get(userIDContextKey)
	if !exists {
		return "", false
	}

	userID, ok := value.(string)
	return userID, ok
}
