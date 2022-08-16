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

package universe

import (
	"errors"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/builtin"
	"github.com/tensorchord/envd/pkg/lang/ir"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

// RegisterenvdRules registers built-in envd rules into the global namespace.
func RegisterenvdRules() {
	starlark.Universe[ruleBase] = starlark.NewBuiltin(ruleBase, ruleFuncBase)
	starlark.Universe[ruleShell] = starlark.NewBuiltin(ruleShell, ruleFuncShell)
	starlark.Universe[ruleRun] = starlark.NewBuiltin(ruleRun, ruleFuncRun)
	starlark.Universe[ruleGitConfig] = starlark.NewBuiltin(ruleGitConfig, ruleFuncGitConfig)
	starlark.Universe[ruleExpose] = starlark.NewBuiltin(ruleExpose, ruleFuncExpose)
}

func RegisterBuildContext(buildContextDir string) {
	starlark.Universe[builtin.BuildContextDir] = starlark.String(buildContextDir)
}

func ruleFuncBase(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var os, language, image starlark.String

	if err := starlark.UnpackArgs(ruleBase, args, kwargs,
		"os?", &os, "language?", &language, "image?", &image); err != nil {
		return nil, err
	}

	osStr := os.GoString()
	langStr := language.GoString()
	imageStr := image.GoString()

	logger.Debugf("rule `%s` is invoked, os=%s, language=%s, image=%s",
		ruleBase, osStr, langStr, imageStr)

	err := ir.Base(osStr, langStr, imageStr)
	return starlark.None, err
}

func ruleFuncRun(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var commands *starlark.List

	if err := starlark.UnpackArgs(ruleRun,
		args, kwargs, "commands?", &commands); err != nil {
		return nil, err
	}

	goCommands := []string{}
	if commands != nil {
		for i := 0; i < commands.Len(); i++ {
			goCommands = append(goCommands, commands.Index(i).(starlark.String).GoString())
		}
	}

	logger.Debugf("rule `%s` is invoked, commands=%v", ruleRun, goCommands)
	if err := ir.Run(goCommands); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func ruleFuncShell(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var shell starlark.String

	if err := starlark.UnpackPositionalArgs(ruleShell, args, kwargs, 1, &shell); err != nil {
		return nil, err
	}

	shellStr := shell.GoString()

	logger.Debugf("rule `%s` is invoked, shell=%s", ruleShell, shellStr)

	err := ir.Shell(shellStr)
	return starlark.None, err
}

func ruleFuncExpose(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		envdPort    starlark.Int
		hostPort    = starlark.MakeInt(0)
		serviceName = starlark.String("")
	)

	if err := starlark.UnpackArgs(ruleExpose,
		args, kwargs, "envd_port", &envdPort, "host_port?", &hostPort, "service?", &serviceName); err != nil {
		return nil, err
	}
	envdPortInt, ok := envdPort.Int64()
	if !ok && envdPortInt < 0 && envdPortInt > 65536 {
		return nil, errors.New("envd_port must be a positive integer less than 65536")
	}
	hostPortInt, ok := hostPort.Int64()
	if !ok && hostPortInt < 0 && hostPortInt > 65536 {
		return nil, errors.New("envd_port must be a positive integer less than 65536")
	}
	serviceNameStr := serviceName.GoString()

	logger.Debugf("rule `%s` is invoked, envd_port=%d, host_port=%d, service=%s", ruleExpose, envdPortInt, hostPortInt, serviceNameStr)
	err := ir.Expose(int(envdPortInt), int(hostPortInt), serviceNameStr)
	return starlark.None, err
}

func ruleFuncGitConfig(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name, email, editor starlark.String

	if err := starlark.UnpackArgs(ruleGitConfig,
		args, kwargs, "name?", &name, "email?", &email, "editor?", &editor); err != nil {
		return nil, err
	}

	nameStr := name.GoString()
	emailStr := email.GoString()
	editorStr := editor.GoString()

	logger.Debugf("rule `%s` is invoked, name=%s, email=%s, editor=%s",
		ruleGitConfig, nameStr, emailStr, editorStr)

	err := ir.Git(nameStr, emailStr, editorStr)
	return starlark.None, err
}
