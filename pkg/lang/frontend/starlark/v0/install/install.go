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
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	ir "github.com/tensorchord/envd/pkg/lang/ir/v0"
	"github.com/tensorchord/envd/pkg/util/starlarkutil"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

var Module = &starlarkstruct.Module{
	Name: "install",
	Members: starlark.StringDict{
		"python_packages":   starlark.NewBuiltin(rulePyPIPackage, ruleFuncPyPIPackage),
		"r_packages":        starlark.NewBuiltin(ruleRPackage, ruleFuncRPackage),
		"apt_packages":      starlark.NewBuiltin(ruleSystemPackage, ruleFuncSystemPackage),
		"cuda":              starlark.NewBuiltin(ruleCUDA, ruleFuncCUDA),
		"vscode_extensions": starlark.NewBuiltin(ruleVSCode, ruleFuncVSCode),
		"conda_packages":    starlark.NewBuiltin(ruleConda, ruleFuncConda),
		"julia_packages":    starlark.NewBuiltin(ruleJulia, ruleFuncJulia),
	},
}

func ruleFuncPyPIPackage(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name *starlark.List
	var requirementsFile starlark.String
	var wheels *starlark.List

	if err := starlark.UnpackArgs(rulePyPIPackage, args, kwargs,
		"name?", &name, "requirements?", &requirementsFile, "local_wheels?", &wheels); err != nil {
		return nil, err
	}

	nameList, err := starlarkutil.ToStringSlice(name)
	if err != nil {
		return nil, err
	}

	requirementsFileStr := requirementsFile.GoString()

	localWheels, err := starlarkutil.ToStringSlice(wheels)
	if err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, name=%v, requirements=%s, local_wheels=%s",
		rulePyPIPackage, nameList, requirementsFileStr, localWheels)

	err = ir.PyPIPackage(nameList, requirementsFileStr, localWheels)
	return starlark.None, err
}

func ruleFuncRPackage(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name *starlark.List

	if err := starlark.UnpackArgs(ruleRPackage,
		args, kwargs, "name", &name); err != nil {
		return nil, err
	}

	nameList, err := starlarkutil.ToStringSlice(name)
	if err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, name=%v", ruleRPackage, nameList)
	ir.RPackage(nameList)

	return starlark.None, nil
}

func ruleFuncJulia(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name *starlark.List

	if err := starlark.UnpackArgs(ruleJulia,
		args, kwargs, "name", &name); err != nil {
		return nil, err
	}

	nameList, err := starlarkutil.ToStringSlice(name)
	if err != nil {
		return nil, err
	}
	logger.Debugf("rule `%s` is invoked, name=%v", ruleJulia, nameList)
	ir.JuliaPackage(nameList)

	return starlark.None, nil
}

func ruleFuncSystemPackage(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name *starlark.List

	if err := starlark.UnpackArgs(ruleSystemPackage,
		args, kwargs, "name?", &name); err != nil {
		return nil, err
	}

	nameList, err := starlarkutil.ToStringSlice(name)
	if err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, name=%v", ruleSystemPackage, nameList)
	ir.SystemPackage(nameList)

	return starlark.None, nil
}

func ruleFuncCUDA(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var version, cudnn string

	if err := starlark.UnpackArgs(ruleCUDA, args, kwargs,
		"version", &version, "cudnn?", &cudnn); err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, version=%s, cudnn=%s",
		ruleCUDA, version, cudnn)
	ir.CUDA(version, cudnn)

	return starlark.None, nil
}

func ruleFuncVSCode(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var plugins *starlark.List

	if err := starlark.UnpackArgs(ruleVSCode,
		args, kwargs, "name", &plugins); err != nil {
		return nil, err
	}

	pluginList, err := starlarkutil.ToStringSlice(plugins)
	if err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, plugins=%v", ruleVSCode, pluginList)
	if err := ir.VSCodePlugins(pluginList); err != nil {
		return starlark.None, err
	}

	return starlark.None, nil
}

func ruleFuncConda(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name, channel *starlark.List
	var envFile starlark.String

	if err := starlark.UnpackArgs(ruleConda,
		args, kwargs, "name?", &name, "channel?", &channel, "env_file?", &envFile); err != nil {
		return nil, err
	}

	nameList, err := starlarkutil.ToStringSlice(name)
	if err != nil {
		return nil, err
	}

	channelList, err := starlarkutil.ToStringSlice(channel)
	if err != nil {
		return nil, err
	}

	envFileStr := envFile.GoString()
	if envFileStr != "" {
		if (len(nameList) != 0) || (len(channelList) != 0) {
			return nil, errors.New("env_file and name/channel are mutually exclusive")
		}
	}

	logger.Debugf("rule `%s` is invoked, name=%v, channel=%v, env_file=%s", ruleConda, nameList, channelList, envFileStr)
	if err := ir.CondaPackage(nameList, channelList, envFileStr); err != nil {
		return starlark.None, err
	}

	return starlark.None, nil
}
