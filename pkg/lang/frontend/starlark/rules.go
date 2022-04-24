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
	"github.com/tensorchord/MIDI/pkg/lang/ir"
	"go.starlark.net/starlark"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

// registerMIDIRules registers built-in MIDI rules into the global namespace.
func registerMIDIRules() {
	starlark.Universe[ruleBase] = starlark.NewBuiltin(ruleBase, ruleFuncBase)
	starlark.Universe[rulePyPIPackage] = starlark.NewBuiltin(
		rulePyPIPackage, ruleFuncPyPIPackage)
	starlark.Universe[ruleSystemPackage] = starlark.NewBuiltin(
		ruleSystemPackage, ruleFuncSystemPackage)
	starlark.Universe[ruleCUDA] = starlark.NewBuiltin(ruleCUDA, ruleFuncCUDA)
}

func ruleFuncBase(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var os, language starlark.String

	if err := starlark.UnpackArgs(ruleBase, args, kwargs, "os?", &os, "language?", &language); err != nil {
		return nil, err
	}

	osStr := ""
	if os != starlark.String("") {
		osStr = os.GoString()
	}
	langStr := ""
	if language != starlark.String("") {
		langStr = language.GoString()
	}

	logger.Debugf("rule `%s` is invoked, os=%s, language=%s", ruleBase,
		osStr, langStr)
	ir.Base(osStr, langStr)

	return starlark.None, nil
}

func ruleFuncSystemPackage(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name *starlark.List

	if err := starlark.UnpackArgs(ruleSystemPackage,
		args, kwargs, "name?", &name); err != nil {
		return nil, err
	}

	nameList := []string{}
	if name != nil {
		for i := 0; i < name.Len(); i++ {
			nameList = append(nameList, name.Index(i).(starlark.String).GoString())
		}
	}

	logger.Debugf("rule `%s` is invoked, name=%v", ruleSystemPackage, nameList)
	ir.SystemPackage(nameList)

	return starlark.None, nil
}

func ruleFuncPyPIPackage(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name *starlark.List

	if err := starlark.UnpackArgs(rulePyPIPackage,
		args, kwargs, "name?", &name); err != nil {
		return nil, err
	}

	nameList := []string{}
	if name != nil {
		for i := 0; i < name.Len(); i++ {
			nameList = append(nameList, name.Index(i).(starlark.String).GoString())
		}
	}

	logger.Debugf("rule `%s` is invoked, name=%v", rulePyPIPackage, nameList)
	ir.PyPIPackage(nameList)

	return starlark.None, nil
}

func ruleFuncCUDA(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var version, cudnn starlark.String

	if err := starlark.UnpackArgs(ruleCUDA, args, kwargs,
		"version?", &version, "cudnn?", &cudnn); err != nil {
		return nil, err
	}

	versionStr := ""
	if version != starlark.String("") {
		versionStr = version.GoString()
	}
	cudnnStr := ""
	if cudnn != starlark.String("") {
		cudnnStr = cudnn.GoString()
	}

	logger.Debugf("rule `%s` is invoked, version=%s, cudnn=%s", ruleCUDA,
		versionStr, cudnnStr)
	ir.CUDA(versionStr, cudnnStr)

	return starlark.None, nil
}
