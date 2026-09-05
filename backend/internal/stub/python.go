package stub

import (
	"fmt"
	"strings"
)

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

func Python(sig Signature) string {
	var params []string
	for _, p := range sig.Params {
		params = append(params, fmt.Sprintf("%s: %s", p.Name, pythonType(p.Type)))
	}

	returnType := pythonType(sig.ReturnType)

	return fmt.Sprintf("def %s(%s) -> %s:\n    pass\n",
		sig.FunctionName, strings.Join(params, ", "), returnType)
}
