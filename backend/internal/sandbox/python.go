package sandbox

import "context"

func runWithPython(ctx context.Context, code string) (string, error) {
	return runInContainer(ctx, pythonImage, "main.py", code, []string{"python3", "main.py"})
}
