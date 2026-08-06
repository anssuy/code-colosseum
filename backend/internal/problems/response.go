package problems

import (
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
)

type ProblemResponse struct {
	ID            pgtype.UUID        `json:"id"`
	Title         string             `json:"title"`
	Slug          string             `json:"slug"`
	Difficulty    dbgen.Difficulty   `json:"difficulty"`
	Description   string             `json:"description"`
	TimeLimitMs   int32              `json:"timeLimitMs"`
	MemoryLimitMb int32              `json:"memoryLimitMb"`
	CreatedAt     pgtype.Timestamptz `json:"createdAt"`
}

func ProblemResponseFrom(problem dbgen.Problem) ProblemResponse {
	return ProblemResponse{
		ID:            problem.ID,
		Title:         problem.Title,
		Slug:          problem.Slug,
		Difficulty:    problem.Difficulty,
		Description:   problem.Description,
		TimeLimitMs:   problem.TimeLimitMs,
		MemoryLimitMb: problem.MemoryLimitMb,
		CreatedAt:     problem.CreatedAt,
	}
}
