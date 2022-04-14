package starlark

import (
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

type Interpreter interface {
	Eval(script string) (interface{}, error)
}

// generalInterpreter is the interpreter implementation for Starlark.
// Please refer to https://github.com/google/starlark-go
type generalInterpreter struct {
	*starlark.Thread
}

func NewInterpreter() Interpreter {
	starlark.Universe["midi"] = Module
	return &generalInterpreter{
		Thread: &starlark.Thread{Load: repl.MakeLoad()},
	}
}

func (s generalInterpreter) Eval(script string) (interface{}, error) {
	globals, err := starlark.ExecFile(s.Thread, "", script, nil)
	if err != nil {
		repl.PrintError(err)
		return globals, err
	}
	return globals, nil
}
