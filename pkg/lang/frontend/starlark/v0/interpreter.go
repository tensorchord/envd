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

package v0

import (
	"bytes"
	"hash/fnv"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"

	interp "github.com/tensorchord/envd/pkg/lang/frontend/starlark"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v0/config"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v0/data"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v0/install"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v0/io"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v0/runtime"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v0/universe"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

type entry struct {
	globals starlark.StringDict
	err     error
}

// generalInterpreter is the interpreter implementation for Starlark.
// Please refer to https://github.com/google/starlark-go
type generalInterpreter struct {
	predeclared     starlark.StringDict
	buildContextDir string
	cache           map[string]*entry
}

func NewInterpreter(buildContextDir string) interp.Interpreter {
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
		cache:           make(map[string]*entry),
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
	// There are two cases:
	// 1. module exists
	// 2. there's an explicit `nil` placeholder for module in s.cache
	// 3. module does not exist in s.cache
	e, ok := s.cache[module]

	// Case 1.
	if e != nil {
		return e.globals, e.err
	}

	// Case 2.
	// There is an explicit `nil` for module, which means we are in the middle of exec module.
	if ok {
		return nil, errors.Newf("Detected cycle import during parsing %s", module)
	}

	// Case 3.
	// Add a placeholder to indicate "load in progress".
	s.cache[module] = nil

	if !strings.HasPrefix(module, universe.GitPrefix) {
		var data interface{}
		globals, err := starlark.ExecFile(thread, module, data, s.predeclared)
		e = &entry{globals, err}
	} else {
		// exec remote git repo
		url := module[len(universe.GitPrefix):]
		path, err := fileutil.DownloadOrUpdateGitRepo(url)
		if err != nil {
			return nil, err
		}
		globals, err := s.loadGitModule(thread, path)
		e = &entry{globals, err}
	}

	// Update the cache.
	s.cache[module] = e

	return e.globals, e.err
}

func (s *generalInterpreter) loadGitModule(thread *starlark.Thread, path string) (globals starlark.StringDict, err error) {
	var src interface{}
	globals = starlark.StringDict{}
	logger := logrus.WithField("file", thread.Name)
	logger.Debugf("load git module from: %s", path)
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".envd") {
			return nil
		}
		dict, err := starlark.ExecFile(thread, path, src, s.predeclared)
		if err != nil {
			return err
		}
		for key, val := range dict {
			if _, exist := globals[key]; exist {
				return errors.Newf("found duplicated object name: %s in %s", key, path)
			}
			if !strings.HasPrefix(key, "_") {
				globals[key] = val
			}
		}
		return nil
	})
	return
}

func (s generalInterpreter) ExecFile(filename string, funcname string) (interface{}, error) {
	logrus.WithField("filename", filename).Debug("interpret the file")
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
				_, err := starlark.Call(thread, fn, nil, nil)
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
	thread := s.NewThread(script)
	return starlark.ExecFile(thread, "", script, s.predeclared)
}

func GetEnvdProgramHash(filename string) (string, error) {
	envdSrc, err := os.ReadFile(filename)
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
