package stub

import (
	"fmt"
	"strings"
)

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

func TypeScript(sig Signature) string {
	var params []string
	for _, p := range sig.Params {
		params = append(params, fmt.Sprintf("%s: %s", p.Name, typescriptType(p.Type)))
	}

	returnType := typescriptType(sig.ReturnType)

	return fmt.Sprintf("function %s(%s): %s {\n\n}\n",
		sig.FunctionName, strings.Join(params, ", "), returnType)
}
