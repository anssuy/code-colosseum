package judge

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/anssuy/code-colosseum/backend/internal/sandbox"
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

func Run(ctx context.Context, language, code string, testCases []TestCase) Result {
	result := Result{
		Status:     Accepted,
		TotalTests: int32(len(testCases)),
	}

	start := time.Now()

	for _, tc := range testCases {
		output, err := sandbox.Run(ctx, language, code, tc.Input)

		if err != nil {
			log.Printf("judge error: %v\noutput: %s", err, output)

			if errors.Is(err, sandbox.ErrTimeout) {
				result.Status = TimeLimitExceeded
			} else {
				result.Status = RuntimeError
			}

			result.ExecutionTimeMS = time.Since(start).Milliseconds()
			return result
		}

		if normalizeOutput(output) != normalizeOutput(tc.ExpectedOutput) {
			result.Status = WrongAnswer
			result.ExecutionTimeMS = time.Since(start).Milliseconds()
			return result
		}

		result.PassedTests++
	}

	result.ExecutionTimeMS = time.Since(start).Milliseconds()
	return result
}

func normalizeOutput(output string) string {
	return strings.TrimSpace(strings.ReplaceAll(output, "\r\n", "\n"))
}
