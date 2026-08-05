package httpx

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func InternalError(c *gin.Context, logMsg string, err error, userMsg string) {
	log.Printf("%s: %v", logMsg, err)
	WriteError(c, http.StatusInternalServerError, userMsg)
}
