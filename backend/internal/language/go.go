package language

import (
	"context"
	"fmt"
	"strings"
)

const goImage = "sandbox-go:latest"

func goType(t string) string {
	if before, ok := strings.CutSuffix(t, "[]"); ok {
		return "[]" + goType(before)
	}
	switch t {
	case "int":
		return "int"
	case "float":
		return "float64"
	case "string":
		return "string"
	case "bool":
		return "bool"
	default:
		return "interface{}"
	}
}

func goStub(sig Signature) string {
	var params []string
	for _, p := range sig.Params {
		params = append(params, fmt.Sprintf("%s %s", p.Name, goType(p.Type)))
	}

	return fmt.Sprintf("func %s(%s) %s {\n\n}\n",
		sig.FunctionName, strings.Join(params, ", "), goType(sig.ReturnType))
}

func goHarness(sig Signature, userCode string) string {
	var decls strings.Builder
	var callArgs []string

	for i, p := range sig.Params {
		varName := fmt.Sprintf("arg%d", i)
		fmt.Fprintf(&decls, "\tvar %s %s\n", varName, goType(p.Type))
		fmt.Fprintf(&decls, "\tjson.Unmarshal(argsRaw[%d], &%s)\n", i, varName)
		callArgs = append(callArgs, varName)
	}

	return fmt.Sprintf(`package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

%s

func main() {
	input, _ := io.ReadAll(os.Stdin)

	var argsRaw []json.RawMessage
	json.Unmarshal(input, &argsRaw)

%s
	result := %s(%s)

	out, _ := json.Marshal(result)
	fmt.Println(string(out))
}
`, userCode, decls.String(), sig.FunctionName, strings.Join(callArgs, ", "))
}

func goRun(ctx context.Context, code string, cfg runConfig) (string, error) {
	return runInContainer(ctx, goImage, "main.go", code,
		[]string{"go", "run", "main.go"},
		cfg)
}
