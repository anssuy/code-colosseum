package tags

import (
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
)

type TagResponse struct {
	ID   pgtype.UUID `json:"id"`
	Slug string      `json:"slug"`
	Name string      `json:"name"`
}

func TagResponseFrom(tag dbgen.Tag) TagResponse {
	return TagResponse{
		ID:   tag.ID,
		Slug: tag.Slug,
		Name: tag.Name,
	}
}
