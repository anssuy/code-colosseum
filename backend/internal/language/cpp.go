package language

import (
	"context"
	"fmt"
	"strings"
)

const cppImage = "sandbox-cpp:latest"

func cppType(t string) string {
	if before, ok := strings.CutSuffix(t, "[]"); ok {
		return "std::vector<" + cppType(before) + ">"
	}
	switch t {
	case "int":
		return "int"
	case "float":
		return "double"
	case "string":
		return "std::string"
	case "bool":
		return "bool"
	default:
		return "auto"
	}
}

func cppStub(sig Signature) string {
	var params []string
	for _, p := range sig.Params {
		params = append(params, fmt.Sprintf("%s %s", cppType(p.Type), p.Name))
	}

	return fmt.Sprintf("%s %s(%s) {\n\n}\n",
		cppType(sig.ReturnType), sig.FunctionName, strings.Join(params, ", "))
}

func cppHarness(sig Signature, userCode string) string {
	var decls strings.Builder
	var callArgs []string

	for i, p := range sig.Params {
		varName := fmt.Sprintf("arg%d", i)
		fmt.Fprintf(&decls, "    %s %s = args[%d].get<%s>();\n",
			cppType(p.Type), varName, i, cppType(p.Type))
		callArgs = append(callArgs, varName)
	}

	return fmt.Sprintf(`#include <iostream>
#include <vector>
#include <string>
#include <climits>
#include <algorithm>
#include "json.hpp"

using json = nlohmann::json;

%s

int main() {
    std::string input((std::istreambuf_iterator<char>(std::cin)), std::istreambuf_iterator<char>());
    json args = json::parse(input);

%s
    auto result = %s(%s);

    json output = result;
    std::cout << output.dump() << std::endl;

    return 0;
}
`, userCode, decls.String(), sig.FunctionName, strings.Join(callArgs, ", "))
}

func cppRun(ctx context.Context, code string, cfg runConfig) (string, error) {
	return runInContainer(ctx, cppImage, "main.cpp", code,
		[]string{"sh", "-c", "g++ -std=c++17 main.cpp -I/opt/json -o /tmp/main && /tmp/main"},
		cfg)
}
