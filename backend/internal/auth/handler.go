package auth

import (
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
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
