// Copyright 2022 The MIDI Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package starlark

import (
	"github.com/sirupsen/logrus"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

type Interpreter interface {
	Eval(script string) (interface{}, error)
	ExecFile(filename string) (interface{}, error)
}

// generalInterpreter is the interpreter implementation for Starlark.
// Please refer to https://github.com/google/starlark-go
type generalInterpreter struct {
	*starlark.Thread
}

func NewInterpreter() Interpreter {
	// Register MIDI rules to Starlark.
	registerMIDIRules()
	return &generalInterpreter{
		Thread: &starlark.Thread{Load: repl.MakeLoad()},
	}
}

func (s generalInterpreter) ExecFile(filename string) (interface{}, error) {
	var src interface{}
	globals, err := starlark.ExecFile(s.Thread, filename, src, nil)
	if err != nil {
		s.printEvalError(err)
		return globals, err
	}
	return globals, nil
}

func (s generalInterpreter) Eval(script string) (interface{}, error) {
	globals, err := starlark.ExecFile(s.Thread, "", script, nil)
	if err != nil {
		s.printEvalError(err)
		return globals, err
	}
	return globals, nil
}

func (s generalInterpreter) printEvalError(err error) {
	if evalErr, ok := err.(*starlark.EvalError); ok {
		logrus.Error(evalErr.Backtrace())
	} else {
		logrus.Error(err)
	}
}
