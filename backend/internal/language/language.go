package language

import (
	"context"
	"fmt"
)

type Param struct {
	Name string
	Type string
}

type Signature struct {
	FunctionName string
	Params       []Param
	ReturnType   string
}

type langDef struct {
	Stub          func(Signature) string
	Harness       func(Signature, string) string
	Run           func(ctx context.Context, code string, cfg runConfig) (string, error)
	Image         string
	DockerContext string
}

var registry = map[string]langDef{
	"python": {
		Stub:          pythonStub,
		Harness:       pythonHarness,
		Run:           pythonRun,
		Image:         pythonImage,
		DockerContext: "internal/language/dockerfiles/python",
	},
	"javascript": {
		Stub:          javascriptStub,
		Harness:       javascriptHarness,
		Run:           javascriptRun,
		Image:         nodeImage,
		DockerContext: "internal/language/dockerfiles/node",
	},
	"typescript": {
		Stub:          typescriptStub,
		Harness:       typescriptHarness,
		Run:           typescriptRun,
		Image:         nodeImage,
		DockerContext: "internal/language/dockerfiles/node",
	},
	"java": {
		Stub:          javaStub,
		Harness:       javaHarness,
		Run:           javaRun,
		Image:         javaImage,
		DockerContext: "internal/language/dockerfiles/java",
	},
	"cpp": {
		Stub:          cppStub,
		Harness:       cppHarness,
		Run:           cppRun,
		Image:         cppImage,
		DockerContext: "internal/language/dockerfiles/cpp",
	},
}

func Registry() map[string]langDef {
	return registry
}

func Stubs(sig Signature) map[string]string {
	stubs := make(map[string]string, len(registry))
	for name, def := range registry {
		stubs[name] = def.Stub(sig)
	}
	return stubs
}

func IsValid(s string) bool {
	_, ok := registry[s]
	return ok
}

func Stub(lang string, sig Signature) (string, bool) {
	d, ok := registry[lang]
	if !ok {
		return "", false
	}
	return d.Stub(sig), true
}

func Harness(lang string, sig Signature, userCode string) (string, bool) {
	d, ok := registry[lang]
	if !ok {
		return "", false
	}
	return d.Harness(sig, userCode), true
}

func Run(ctx context.Context, lang, code, stdin string) (string, error) {
	d, ok := registry[lang]
	if !ok {
		return "", fmt.Errorf("unsupported language")
	}
	cfg := defaultRunConfig(lang)
	cfg.stdin = stdin
	return d.Run(ctx, code, cfg)
}
