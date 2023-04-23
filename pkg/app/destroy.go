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
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	buildutil "github.com/tensorchord/envd/pkg/app/build"
	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/syncthing"
)

var CommandDestroy = &cli.Command{
	Name:     "destroy",
	Category: CategoryBasic,
	Aliases:  []string{"down", "d"},
	Usage:    "Destroy the envd environment",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:        "path",
			Usage:       "Path to the directory containing the build.envd",
			Aliases:     []string{"p"},
			DefaultText: "current directory",
		},
		&cli.PathFlag{
			Name:    "name",
			Usage:   "Name of the environment or container ID",
			Aliases: []string{"n"},
		},
	},

	Action: destroy,
}

// Prompts the user to confirm an operation with [Y/n].
// If the output is not tty, it will return false automatically.
func confirm(prompt string) bool {
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	if !isTerminal {
		return false
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [Y/n] ", prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = response[:len(response)-1] // Remove newline character
	return response == "y" || response == "Y" || response == "yes" || response == "Yes"
}

func destroy(clicontext *cli.Context) error {
	path := clicontext.Path("path")
	name := clicontext.String("name")
	if path != "" && name != "" {
		return errors.New("Cannot specify --path and --name at the same time.")
	}

	var ctrName string
	if name != "" {
		ctrName = name
	} else if path != "" {
		buildContext, err := filepath.Abs(path)
		if err != nil {
			return errors.Wrap(err, "failed to get absolute path of the build context")
		}
		ctrName, err = buildutil.CreateEnvNameFromDir(buildContext)
		if err != nil {
			return errors.Wrap(err, "failed to create an env name")
		}
	} else {
		// Both path and name are empty
		// Destroy the environment in the current directory only if user confirms
		buildContext, err := filepath.Abs(".")
		if err != nil {
			return errors.Wrap(err, "failed to get absolute path of the build context")
		}
		ctrName, err = buildutil.CreateEnvNameFromDir(buildContext)
		if err != nil {
			return errors.Wrap(err, "failed to create an env name")
		}
		if !confirm(fmt.Sprintf("Are you sure you want to destroy container %s in the current directory?", ctrName)) {
			return nil
		}
	}

	context, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}
	telemetry.GetReporter().Telemetry("destroy", telemetry.AddField("runner", context.Runner))

	opt := envd.Options{
		Context: context,
	}
	envdEngine, err := envd.New(clicontext.Context, opt)
	if err != nil {
		return errors.Wrap(err, "failed to create envd engine")
	}

	if ctrName, err := envdEngine.Destroy(clicontext.Context, ctrName); err != nil {
		return errors.Wrapf(err, "failed to destroy the environment: %s", ctrName)
	} else if ctrName != "" {
		logrus.Infof("environment(%s) is destroyed", ctrName)
	}

	if err = sshconfig.RemoveEntry(ctrName); err != nil {
		logrus.Infof("failed to remove entry %s from your SSH config file: %s", ctrName, err)
		return errors.Wrap(err, "failed to remove entry from your SSH config file")
	}

	err = syncthing.CleanLocalConfig(name)
	if err != nil {
		return errors.Wrap(err, "failed to remove syncthing config file")
	}

	return nil
}
