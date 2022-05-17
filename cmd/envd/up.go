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
	"fmt"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/builder"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/ssh"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandUp = &cli.Command{
	Name:    "up",
	Aliases: []string{"u"},
	Usage:   "build and run the envd environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "tag",
			Usage:   "Name and optionally a tag in the 'name:tag' format",
			Aliases: []string{"t"},
			Value:   "envd:dev",
		},
		&cli.PathFlag{
			Name:    "path",
			Usage:   "Path to the directory containing the build.envd",
			Aliases: []string{"p"},
			Value:   ".",
		},
		&cli.StringSliceFlag{
			Name:    "volume",
			Usage:   "Mount host directory into container",
			Aliases: []string{"v"},
		},
		&cli.PathFlag{
			Name:    "file",
			Usage:   "Name of the build.envd",
			Aliases: []string{"f"},
			Value:   "build.envd",
		},
		&cli.BoolFlag{
			Name:  "auth",
			Usage: "Enable authentication for ssh",
			Value: false,
		},
		&cli.PathFlag{
			Name:    "private-key",
			Usage:   "Path to the private key",
			Aliases: []string{"k"},
			Value:   "~/.ssh/id_rsa",
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "Timeout of container creation",
			Value: time.Second * 30,
		},
	},

	Action: up,
}

func up(clicontext *cli.Context) error {
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build context")
	}

	manifest, err := filepath.Abs(filepath.Join(buildContext, clicontext.Path("file")))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build file")
	}
	if manifest == "" {
		return errors.New("file does not exist")
	}

	config := home.GetManager().ConfigFile()

	tag := clicontext.String("tag")
	if tag == "" {
		logrus.Debug("tag not specified, using default")
		tag = fmt.Sprintf("%s:%s", fileutil.Base(buildContext), "dev")
	}
	ctr := fileutil.Base(buildContext)

	logger := logrus.WithFields(logrus.Fields{
		"build-context":  buildContext,
		"build-file":     manifest,
		"config":         config,
		"tag":            tag,
		"container-name": ctr,
	})
	logger.Debug("starting up")

	builder, err := builder.New(clicontext.Context, config, manifest, buildContext, tag)
	if err != nil {
		return errors.Wrap(err, "failed to create the builder")
	}

	if err := builder.Build(clicontext.Context); err != nil {
		return err
	}
	gpu := builder.GPUEnabled()

	dockerClient, err := docker.NewClient(clicontext.Context)
	if err != nil {
		return err
	}
	containerID, containerIP, err := dockerClient.StartEnvd(clicontext.Context,
		tag, ctr, buildContext, gpu, *ir.DefaultGraph, clicontext.Duration("timeout"),
		clicontext.StringSlice("volume"))
	if err != nil {
		return err
	}
	logrus.Debugf("container %s is running", containerID)

	sshClient, err := ssh.NewClient(
		containerIP, "envd", 2222, clicontext.Bool("auth"), clicontext.Path("private-key"), "")
	if err != nil {
		return err
	}
	if err := sshClient.Attach(); err != nil {
		return err
	}

	return nil
}
