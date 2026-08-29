package problems

import (
	"encoding/json"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/stub"

	"github.com/jackc/pgx/v5/pgtype"
)

type ParamSpec struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type Stubs struct {
	Python     string `json:"python"`
	JavaScript string `json:"javascript"`
	TypeScript string `json:"typescript"`
}

type ProblemResponse struct {
	ID            pgtype.UUID        `json:"id"`
	Title         string             `json:"title"`
	Slug          string             `json:"slug"`
	Difficulty    dbgen.Difficulty   `json:"difficulty"`
	Description   string             `json:"description"`
	TimeLimitMs   int32              `json:"timeLimitMs"`
	MemoryLimitMb int32              `json:"memoryLimitMb"`
	FunctionName  string             `json:"functionName"`
	Params        []ParamSpec        `json:"params"`
	ReturnType    string             `json:"returnType"`
	Stubs         Stubs              `json:"stubs"`
	CreatedAt     pgtype.Timestamptz `json:"createdAt"`
}

func ProblemResponseFrom(problem dbgen.Problem) ProblemResponse {
	var params []ParamSpec
	_ = json.Unmarshal(problem.Params, &params)

	stubParams := make([]stub.Param, len(params))
	for i, p := range params {
		stubParams[i] = stub.Param{Name: p.Name, Type: p.Type}
	}

	sig := stub.Signature{
		FunctionName: problem.FunctionName,
		Params:       stubParams,
		ReturnType:   problem.ReturnType,
	}

	return ProblemResponse{
		ID:            problem.ID,
		Title:         problem.Title,
		Slug:          problem.Slug,
		Difficulty:    problem.Difficulty,
		Description:   problem.Description,
		TimeLimitMs:   problem.TimeLimitMs,
		MemoryLimitMb: problem.MemoryLimitMb,
		FunctionName:  problem.FunctionName,
		Params:        params,
		ReturnType:    problem.ReturnType,
		Stubs: Stubs{
			Python:     stub.Python(sig),
			JavaScript: stub.JavaScript(sig),
			TypeScript: stub.TypeScript(sig),
		},
		CreatedAt: problem.CreatedAt,
	}
}
