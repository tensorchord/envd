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

	buildutil "github.com/tensorchord/envd/pkg/app/build"
	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
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
		&cli.StringFlag{
			Name:  "name",
			Usage: "environment name",
			Value: "",
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
		&cli.IntFlag{
			Name:  "shm-size",
			Usage: "Configure the shared memory size (megabyte)",
			Value: 2048,
		},
		&cli.StringFlag{
			Name:  "cpus",
			Usage: "Request CPU resources (number of cores), such as 0.5, 1, 2",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "cpu-set",
			Usage: "Limit the specific CPUs or cores the environment can use, such as `0-3`, `1,3`",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "memory",
			Usage: "Request Memory, such as 512Mb, 2Gb",
			Value: "",
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
		&cli.StringFlag{
			Name:  "host",
			Usage: "Assign the host address for environment ssh acesss server listening",
			Value: envd.Localhost,
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

	buildOpt, err := buildutil.ParseBuildOpt(clicontext)
	if err != nil {
		return errors.Wrap(err, "failed to parse the build options")
	}

	// Always push image to registry when envd-server is the runner.
	if c.Runner == types.RunnerTypeEnvdServer {
		buildOpt.OutputOpts = fmt.Sprintf("type=image,name=%s,push=true", buildOpt.Tag)

		// Unable to modify sshd host when runner is envd-server.
		if clicontext.String("host") != envd.Localhost {
			return errors.New("Failed to modify the sshd host when runner is envd-server.")
		}
	}
	start := time.Now()

	ctr := filepath.Base(buildOpt.BuildContextDir)
	detach := clicontext.Bool("detach")
	logger := logrus.WithFields(logrus.Fields{
		"builder-options": buildOpt,
		"container-name":  ctr,
		"detach":          detach,
	})
	logger.Debug("starting up command")

	builder, err := buildutil.GetBuilder(clicontext, buildOpt)
	if err != nil {
		return err
	}
	if err = buildutil.InterpretEnvdDef(builder); err != nil {
		return err
	}
	if err = buildutil.DetectEnvironment(clicontext, buildOpt); err != nil {
		return err
	}
	if err = buildutil.BuildImage(clicontext, builder); err != nil {
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
	name := clicontext.String("name")
	if name == "" {
		name, err = buildutil.CreateEnvNameFromDir(buildOpt.BuildContextDir)
		if err != nil {
			return errors.Wrapf(err, "failed to create the env name from %s", buildOpt.BuildContextDir)
		}
	}
	startOptions := envd.StartOptions{
		EnvironmentName: name,
		BuildContext:    buildOpt.BuildContextDir,
		Image:           buildOpt.Tag,
		NumGPU:          numGPU,
		Forced:          clicontext.Bool("force"),
		Timeout:         clicontext.Duration("timeout"),
		SshdHost:        clicontext.String("host"),
		ShmSize:         clicontext.Int("shm-size"),
		NumCPU:          clicontext.String("cpus"),
		NumMem:          clicontext.String("memory"),
		CPUSet:          clicontext.String("cpu-set"),
	}
	if len(startOptions.NumCPU) > 0 && len(startOptions.CPUSet) > 0 {
		return errors.New("`--cpus` and `--cpu-set` are mutually exclusive")
	}

	if c.Runner != types.RunnerTypeEnvdServer {
		startOptions.EngineSource = envd.EngineSource{
			DockerSource: &envd.DockerSource{
				Graph:        builder.GetGraph(),
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
	hostname, err := c.GetSSHHostname(startOptions.SshdHost)
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
	telemetry.GetReporter().Telemetry(
		"up",
		telemetry.AddField("runner", c.Runner),
		telemetry.AddField("duration", time.Since(start).Seconds()))

	if !detach {
		if err := engine.Attach(ctr, hostname,
			clicontext.Path("private-key"), res, builder.GetGraph()); err != nil {
			return errors.Wrap(err, "failed to attach to the ssh target")
		}
		logrus.Infof("Detached successfully. You can attach to the container with command `ssh %s.envd`\n",
			startOptions.EnvironmentName)
	}

	return nil
}
