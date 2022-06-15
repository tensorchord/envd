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
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"

	"github.com/tensorchord/envd/pkg/lang/ir"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

// registerenvdRules registers built-in envd rules into the global namespace.
func RegisterenvdRules() {
	starlark.Universe[ruleBase] = starlark.NewBuiltin(ruleBase, ruleFuncBase)
	starlark.Universe[ruleShell] = starlark.NewBuiltin(ruleShell, ruleFuncShell)
	starlark.Universe[ruleRun] = starlark.NewBuiltin(ruleRun, ruleFuncRun)
	starlark.Universe[ruleGitConfig] = starlark.NewBuiltin(ruleGitConfig, ruleFuncGitConfig)
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
	err := ir.Base(osStr, langStr)
	if err != nil {
		return starlark.None, err
	}

	return starlark.None, nil
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

	shellStr := ""
	if shell != starlark.String("") {
		shellStr = shell.GoString()
	}

	logger.Debugf("rule `%s` is invoked, shell=%s", ruleShell,
		shellStr)
	if err := ir.Shell(shellStr); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func ruleFuncGitConfig(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name, email, editor starlark.String

	if err := starlark.UnpackArgs(ruleGitConfig,
		args, kwargs, "name?", &name, "email?", &email, "editor?", &editor); err != nil {
		return nil, err
	}

	nameStr := ""
	if name != starlark.String("") {
		nameStr = name.GoString()
	}

	emailStr := ""
	if email != starlark.String("") {
		nameStr = email.GoString()
	}

	editorStr := ""
	if editor != starlark.String("") {
		editorStr = editor.GoString()
	}

	logger.Debugf("rule `%s` is invoked, name=%s, email=%s, editor=%s",
		ruleGitConfig, nameStr, emailStr, editorStr)
	if err := ir.Git(nameStr, emailStr, editorStr); err != nil {
		return nil, err
	}

	return starlark.None, nil
}
