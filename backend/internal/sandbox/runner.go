package sandbox

import "context"

type runner func(ctx context.Context, code string) (string, error)

var runners = map[string]runner{
	"javascript": runWithNode,
	"typescript": runWithTypeScript,
	"python":     runWithPython,
}
