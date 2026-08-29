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
	needsList := false
	for _, p := range sig.Params {
		if strings.Contains(p.Type, "[]") {
			needsList = true
		}
	}
	if strings.Contains(sig.ReturnType, "[]") {
		needsList = true
	}

	var params []string
	for _, p := range sig.Params {
		params = append(params, fmt.Sprintf("%s: %s", p.Name, pythonType(p.Type)))
	}

	returnType := pythonType(sig.ReturnType)

	var b strings.Builder
	if needsList {
		b.WriteString("from typing import List\n\n")
	}
	b.WriteString(fmt.Sprintf("def %s(%s) -> %s:\n    pass\n",
		sig.FunctionName, strings.Join(params, ", "), returnType))

	return b.String()
}
