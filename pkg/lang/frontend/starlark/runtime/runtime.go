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
