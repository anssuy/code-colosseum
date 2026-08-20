package sandbox

import "context"

func runWithPython(ctx context.Context, code string, cfg ...runConfig) (string, error) {
	return runInContainer(ctx, pythonImage, "main.py", code, []string{"python3", "main.py"}, getRunConfig(cfg))
}
