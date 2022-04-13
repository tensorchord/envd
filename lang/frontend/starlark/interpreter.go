package starlark

import (
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

type Interpreter interface {
	Eval(script string) (interface{}, error)
}

// StarlarkGo is the interpreter implementation for Starlark.
// Please refer to https://github.com/google/starlark-go
type StarlarkGo struct {
	*starlark.Thread
}

func NewInterpreter() Interpreter {
	starlark.Universe["midi"] = Module
	return &StarlarkGo{
		Thread: &starlark.Thread{Load: repl.MakeLoad()},
	}
}

func (s StarlarkGo) Eval(script string) (interface{}, error) {
	globals, err := starlark.ExecFile(s.Thread, "", script, nil)
	if err != nil {
		repl.PrintError(err)
		return globals, err
	}
	return globals, nil
}
