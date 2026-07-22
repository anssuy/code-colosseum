package auth

import (
	"log"
	"net/http"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	queries *dbgen.Queries
	tokens  *TokenManager
}

func NewHandler(queries *dbgen.Queries, tokens *TokenManager) *Handler {
	return &Handler{
		queries: queries,
		tokens:  tokens,
	}
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func internalError(c *gin.Context, logMsg string, err error, userMsg string) {
	log.Printf("%s: %v", logMsg, err)
	writeError(c, http.StatusInternalServerError, userMsg)
}

func userJSON(user dbgen.User) gin.H {
	return gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"rating":     user.Rating,
		"wins":       user.Wins,
		"losses":     user.Losses,
		"created_at": user.CreatedAt,
	}
}
