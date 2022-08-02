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

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/ssh"
)

var CommandRun = &cli.Command{
	Name:     "run",
	Category: CategoryBasic,
	Usage:    "Spawns a command installed into the environment.",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "name",
			Usage:   "Name of the environment",
			Aliases: []string{"n"},
		},
		&cli.StringFlag{
			Name:    "command",
			Usage:   "Command to execute",
			Aliases: []string{"c"},
		},
	},

	Action: run,
}

func run(clicontext *cli.Context) error {
	name := clicontext.String("name")

	// Check if the container is running.
	dockerClient, err := docker.NewClient(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to create the docker client")
	}
	if isRunning, err :=
		dockerClient.IsRunning(clicontext.Context, name); err != nil {
		return errors.Wrapf(
			err, "failed to check if the environment %s is running", name)
	} else if !isRunning {
		return errors.Newf("the environment %s is not running", name)
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
	if bytes, err := sshClient.ExecWithOutput(clicontext.String("command")); err != nil {
		return errors.Wrap(err, "failed to execute the command")
	} else {
		fmt.Fprint(clicontext.App.Writer, string(bytes))
	}
	return nil
}
