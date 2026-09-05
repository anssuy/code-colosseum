package stub

type Param struct {
	Name string
	Type string
}

type Signature struct {
	FunctionName string
	Params       []Param
	ReturnType   string
}
