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
	"go.starlark.net/starlark"

	"github.com/tensorchord/MIDI/pkg/lang/ir"
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
	starlark.Universe[ruleVSCode] = starlark.NewBuiltin(ruleVSCode, ruleFuncVSCode)
	starlark.Universe[ruleUbuntuAPT] = starlark.NewBuiltin(ruleUbuntuAPT, ruleFuncUbuntuAPT)
	starlark.Universe[rulePyPIMirror] = starlark.NewBuiltin(rulePyPIMirror, ruleFuncPyPIMirror)
}

func ruleFuncBase(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

func ruleFuncSystemPackage(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

func ruleFuncPyPIPackage(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

func ruleFuncCUDA(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

func ruleFuncVSCode(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var plugins *starlark.List

	if err := starlark.UnpackArgs(ruleVSCode,
		args, kwargs, "plugins?", &plugins); err != nil {
		return nil, err
	}

	pluginList := []string{}
	if plugins != nil {
		for i := 0; i < plugins.Len(); i++ {
			pluginList = append(pluginList, plugins.Index(i).(starlark.String).GoString())
		}
	}

	logger.Debugf("rule `%s` is invoked, plugins=%v", ruleVSCode, pluginList)
	if err := ir.VSCodePlugins(pluginList); err != nil {
		return starlark.None, err
	}

	return starlark.None, nil
}

func ruleFuncUbuntuAPT(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var mode, source starlark.String

	if err := starlark.UnpackArgs(ruleUbuntuAPT, args, kwargs,
		"mode?", &mode, "source?", &source); err != nil {
		return nil, err
	}

	modeStr := ""
	if mode != starlark.String("") {
		modeStr = mode.GoString()
	}
	sourceStr := ""
	if source != starlark.String("") {
		sourceStr = source.GoString()
	}

	logger.Debugf("rule `%s` is invoked, mode=%s, source=%s", ruleUbuntuAPT,
		modeStr, sourceStr)
	ir.UbuntuAPT(modeStr, sourceStr)

	return starlark.None, nil
}

func ruleFuncPyPIMirror(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var mode, mirror starlark.String

	if err := starlark.UnpackArgs(rulePyPIMirror, args, kwargs,
		"mode?", &mode, "mirror?", &mirror); err != nil {
		return nil, err
	}

	modeStr := ""
	if mode != starlark.String("") {
		modeStr = mode.GoString()
	}
	mirrorStr := ""
	if mirror != starlark.String("") {
		mirrorStr = mirror.GoString()
	}

	logger.Debugf("rule `%s` is invoked, mode=%s, mirror=%s", rulePyPIMirror,
		modeStr, mirrorStr)
	ir.PyPIMirror(modeStr, mirrorStr)

	return starlark.None, nil
}
