package judge

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

const pythonImportCheckScript = `
import ast
import sys

try:
    tree = ast.parse(sys.stdin.read())
except SyntaxError:
    print("SYNTAX_ERROR")
    sys.exit(1)

for node in ast.walk(tree):
    if isinstance(node, (ast.Import, ast.ImportFrom)):
        print("IMPORT_DETECTED")
        sys.exit(1)

print("OK")
`

func checkPythonImports(ctx context.Context, code string) error {
	cmd := exec.CommandContext(ctx, "python3", "-c", pythonImportCheckScript)
	cmd.Stdin = strings.NewReader(code)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		out := strings.TrimSpace(out.String())
		switch out {
		case "IMPORT_DETECTED":
			return fmt.Errorf("imports are not allowed")
		case "SYNTAX_ERROR":
			return fmt.Errorf("syntax error in submitted code")
		default:
			return fmt.Errorf("could not validate code")
		}
	}

	return nil
}

func checkImports(ctx context.Context, language, code string) error {
	switch language {
	case "python":
		return checkPythonImports(ctx, code)
	default:
		return nil
	}
}
