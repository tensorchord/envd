// Copyright 2022 The envd Authors
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
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/config"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/install"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/universe"
)

type Interpreter interface {
	Eval(script string) (interface{}, error)
	ExecFile(filename string, funcname string) (interface{}, error)
}

// generalInterpreter is the interpreter implementation for Starlark.
// Please refer to https://github.com/google/starlark-go
type generalInterpreter struct {
	*starlark.Thread
	predeclared starlark.StringDict
}

func NewInterpreter() Interpreter {
	// Register envd rules to Starlark.
	universe.RegisterenvdRules()
	return &generalInterpreter{
		Thread: &starlark.Thread{Load: repl.MakeLoad()},
		predeclared: starlark.StringDict{
			"install": install.Module,
			"config":  config.Module,
		},
	}
}

func (s generalInterpreter) ExecFile(filename string, funcname string) (interface{}, error) {
	logrus.WithField("filename", filename).Debug("interprete the file")
	var src interface{}
	globals, err := starlark.ExecFile(s.Thread, filename, src, s.predeclared)
	if err != nil {
		return nil, err
	}
	if funcname != "" {
		logrus.Debugf("Execute %s func", funcname)
		if globals.Has(funcname) {
			buildVar := globals[funcname]
			if fn, ok := buildVar.(*starlark.Function); ok {
				_, err := starlark.Call(s.Thread, fn, nil, nil)
				if err != nil {
					return nil, errors.Wrapf(err, "Exception when exec %s func", funcname)
				}
			} else {
				return nil, errors.Errorf("%s is not a function", funcname)
			}
		} else {
			return nil, errors.Errorf("envd file doesn't has %s function", funcname)
		}

	}
	return globals, nil
}

func (s generalInterpreter) Eval(script string) (interface{}, error) {
	globals, err := starlark.ExecFile(s.Thread, "", script, s.predeclared)
	if err != nil {
		return globals, err
	}
	return globals, nil
}
