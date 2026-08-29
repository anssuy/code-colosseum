package harness

import "fmt"

func TypeScript(sig Signature, userCode string) string {
	return fmt.Sprintf(`%s

const args = JSON.parse(require('fs').readFileSync(0, 'utf8'));
const result = %s(...args);
console.log(JSON.stringify(result));
`, userCode, sig.FunctionName)
}
