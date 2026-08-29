package stub

import (
	"fmt"
	"strings"
)

func JavaScript(sig Signature) string {
	var params []string
	for _, p := range sig.Params {
		params = append(params, p.Name)
	}

	return fmt.Sprintf("function %s(%s) {\n\n}\n",
		sig.FunctionName, strings.Join(params, ", "))
}
