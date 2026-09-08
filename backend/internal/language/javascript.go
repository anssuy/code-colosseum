package language

import (
	"context"
	"fmt"
	"strings"
)

const nodeImage = "sandbox-node:latest"

func javascriptRun(ctx context.Context, code string, cfg runConfig) (string, error) {
	return runInContainer(ctx, nodeImage, "main.js", code, []string{"node", "main.js"}, cfg)
}

func javascriptStub(sig Signature) string {
	var params []string
	for _, p := range sig.Params {
		params = append(params, p.Name)
	}

	return fmt.Sprintf("function %s(%s) {\n\n}\n",
		sig.FunctionName, strings.Join(params, ", "))
}

func javascriptHarness(sig Signature, userCode string) string {
	return fmt.Sprintf(`%s

const args = JSON.parse(require('fs').readFileSync(0, 'utf8'));
const result = %s(...args);
console.log(JSON.stringify(result));
`, userCode, sig.FunctionName)
}
