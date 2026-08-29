package harness

import "fmt"

func Python(sig Signature, userCode string) string {
	return fmt.Sprintf(`import json
import sys

%s

if __name__ == "__main__":
    args = json.loads(sys.stdin.read())
    result = %s(*args)
    print(json.dumps(result))
`, userCode, sig.FunctionName)
}
