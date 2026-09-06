package judge

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/anssuy/code-colosseum/backend/internal/language"
)

const (
	Accepted          = "accepted"
	WrongAnswer       = "wrong_answer"
	RuntimeError      = "runtime_error"
	TimeLimitExceeded = "time_limit_exceeded"
)

type TestCase struct {
	Input          string
	ExpectedOutput string
}

type Result struct {
	Status          string
	PassedTests     int32
	TotalTests      int32
	ExecutionTimeMS int64
}

func wrapCode(lang, functionName string, params []language.Param, returnType, code string) (string, error) {
	sig := language.Signature{
		FunctionName: functionName,
		Params:       params,
		ReturnType:   returnType,
	}

	wrapped, ok := language.Harness(lang, sig, code)
	if !ok {
		return "", errors.New("unsupported language")
	}
	return wrapped, nil
}

func Run(ctx context.Context, lang, functionName string, params []language.Param, returnType, code string, testCases []TestCase) Result {
	result := Result{
		Status:     Accepted,
		TotalTests: int32(len(testCases)),
	}

	if err := checkImports(ctx, lang, code); err != nil {
		result.Status = RuntimeError
		return result
	}

	wrapped, err := wrapCode(lang, functionName, params, returnType, code)
	if err != nil {
		result.Status = RuntimeError
		return result
	}

	start := time.Now()

	for _, tc := range testCases {
		output, err := language.Run(ctx, lang, wrapped, tc.Input)

		if err != nil {
			log.Printf("judge error: %v\noutput: %s", err, output)

			if errors.Is(err, language.ErrTimeout) {
				result.Status = TimeLimitExceeded
			} else {
				result.Status = RuntimeError
			}

			result.ExecutionTimeMS = time.Since(start).Milliseconds()
			return result
		}

		if !jsonEqual(output, tc.ExpectedOutput) {
			result.Status = WrongAnswer
			result.ExecutionTimeMS = time.Since(start).Milliseconds()
			return result
		}

		result.PassedTests++
	}

	result.ExecutionTimeMS = time.Since(start).Milliseconds()
	return result
}

func jsonEqual(a, b string) bool {
	var va, vb interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(a)), &va); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(b)), &vb); err != nil {
		return false
	}

	normA, _ := json.Marshal(va)
	normB, _ := json.Marshal(vb)
	return string(normA) == string(normB)
}
