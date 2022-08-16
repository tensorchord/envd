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

package runtime

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/tensorchord/envd/pkg/lang/ir"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

var Module = &starlarkstruct.Module{
	Name: "runtime",
	Members: starlark.StringDict{
		"command": starlark.NewBuiltin(ruleCommand, ruleFuncCommand),
		"daemon":  starlark.NewBuiltin(ruleDaemon, ruleFuncDaemon),
	},
}

func ruleFuncCommand(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var commands starlark.IterableMapping

	if err := starlark.UnpackArgs(ruleCommand, args, kwargs,
		"commands?", &commands); err != nil {
		return nil, err
	}

	commandsMap := make(map[string]string)
	for _, tuple := range commands.Items() {
		if len(tuple) != 2 {
			return nil, fmt.Errorf("invalid command in %s", ruleCommand)
		}

		commandsMap[tuple[0].(starlark.String).GoString()] =
			tuple[1].(starlark.String).GoString()
	}

	logger.Debugf("rule `%s` is invoked, commands: %v",
		ruleCommand, commandsMap)

	ir.RuntimeCommands(commandsMap)
	return starlark.None, nil
}

func ruleFuncDaemon(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var commands *starlark.List

	if err := starlark.UnpackArgs(ruleDaemon, args, kwargs, "commands", &commands); err != nil {
		return nil, err
	}

	commandList := [][]string{}
	if commands != nil {
		for i := 0; i < commands.Len(); i++ {
			args, ok := commands.Index(i).(*starlark.List)
			if !ok {
				logger.Warnf("cannot parse %s into a list of string", commands.Index(i).String())
				continue
			}
			argList := []string{}
			for j := 0; j < args.Len(); j++ {
				argList = append(argList, args.Index(j).(starlark.String).GoString())
			}
			commandList = append(commandList, argList)
		}

		logger.Debugf("rule `%s` is invoked, commands=%v", ruleDaemon, commandList)
		ir.RuntimeDaemon(commandList)
	}
	return starlark.None, nil
}
