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
	"github.com/tensorchord/envd/pkg/lang/ir"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandUp = &cli.Command{
	Name:     "up",
	Category: CategoryBasic,
	Aliases:  []string{"u"},
	Usage:    "Build and run the envd environment",
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
		&cli.BoolFlag{
			Name:    "use-proxy",
			Usage:   "Use HTTPS_PROXY/HTTP_PROXY/NO_PROXY in the build process",
			Aliases: []string{"proxy"},
			Value:   false,
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
	c, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}

	buildOpt, err := ParseBuildOpt(clicontext)
	if err != nil {
		return errors.Wrap(err, "failed to parse the build options")
	}

	// Always push image to registry when envd-server is the runner.
	if c.Runner == types.RunnerTypeEnvdServer {
		buildOpt.OutputOpts = fmt.Sprintf("type=image,name=%s,push=true", buildOpt.Tag)
	}
	r := string(c.Runner)
	telemetry.GetReporter().Telemetry("up", &r)

	ctr := filepath.Base(buildOpt.BuildContextDir)
	detach := clicontext.Bool("detach")
	logger := logrus.WithFields(logrus.Fields{
		"builder-options": buildOpt,
		"container-name":  ctr,
		"detach":          detach,
	})
	logger.Debug("starting up command")

	builder, err := GetBuilder(clicontext, buildOpt)
	if err != nil {
		return err
	}
	if err = InterpretEnvdDef(builder); err != nil {
		return err
	}
	if err = DetectEnvironment(clicontext, buildOpt); err != nil {
		return err
	}
	if err = BuildImage(clicontext, builder); err != nil {
		return err
	}

	logrus.Debug("start running the environment")
	// Do not attach GPU if the flag is set.
	gpuEnable := clicontext.Bool("no-gpu")
	var gpu bool
	if gpuEnable {
		gpu = false
	} else {
		gpu = builder.GPUEnabled()
	}
	numGPU := 0
	if gpu {
		numGPU = 1
	}

	opt := envd.Options{
		Context: c,
	}
	engine, err := envd.New(clicontext.Context, opt)
	if err != nil {
		return errors.Wrap(err, "failed to create the docker client")
	}

	startOptions := envd.StartOptions{
		EnvironmentName: filepath.Base(buildOpt.BuildContextDir),
		BuildContext:    buildOpt.BuildContextDir,
		Image:           buildOpt.Tag,
		NumGPU:          numGPU,
		Forced:          clicontext.Bool("force"),
		Timeout:         clicontext.Duration("timeout"),
	}
	if c.Runner != types.RunnerTypeEnvdServer {
		startOptions.EngineSource = envd.EngineSource{
			DockerSource: &envd.DockerSource{
				Graph:        *ir.DefaultGraph,
				MountOptions: clicontext.StringSlice("volume"),
			},
		}
	} else if c.Runner == types.RunnerTypeEnvdServer {
		startOptions.EnvdServerSource = &envd.EnvdServerSource{}
	}

	res, err := engine.StartEnvd(clicontext.Context, startOptions)
	if err != nil {
		return errors.Wrap(err, "failed to start the envd environment")
	}
	logrus.Debugf("container %s is running", res.Name)

	logrus.Debugf("add entry %s to SSH config.", ctr)
	hostname, err := c.GetSSHHostname()
	if err != nil {
		return errors.Wrap(err, "failed to get the ssh hostname")
	}

	eo, err := engine.GenerateSSHConfig(ctr, hostname,
		clicontext.Path("private-key"), res)
	if err != nil {
		return errors.Wrap(err, "failed to get the ssh entry")
	}
	if err = sshconfig.AddEntry(eo); err != nil {
		logrus.Infof("failed to add entry %s to your SSH config file: %s", ctr, err)
		return errors.Wrap(err, "failed to add entry to your SSH config file")
	}

	if !detach {
		if err := engine.Attach(ctr, hostname,
			clicontext.Path("private-key"), res); err != nil {
			return errors.Wrap(err, "failed to attach to the ssh target")
		}
	}

	return nil
}
