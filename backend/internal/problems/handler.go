package problems

import dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"

type Handler struct {
	queries *dbgen.Queries
}

func NewHandler(queries *dbgen.Queries) *Handler {
	return &Handler{queries: queries}
}
