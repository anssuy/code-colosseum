package language

import (
	"context"
	"fmt"
	"strings"
)

func typescriptRun(ctx context.Context, code string, cfg runConfig) (string, error) {
	return runInContainer(ctx, nodeImage, "main.ts", code, []string{"tsx", "main.ts"}, cfg)
}

func typescriptType(t string) string {
	if strings.HasSuffix(t, "[]") {
		inner := typescriptType(strings.TrimSuffix(t, "[]"))
		return inner + "[]"
	}
	switch t {
	case "int", "float":
		return "number"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	default:
		return "any"
	}
}

func typescriptStub(sig Signature) string {
	var params []string
	for _, p := range sig.Params {
		params = append(params, fmt.Sprintf("%s: %s", p.Name, typescriptType(p.Type)))
	}

	returnType := typescriptType(sig.ReturnType)

	return fmt.Sprintf("function %s(%s): %s {\n\n}\n",
		sig.FunctionName, strings.Join(params, ", "), returnType)
}

func typescriptHarness(sig Signature, userCode string) string {
	return fmt.Sprintf(`%s

const args = JSON.parse(require('fs').readFileSync(0, 'utf8'));
const result = %s(...args);
console.log(JSON.stringify(result));
`, userCode, sig.FunctionName)
}
