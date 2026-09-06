package language

import (
	"context"
	"fmt"
	"strings"
)

const javaImage = "sandbox-java:latest"

func javaType(t string) string {
	if strings.HasSuffix(t, "[]") {
		return javaType(strings.TrimSuffix(t, "[]")) + "[]"
	}
	switch t {
	case "int":
		return "int"
	case "float":
		return "double"
	case "string":
		return "String"
	case "bool":
		return "boolean"
	default:
		return "Object"
	}
}

func javaStub(sig Signature) string {
	var params []string
	for _, p := range sig.Params {
		params = append(params, fmt.Sprintf("%s %s", javaType(p.Type), p.Name))
	}

	return fmt.Sprintf("public %s %s(%s) {\n\n}\n",
		javaType(sig.ReturnType), sig.FunctionName, strings.Join(params, ", "))
}

func javaHarness(sig Signature, userCode string) string {
	var decls strings.Builder
	var callArgs []string

	for i, p := range sig.Params {
		varName := fmt.Sprintf("arg%d", i)
		decls.WriteString(fmt.Sprintf("        %s %s = %s;\n",
			javaType(p.Type), varName, javaGetter(p.Type, fmt.Sprintf("argsArray.get(%d)", i))))
		callArgs = append(callArgs, varName)
	}

	return fmt.Sprintf(`import org.json.JSONArray;
import java.util.Scanner;

public class Main {

    static int[] toIntArray(JSONArray a) {
        int[] r = new int[a.length()];
        for (int i = 0; i < a.length(); i++) r[i] = a.getInt(i);
        return r;
    }

    static String[] toStringArray(JSONArray a) {
        String[] r = new String[a.length()];
        for (int i = 0; i < a.length(); i++) r[i] = a.getString(i);
        return r;
    }

    static boolean[] toBoolArray(JSONArray a) {
        boolean[] r = new boolean[a.length()];
        for (int i = 0; i < a.length(); i++) r[i] = a.getBoolean(i);
        return r;
    }

    static double[] toDoubleArray(JSONArray a) {
        double[] r = new double[a.length()];
        for (int i = 0; i < a.length(); i++) r[i] = a.getDouble(i);
        return r;
    }

%s

    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        String input = scanner.hasNextLine() ? scanner.nextLine() : "";
        JSONArray argsArray = new JSONArray(input);

        Main solution = new Main();
%s
        Object result = solution.%s(%s);
        System.out.println(result);
    }
}
`, userCode, decls.String(), sig.FunctionName, strings.Join(callArgs, ", "))
}

func javaGetter(t, source string) string {
	if strings.HasSuffix(t, "[]") {
		inner := strings.TrimSuffix(t, "[]")
		switch inner {
		case "int":
			return fmt.Sprintf("toIntArray((JSONArray) %s)", source)
		case "string":
			return fmt.Sprintf("toStringArray((JSONArray) %s)", source)
		case "bool":
			return fmt.Sprintf("toBoolArray((JSONArray) %s)", source)
		case "float":
			return fmt.Sprintf("toDoubleArray((JSONArray) %s)", source)
		}
	}
	switch t {
	case "int":
		return fmt.Sprintf("((Number) %s).intValue()", source)
	case "float":
		return fmt.Sprintf("((Number) %s).doubleValue()", source)
	case "string":
		return fmt.Sprintf("(String) %s", source)
	case "bool":
		return fmt.Sprintf("(Boolean) %s", source)
	}
	return source
}

func javaRun(ctx context.Context, code string, cfg runConfig) (string, error) {
	return runInContainer(ctx, javaImage, "Main.java", code,
		[]string{"sh", "-c", "javac -cp /opt/lib/json.jar -d /tmp Main.java && java -cp /tmp:/opt/lib/json.jar Main"},
		cfg)
}
