// Copyright 2022 The MIDI Authors
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
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/MIDI/pkg/builder"
	"github.com/tensorchord/MIDI/pkg/docker"
	"github.com/tensorchord/MIDI/pkg/home"
	"github.com/tensorchord/MIDI/pkg/lang/ir"
	"github.com/tensorchord/MIDI/pkg/ssh"
)

var CommandUp = &cli.Command{
	Name:    "up",
	Aliases: []string{"u"},
	Usage:   "build and run the MIDI environment",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "volume",
			Usage:   "Mount host directory into container",
			Aliases: []string{"v"},
		},
		&cli.StringFlag{
			Name:    "tag",
			Usage:   "Name and optionally a tag in the 'name:tag' format",
			Aliases: []string{"t"},
			Value:   "midi:dev",
		},
		&cli.PathFlag{
			Name:    "file",
			Usage:   "Name of the build.MIDI (Default is 'PATH/build.MIDI')",
			Aliases: []string{"f"},
			Value:   "./build.MIDI",
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
	path, err := filepath.Abs(clicontext.Path("file"))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build file")
	}
	if path == "" {
		return errors.New("file does not exist")
	}

	config := home.GetManager().ConfigFile()

	tag := clicontext.String("tag")

	builder, err := builder.New(clicontext.Context, config, path, tag)
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
	containerID, containerIP, err := dockerClient.StartMIDI(
		clicontext.Context, tag, "midi", gpu, *ir.DefaultGraph, clicontext.Duration("timeout"), clicontext.StringSlice("volume"))
	if err != nil {
		return err
	}
	logrus.Debugf("container %s is running", containerID)

	sshClient, err := ssh.NewClient(
		containerIP, "root", 2222, clicontext.Bool("auth"), clicontext.Path("private-key"), "")
	if err != nil {
		return err
	}
	if err := sshClient.Attach(); err != nil {
		return err
	}

	return nil
}
