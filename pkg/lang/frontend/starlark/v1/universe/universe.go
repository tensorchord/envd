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
	"fmt"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v1/builtin"
	ir "github.com/tensorchord/envd/pkg/lang/ir/v1"
	"github.com/tensorchord/envd/pkg/util/starlarkutil"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

// RegisterEnvdRules registers built-in envd rules into the global namespace.
func RegisterEnvdRules() {
	starlark.Universe[ruleBase] = starlark.NewBuiltin(ruleBase, ruleFuncBase)
	starlark.Universe[ruleShell] = starlark.NewBuiltin(ruleShell, ruleFuncShell)
	starlark.Universe[ruleRun] = starlark.NewBuiltin(ruleRun, ruleFuncRun)
	starlark.Universe[ruleGitConfig] = starlark.NewBuiltin(ruleGitConfig, ruleFuncGitConfig)
	starlark.Universe[ruleInclude] = starlark.NewBuiltin(ruleInclude, ruleFuncInclude)
}

func RegisterBuildContext(buildContextDir string) {
	starlark.Universe[builtin.BuildContextDir] = starlark.String(buildContextDir)
}

func ruleFuncBase(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var image string
	var dev bool

	if err := starlark.UnpackArgs(ruleBase, args, kwargs, "image?", &image, "dev?", &dev); err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, image=%s, dev=%t\n", ruleBase, image, dev)

	err := ir.Base(image, dev)
	return starlark.None, err
}

func ruleFuncRun(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var commands *starlark.List
	mountHost := false

	if err := starlark.UnpackArgs(ruleRun,
		args, kwargs, "commands", &commands, "mount_host?", &mountHost); err != nil {
		return nil, err
	}

	goCommands, err := starlarkutil.ToStringSlice(commands)
	if err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, commands=%v, mount_host=%t", ruleRun, goCommands, mountHost)
	if err := ir.Run(goCommands, mountHost); err != nil {
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

func ruleFuncInclude(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var gitRepo string

	if err := starlark.UnpackArgs(ruleInclude,
		args, kwargs, "git?", &gitRepo); err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, git=%s", ruleInclude, gitRepo)

	globals, err := thread.Load(thread, fmt.Sprintf("%s%s", GitPrefix, gitRepo))
	if err != nil {
		return nil, err
	}
	module := &starlarkstruct.Module{
		Name:    gitRepo,
		Members: globals,
	}
	return module, nil
}
