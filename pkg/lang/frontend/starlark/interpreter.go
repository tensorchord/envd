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
	"bytes"
	"hash/fnv"
	"io/ioutil"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/config"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/install"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/io"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/runtime"
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
	predeclared     starlark.StringDict
	buildContextDir string
}

func NewInterpreter(buildContextDir string) Interpreter {
	// Register envd rules and built-in variables to Starlark.
	universe.RegisterenvdRules()
	universe.RegisterBuildContext(buildContextDir)

	return &generalInterpreter{
		Thread: &starlark.Thread{Load: repl.MakeLoad()},
		predeclared: starlark.StringDict{
			"install": install.Module,
			"config":  config.Module,
			"io":      io.Module,
			"runtime": runtime.Module,
		},
		buildContextDir: buildContextDir,
	}
}

func GetEnvdProgramHash(filename string) (string, error) {
	envdSrc, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	// No Check builtin or predeclared for now
	funcAlwaysHas := func(x string) bool {
		return true
	}
	_, prog, err := starlark.SourceProgram(filename, envdSrc, funcAlwaysHas)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = prog.Write(buf)
	if err != nil {
		return "", err
	}
	h := fnv.New64a()
	h.Write(buf.Bytes())
	hashsum := h.Sum64()
	return strconv.FormatUint(hashsum, 16), nil
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
	return starlark.ExecFile(s.Thread, "", script, s.predeclared)
}
