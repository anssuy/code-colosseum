package testcases

import (
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type TestCaseResponse struct {
	ID             pgtype.UUID `json:"id"`
	ProblemID      pgtype.UUID `json:"problemId"`
	Input          string      `json:"input"`
	ExpectedOutput string      `json:"expectedOutput,omitempty"`
	IsSample       bool        `json:"isSample"`
	Ord            int32       `json:"ord"`
}

func TestCaseResponseFrom(tc dbgen.ProblemTestCase, includeHidden bool) TestCaseResponse {
	resp := TestCaseResponse{
		ID:        tc.ID,
		ProblemID: tc.ProblemID,
		Input:     tc.Input,
		IsSample:  tc.IsSample,
		Ord:       tc.Ord,
	}
	if tc.IsSample || includeHidden {
		resp.ExpectedOutput = tc.ExpectedOutput
	}
	return resp
}
