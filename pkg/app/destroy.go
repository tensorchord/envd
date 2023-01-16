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
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	buildutil "github.com/tensorchord/envd/pkg/app/build"
	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
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

func destroy(clicontext *cli.Context) error {
	path := clicontext.Path("path")
	name := clicontext.String("name")
	if path != "" && name != "" {
		return errors.New("Cannot specify --path and --name at the same time.")
	}
	if path == "" && name == "" {
		path = "."
	}
	var ctrName string
	if name != "" {
		ctrName = name
	} else {
		buildContext, err := filepath.Abs(path)
		if err != nil {
			return errors.Wrap(err, "failed to get absolute path of the build context")
		}
		ctrName, err = buildutil.CreateEnvNameFromDir(buildContext)
		if err != nil {
			return errors.Wrap(err, "failed to create an env name")
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
	return nil
}
