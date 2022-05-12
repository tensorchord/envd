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

package main

import (
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandDestroy = &cli.Command{
	Name:    "destroy",
	Aliases: []string{"d"},
	Usage:   "destroys the envd environment",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "path",
			Usage:   "Path to the directory containing the build.envd",
			Aliases: []string{"p"},
			Value:   ".",
		},
	},

	Action: destroy,
}

func destroy(clicontext *cli.Context) error {
	dockerClient, err := docker.NewClient(clicontext.Context)
	if err != nil {
		return err
	}

	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build context")
	}

	ctr := fileutil.Base(buildContext)

	if err := dockerClient.Destroy(clicontext.Context, ctr); err != nil {
		return errors.Wrapf(err, "failed to destroy the environment: %s", ctr)
	}
	logrus.Info("envd environment destroyed")
	return nil
}
