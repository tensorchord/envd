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

package app

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/ssh"
)

var CommandRun = &cli.Command{
	Name:     "exec",
	Category: CategoryExpert,
	Usage:    "Spawns a command installed into the environment.",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "name",
			Usage:   "Name of the environment",
			Aliases: []string{"n"},
		},
		&cli.StringFlag{
			Name:    "command",
			Usage:   "Command defined in build.envd to execute",
			Aliases: []string{"c"},
		},
		&cli.PathFlag{
			Name:    "from",
			Usage:   "Function to execute, format `file:func`",
			Aliases: []string{"f"},
			Value:   "build.envd:build",
		},
		&cli.PathFlag{
			Name:    "path",
			Usage:   "Path to the directory containing the build.envd",
			Aliases: []string{"p"},
			Value:   ".",
		},
		&cli.StringFlag{
			Name:    "raw",
			Usage:   "Raw command to execute",
			Aliases: []string{"r"},
		},
	},

	Action: exec,
}

func exec(clicontext *cli.Context) error {
	name := clicontext.String("name")
	command := clicontext.String("command")
	rawCommand := clicontext.String("raw")
	path := clicontext.String("path")
	defer func(start time.Time) {
		telemetry.GetReporter().Telemetry(
			"exec", telemetry.AddField("duration", time.Since(start).Seconds()))
	}(time.Now())

	if command != "" && rawCommand != "" {
		return errors.New("--raw and --command are mutually exclusive and may only be used once")
	}

	resultCommand := rawCommand
	if name == "" {
		name = filepath.Base(path)
	}

	// Check if the container is running.
	context, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}
	envdOpt := envd.Options{
		Context: context,
	}
	engine, err := envd.New(clicontext.Context, envdOpt)
	if err != nil {
		return errors.Wrap(err, "failed to create the docker client")
	}
	if isRunning, err :=
		engine.IsRunning(clicontext.Context, name); err != nil {
		return errors.Wrapf(
			err, "failed to check if the environment %s is running", name)
	} else if !isRunning {
		return errors.Newf("the environment %s is not running", name)
	}

	if command != "" {
		rg, err := engine.ListEnvRuntimeGraph(clicontext.Context, name)
		if err != nil {
			return errors.Wrapf(err, "failed to get environment: %s", name)
		}

		logrus.Debugf("runtime commands: %s", rg.RuntimeCommands)
		if cmd, ok := rg.RuntimeCommands[command]; !ok {
			return errors.Newf("command %s does not exist", command)
		} else {
			resultCommand = cmd
		}
	}

	opt, err := ssh.GetOptions(name)
	if err != nil {
		return errors.Wrap(err, "failed to get the ssh options")
	}
	// SSH into the container and execute the command.
	sshClient, err := ssh.NewClient(*opt)
	if err != nil {
		return errors.Wrap(err, "failed to get the ssh client")
	}
	if bytes, err := sshClient.ExecWithOutput(resultCommand); err != nil {
		fmt.Fprintln(clicontext.App.Writer, string(bytes))
		return errors.Wrapf(err,
			"failed to execute the command `%s`", resultCommand)
	} else {
		fmt.Fprint(clicontext.App.Writer, string(bytes))
	}
	return nil
}
