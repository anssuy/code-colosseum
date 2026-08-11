package submissions

import (
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type SubmissionResponse struct {
	ID              pgtype.UUID        `json:"id"`
	ProblemID       pgtype.UUID        `json:"problemId"`
	Language        string             `json:"language"`
	Status          string             `json:"status"`
	PassedTests     int32              `json:"passedTests"`
	TotalTests      int32              `json:"totalTests"`
	ExecutionTimeMS pgtype.Int8        `json:"executionTimeMs"`
	CreatedAt       pgtype.Timestamptz `json:"createdAt"`
}

func SubmissionResponseFrom(s dbgen.Submission) SubmissionResponse {
	return SubmissionResponse{
		ID:              s.ID,
		ProblemID:       s.ProblemID,
		Language:        s.Language,
		Status:          s.Status,
		PassedTests:     s.PassedTests,
		TotalTests:      s.TotalTests,
		ExecutionTimeMS: s.ExecutionTimeMs,
		CreatedAt:       s.CreatedAt,
	}
}
