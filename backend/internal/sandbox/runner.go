package sandbox

import (
	"context"
	"errors"
	"fmt"
)

var ErrTimeout = errors.New("execution timed out")

type runner func(ctx context.Context, code string, cfg ...runConfig) (string, error)

var runners = map[string]runner{
	"javascript": runWithNode,
	"typescript": runWithTypeScript,
	"python":     runWithPython,
}

func getRunConfig(cfg []runConfig) runConfig {
	if len(cfg) > 0 {
		return cfg[0]
	}
	return defaultRunConfig()
}

func Run(ctx context.Context, language, code, stdin string) (string, error) {
	run, ok := runners[language]
	if !ok {
		return "", fmt.Errorf("unsupported language")
	}

	cfg := defaultRunConfig()
	cfg.stdin = stdin

	return run(ctx, code, cfg)
}

func IsSupported(language string) bool {
	_, ok := runners[language]
	return ok
}
