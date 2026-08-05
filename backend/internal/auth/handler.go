package auth

import (
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

func userJSON(user dbgen.User) gin.H {
	return gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"rating":    user.Rating,
		"wins":      user.Wins,
		"losses":    user.Losses,
		"createdAt": user.CreatedAt,
	}
}
