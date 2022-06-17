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
	"github.com/spf13/viper"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/builder"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/ssh"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
	"github.com/tensorchord/envd/pkg/util/netutil"
)

const (
	localhost = "127.0.0.1"
)

var CommandUp = &cli.Command{
	Name:    "up",
	Aliases: []string{"u"},
	Usage:   "build and run the envd environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "tag",
			Usage:   "Name and optionally a tag in the 'name:tag' format (default: PROJECT:dev)",
			Aliases: []string{"t"},
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
		// &cli.BoolFlag{
		// 	Name:  "auth",
		// 	Usage: "Enable authentication for ssh",
		// 	Value: false,
		// },
		&cli.PathFlag{
			Name:    "private-key",
			Usage:   "Path to the private key",
			Aliases: []string{"k"},
			Value:   sshconfig.GetPrivateKey(),
		},
		&cli.PathFlag{
			Name:    "public-key",
			Usage:   "Path to the public key",
			Aliases: []string{"pubk"},
			Value:   sshconfig.GetPublicKey(),
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "Timeout of container creation",
			Value: time.Second * 30,
		},
		&cli.BoolFlag{
			Name:  "detach",
			Usage: "detach from the container",
			Value: false,
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

	detach := clicontext.Bool("detach")

	logger := logrus.WithFields(logrus.Fields{
		"build-context":             buildContext,
		"build-file":                manifest,
		"config":                    config,
		"tag":                       tag,
		"container-name":            ctr,
		"detach":                    detach,
		flag.FlagBuildkitdImage:     viper.GetString(flag.FlagBuildkitdImage),
		flag.FlagBuildkitdContainer: viper.GetString(flag.FlagBuildkitdContainer),
	})
	logger.Debug("starting up command")
	debug := clicontext.Bool("debug")
	builder, err := builder.New(clicontext.Context, config, manifest, buildContext, tag, "", debug)
	if err != nil {
		return errors.Wrap(err, "failed to create the builder")
	}

	if err := builder.Build(clicontext.Context, clicontext.Path("public-key")); err != nil {
		return errors.Wrap(err, "failed to build the image")
	}
	gpu := builder.GPUEnabled()

	dockerClient, err := docker.NewClient(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to create the docker client")
	}

	if gpu {
		nvruntimeExists, err := dockerClient.GPUEnabled(clicontext.Context)
		if err != nil {
			return errors.Wrap(err, "failed to check if nvidia-runtime is installed")
		}
		if !nvruntimeExists {
			return errors.New("GPU is required but nvidia container runtime is not installed, please refer to https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#docker")
		}
	}

	sshPort, err := netutil.GetFreePort()
	if err != nil {
		return errors.Wrap(err, "failed to get a free port")
	}
	numGPUs := builder.NumGPUs()

	containerID, containerIP, err := dockerClient.StartEnvd(clicontext.Context,
		tag, ctr, buildContext, gpu, numGPUs, sshPort, *ir.DefaultGraph, clicontext.Duration("timeout"),
		clicontext.StringSlice("volume"))
	if err != nil {
		return errors.Wrap(err, "failed to start the envd environment")
	}
	logrus.Debugf("container %s is running", containerID)

	logrus.Debugf("Add entry %s to SSH config. at %s", buildContext, containerIP)
	if err = sshconfig.AddEntry(
		ctr, localhost, sshPort, clicontext.Path("private-key")); err != nil {
		logrus.Infof("failed to add entry %s to your SSH config file: %s", ctr, err)
		return errors.Wrap(err, "failed to add entry to your SSH config file")
	}

	if !detach {
		sshClient, err := ssh.NewClient(
			localhost, "envd", sshPort, true, clicontext.Path("private-key"), "")
		if err != nil {
			return errors.Wrap(err, "failed to create the ssh client")
		}
		if err := sshClient.Attach(); err != nil {
			return errors.Wrap(err, "failed to attach to the container")
		}
	}

	return nil
}
