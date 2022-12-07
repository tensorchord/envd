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
	"net"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/v1/data"
	ir "github.com/tensorchord/envd/pkg/lang/ir/v1"
	"github.com/tensorchord/envd/pkg/util/fileutil"
	"github.com/tensorchord/envd/pkg/util/starlarkutil"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

var Module = &starlarkstruct.Module{
	Name: "runtime",
	Members: starlark.StringDict{
		"command": starlark.NewBuiltin(ruleCommand, ruleFuncCommand),
		"daemon":  starlark.NewBuiltin(ruleDaemon, ruleFuncDaemon),
		"expose":  starlark.NewBuiltin(ruleExpose, ruleFuncExpose),
		"environ": starlark.NewBuiltin(ruleEnviron, ruleFuncEnviron),
		"mount":   starlark.NewBuiltin(ruleMount, ruleFuncMount),
		"init":    starlark.NewBuiltin(ruleInitScript, ruleFuncInitScript),
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
			return nil, errors.Newf("invalid command in %s", ruleCommand)
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
				return nil, errors.Newf("invalid daemon commands (%s)", commands.Index(i).String())
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

func ruleFuncExpose(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		envdPort      starlark.Int
		hostPort      = starlark.MakeInt(0) // 0 means envd can randomly choose a free port
		serviceName   = starlark.String("")
		listeningAddr = starlark.String("127.0.0.1") // default to lisen only on local loopback interface
	)

	if err := starlark.UnpackArgs(ruleExpose,
		args, kwargs, "envd_port", &envdPort, "host_port?", &hostPort, "service?", &serviceName, "listen_addr?", &listeningAddr); err != nil {
		return nil, err
	}
	envdPortInt, ok := envdPort.Int64()
	if !ok || envdPortInt < 1 || envdPortInt > 65535 {
		return nil, errors.New("envd_port must be a positive integer less than 65535")
	}
	hostPortInt, ok := hostPort.Int64()
	if !ok || hostPortInt < 0 || hostPortInt > 65535 {
		return nil, errors.New("host_port must be a positive integer less than 65535")
	}
	serviceNameStr := serviceName.GoString()
	listeningAddrStr := listeningAddr.GoString()
	if net.ParseIP(listeningAddrStr) == nil {
		return nil, errors.New("listening_addr must be a valid IP address")
	}

	logger.Debugf("rule `%s` is invoked, envd_port=%d, host_port=%d, service=%s", ruleExpose, envdPortInt, hostPortInt, serviceNameStr)
	err := ir.RuntimeExpose(int(envdPortInt), int(hostPortInt), serviceNameStr, listeningAddrStr)
	return starlark.None, err
}

func ruleFuncEnviron(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var env starlark.IterableMapping
	var path *starlark.List

	if err := starlark.UnpackArgs(ruleCommand, args, kwargs,
		"env?", &env, "extra_path?", &path); err != nil {
		return nil, err
	}

	envMap := make(map[string]string)
	if env != nil {
		for _, tuple := range env.Items() {
			if len(tuple) != 2 {
				return nil, errors.Newf("invalid env (%s)", tuple.String())
			}
			envMap[tuple[0].(starlark.String).GoString()] = tuple[1].(starlark.String).GoString()
		}
	}

	pathList, err := starlarkutil.ToStringSlice(path)
	if err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, env: %v, extra_path: %v", ruleEnviron, envMap, pathList)
	ir.RuntimeEnviron(envMap, pathList)
	return starlark.None, nil
}

func ruleFuncMount(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var source starlark.Value
	var destination starlark.String

	if err := starlark.UnpackArgs(ruleMount, args, kwargs,
		"host_path?", &source, "envd_path?", &destination); err != nil {
		return nil, err
	}

	var sourceStr string
	var err error

	if v, ok := source.(*data.DataSourceValue); ok {
		err = v.Init()
		if err != nil {
			return starlark.None, err
		}
		sourceStr, err = v.GetHostDir()
		if err != nil {
			return starlark.None, err
		}
	} else if vs, ok := source.(starlark.String); ok {
		sourceStr = vs.GoString()
	} else {
		return starlark.None, errors.New("invalid data source")
	}

	destinationStr := destination.GoString()

	logger.Debugf("rule `%s` is invoked, src=%s, dest=%s",
		ruleMount, sourceStr, destinationStr)

	// Expand source directory based on host user
	usr, _ := user.Current()
	dir := usr.HomeDir
	if sourceStr == "~" {
		sourceStr = dir
	} else if strings.HasPrefix(sourceStr, "~/") {
		sourceStr = filepath.Join(dir, sourceStr[2:])
	}
	// Expand dest directory based on container user envd
	dir = fileutil.EnvdHomeDir()
	if destinationStr == "~" {
		destinationStr = dir
	} else if strings.HasPrefix(destinationStr, "~/") {
		destinationStr = fileutil.EnvdHomeDir(destinationStr[2:])
	}
	ir.Mount(sourceStr, destinationStr)

	return starlark.None, nil
}

func ruleFuncInitScript(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var commands *starlark.List

	if err := starlark.UnpackArgs(ruleCommand, args, kwargs,
		"commands?", &commands); err != nil {
		return nil, err
	}

	commandsSlice, err := starlarkutil.ToStringSlice(commands)
	if err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, commands: %v",
		ruleInitScript, commandsSlice)

	ir.RuntimeInitScript(commandsSlice)
	return starlark.None, nil
}
