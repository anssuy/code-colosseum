package match

import (
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
)

type MatchResponse struct {
	ID          pgtype.UUID       `json:"id"`
	ProblemID   pgtype.UUID       `json:"problemId"`
	PlayerOneID pgtype.UUID       `json:"playerOneId"`
	PlayerTwoID pgtype.UUID       `json:"playerTwoId"`
	Status      dbgen.MatchStatus `json:"status"`
}

func MatchResponseFrom(match dbgen.Match) MatchResponse {
	return MatchResponse{
		ID:          match.ID,
		ProblemID:   match.ProblemID,
		PlayerOneID: match.PlayerOneID,
		PlayerTwoID: match.PlayerTwoID,
		Status:      match.Status,
	}
}
