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
	"go.starlark.net/starlark"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/config"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/data"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/install"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/io"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/runtime"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/universe"
)

type Interpreter interface {
	Eval(script string) (interface{}, error)
	ExecFile(filename string, funcname string) (interface{}, error)
}

type entry struct {
	globals starlark.StringDict
	err error
}

// generalInterpreter is the interpreter implementation for Starlark.
// Please refer to https://github.com/google/starlark-go
type generalInterpreter struct {
	*starlark.Thread
	predeclared     starlark.StringDict
	buildContextDir string
	cache map[string]*entry
}

func NewInterpreter(buildContextDir string) Interpreter {
	// Register envd rules and built-in variables to Starlark.
	universe.RegisterEnvdRules()
	universe.RegisterBuildContext(buildContextDir)

	return &generalInterpreter{
		predeclared: starlark.StringDict{
			"install": install.Module,
			"config":  config.Module,
			"io":      io.Module,
			"runtime": runtime.Module,
			"data":    data.Module,
		},
		buildContextDir: buildContextDir,
		cache: make(map[string]*entry),
	}
}

func (s *generalInterpreter) NewThread(module string) *starlark.Thread {
	thread := &starlark.Thread{
		Name: module,
		Load: s.load,
	}
	return thread
}

func (s *generalInterpreter) load(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	return s.exec(thread, module)
}

func (s *generalInterpreter) exec(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	e, ok := s.cache[module]
	if e != nil {
		return e.globals, e.err
	}
	if ok {
		return nil, errors.Newf("Detect cycling import during parsing %s", module)
	}

	s.cache[module] = nil
	// TODO: find the data
	var data string
	globals, err := starlark.ExecFile(thread, module, data, s.predeclared)
	e = &entry{globals, err}
	return globals, err
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
	thread := s.NewThread(filename)
	globals, err := s.exec(thread, filename)
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
