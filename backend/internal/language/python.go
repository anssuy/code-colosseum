package language

import (
	"context"
	"fmt"
	"strings"
)

const (
	pythonPrelude = "from typing import List, Dict, Optional, Tuple, Set\nfrom collections import *\n\n"
	pythonImage   = "sandbox-python:latest"
)

func pythonRun(ctx context.Context, code string, cfg runConfig) (string, error) {
	return runInContainer(ctx, pythonImage, "main.py", code, []string{"python3", "main.py"}, cfg)
}

func pythonType(t string) string {
	if strings.HasSuffix(t, "[]") {
		inner := pythonType(strings.TrimSuffix(t, "[]"))
		return "List[" + inner + "]"
	}
	switch t {
	case "int":
		return "int"
	case "float":
		return "float"
	case "string":
		return "str"
	case "bool":
		return "bool"
	default:
		return "Any"
	}
}

func pythonStub(sig Signature) string {
	var params []string
	for _, p := range sig.Params {
		params = append(params, fmt.Sprintf("%s: %s", p.Name, pythonType(p.Type)))
	}

	returnType := pythonType(sig.ReturnType)

	return fmt.Sprintf("def %s(%s) -> %s:\n    pass\n",
		sig.FunctionName, strings.Join(params, ", "), returnType)
}

func pythonHarness(sig Signature, userCode string) string {
	return fmt.Sprintf(`%simport json
import sys

%s

if __name__ == "__main__":
    args = json.loads(sys.stdin.read())
    result = %s(*args)
    print(json.dumps(result))
`, pythonPrelude, userCode, sig.FunctionName)
}
