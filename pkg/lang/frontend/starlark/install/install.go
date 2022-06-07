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

package install

import (
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

var Module = &starlarkstruct.Module{
	Name: "install",
	Members: starlark.StringDict{
		"pip_package": starlark.NewBuiltin(
			rulePyPIPackage, ruleFuncPyPIPackage),
		"sys_package": starlark.NewBuiltin(
			ruleSystemPackage, ruleFuncSystemPackage),
		"apt_package": starlark.NewBuiltin(ruleUbuntuAPT, ruleFuncUbuntuAPT),
		"cuda":        starlark.NewBuiltin(ruleCUDA, ruleFuncCUDA),
		"vscode":      starlark.NewBuiltin(ruleVSCode, ruleFuncVSCode),
		"shell":       starlark.NewBuiltin(ruleVSCode, ruleFuncVSCode),
	},
}

func ruleFuncPyPIPackage(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name *starlark.List

	if err := starlark.UnpackArgs(rulePyPIPackage,
		args, kwargs, "name", &name); err != nil {
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
	if err := ir.UbuntuAPT(modeStr, sourceStr); err != nil {
		return nil, err
	}

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
