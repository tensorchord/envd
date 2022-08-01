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

	"github.com/tensorchord/envd/pkg/builder"
	"github.com/tensorchord/envd/pkg/docker"
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
	Usage:   "Build and run the envd environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "tag",
			Usage:       "Name and optionally a tag in the 'name:tag' format",
			Aliases:     []string{"t"},
			DefaultText: "PROJECT:dev",
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
			Name:    "from",
			Usage:   "Function to execute, format `file:func`",
			Aliases: []string{"f"},
			Value:   "build.envd:build",
		},
		&cli.PathFlag{
			Name:    "private-key",
			Usage:   "Path to the private key",
			Aliases: []string{"k"},
			Value:   sshconfig.GetPrivateKeyOrPanic(),
			Hidden:  true,
		},
		&cli.PathFlag{
			Name:    "public-key",
			Usage:   "Path to the public key",
			Aliases: []string{"pubk"},
			Value:   sshconfig.GetPublicKeyOrPanic(),
			Hidden:  true,
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "Timeout of container creation",
			Value: time.Second * 30,
		},
		&cli.BoolFlag{
			Name:  "detach",
			Usage: "Detach from the container",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "no-gpu",
			Usage: "Launch the CPU container",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force rebuild and run the container although the previous container is running",
			Value: false,
		},
		// https://github.com/urfave/cli/issues/1134#issuecomment-1191407527
		&cli.StringFlag{
			Name:    "export-cache",
			Usage:   "Export the cache (e.g. type=registry,ref=<image>)",
			Aliases: []string{"ec"},
		},
		&cli.StringFlag{
			Name:    "import-cache",
			Usage:   "Import the cache (e.g. type=registry,ref=<image>)",
			Aliases: []string{"ic"},
		},
	},

	Action: up,
}

func up(clicontext *cli.Context) error {
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build context")
	}
	fileName, funcName, err := builder.ParseFromStr(clicontext.String("from"))
	if err != nil {
		return err
	}

	manifest, err := filepath.Abs(filepath.Join(buildContext, fileName))
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
	// The current container engine is only Docker. It should be expaned to support other container engines.
	tag, err = docker.NormalizeNamed(tag)
	if err != nil {
		return err
	}
	ctr := fileutil.Base(buildContext)

	detach := clicontext.Bool("detach")
	debug := clicontext.Bool("debug")
	force := clicontext.Bool("force")
	output := ""
	exportCache := clicontext.String("export-cache")
	importCache := clicontext.String("import-cache")

	opt := builder.Options{
		ManifestFilePath: manifest,
		ConfigFilePath:   config,
		BuildFuncName:    funcName,
		BuildContextDir:  buildContext,
		Tag:              tag,
		OutputOpts:       output,
		PubKeyPath:       clicontext.Path("public-key"),
		ProgressMode:     "auto",
		ExportCache:      exportCache,
		ImportCache:      importCache,
	}
	if debug {
		opt.ProgressMode = "plain"
	}

	logger := logrus.WithFields(logrus.Fields{
		"builder-options": opt,
		"container-name":  ctr,
		"detach":          detach,
	})
	logger.Debug("starting up command")
	builder, err := builder.New(clicontext.Context, opt)
	if err != nil {
		return errors.Wrap(err, "failed to create the builder")
	}

	if err := builder.Build(clicontext.Context, force); err != nil {
		return errors.Wrap(err, "failed to build the image")
	}
	// Do not attach GPU if the flag is set.
	gpuEnable := clicontext.Bool("no-gpu")
	var gpu bool
	if gpuEnable {
		gpu = false
	} else {
		gpu = builder.GPUEnabled()
	}

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

	sshPortInHost, err := netutil.GetFreePort()
	if err != nil {
		return errors.Wrap(err, "failed to get a free port")
	}
	numGPUs := builder.NumGPUs()

	err = dockerClient.CleanEnvdIfExists(clicontext.Context, ctr, force)
	if err != nil {
		return errors.Wrap(err, "failed to start the envd environment")
	}
	containerID, containerIP, err := dockerClient.StartEnvd(clicontext.Context,
		tag, ctr, buildContext, gpu, numGPUs, sshPortInHost, *ir.DefaultGraph, clicontext.Duration("timeout"),
		clicontext.StringSlice("volume"))
	if err != nil {
		return errors.Wrap(err, "failed to start the envd environment")
	}
	logrus.Debugf("container %s is running", containerID)

	logrus.Debugf("Add entry %s to SSH config. at %s", buildContext, containerIP)
	if err = sshconfig.AddEntry(
		ctr, localhost, sshPortInHost, clicontext.Path("private-key")); err != nil {
		logrus.Infof("failed to add entry %s to your SSH config file: %s", ctr, err)
		return errors.Wrap(err, "failed to add entry to your SSH config file")
	}

	if !detach {
		opt := ssh.DefaultOptions()
		opt.PrivateKeyPath = clicontext.Path("private-key")
		opt.Port = sshPortInHost
		sshClient, err := ssh.NewClient(opt)
		if err != nil {
			return errors.Wrap(err, "failed to create the ssh client")
		}
		if err := sshClient.Attach(); err != nil {
			return errors.Wrap(err, "failed to attach to the container")
		}
	}

	return nil
}
