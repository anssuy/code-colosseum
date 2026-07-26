package sandbox

import "context"

func runWithNode(ctx context.Context, code string) (string, error) {
	return runInContainer(ctx, nodeImage, "main.js", code, []string{"node", "main.js"})
}

func runWithTypeScript(ctx context.Context, code string) (string, error) {
	return runInContainer(ctx, nodeImage, "main.ts", code, []string{"tsx", "main.ts"})
}
