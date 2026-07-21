package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
)

var isProduction bool

func Init() {
	isProduction = os.Getenv("APP_ENV") == "prod"
}

func setAuthCookie(c *gin.Context, name, value string, maxAge int, path string) {
	c.SetCookie(name, value, maxAge, path, "", isProduction, true)
}

func setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetSameSite(http.SameSiteLaxMode)
	setAuthCookie(c, accessTokenCookieName, accessToken, int(accessTokenTTL.Seconds()), "/")
	setAuthCookie(c, refreshTokenCookieName, refreshToken, int(refreshTokenTTL.Seconds()), "/api/auth")
}

func clearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	setAuthCookie(c, accessTokenCookieName, "", -1, "/")
	setAuthCookie(c, refreshTokenCookieName, "", -1, "/api/auth")
}
