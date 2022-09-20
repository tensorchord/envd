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

	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/envd"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
)

var CommandDestroy = &cli.Command{
	Name:     "destroy",
	Category: CategoryBasic,
	Aliases:  []string{"d"},
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
	dockerClient, err := docker.NewClient(clicontext.Context)
	if err != nil {
		return err
	}
	var ctrName string
	if name != "" {
		ctrName = name
	} else {
		buildContext, err := filepath.Abs(path)
		if err != nil {
			return errors.Wrap(err, "failed to get absolute path of the build context")
		}
		ctrName = filepath.Base(buildContext)
	}

	tag, err := getContainerTag(clicontext, ctrName)
	if err != nil {
		return err
	}

	if ctrName, err := dockerClient.Destroy(clicontext.Context, ctrName); err != nil {
		return errors.Wrapf(err, "failed to destroy the environment: %s", ctrName)
	} else if ctrName != "" {
		logrus.Infof("%s is destroyed", ctrName)
	}
	if err := dockerClient.RemoveImage(clicontext.Context, tag); err != nil {
		return errors.Wrapf(err, "failed to remove the image: %s", tag)
	}
	if err = sshconfig.RemoveEntry(ctrName); err != nil {
		logrus.Infof("failed to remove entry %s from your SSH config file: %s", ctrName, err)
		return errors.Wrap(err, "failed to remove entry from your SSH config file")
	}
	return nil
}

func getContainerTag(clicontext *cli.Context, name string) (string, error) {
	envdEngine, err := envd.New(clicontext.Context)
	if err != nil {
		return "", err
	}
	envs, err := envdEngine.ListEnvironment(clicontext.Context)
	if err != nil {
		return "", err
	}
	for _, env := range envs {
		if env.Name == name {
			return env.Container.Image, nil
		}
	}
	return "", errors.Newf("cannot find the image of %s", name)
}
