package starlark

type Interpreter interface {
	Eval(script string) (interface{}, error)
	ExecFile(filename string, funcname string) (interface{}, error)
}
