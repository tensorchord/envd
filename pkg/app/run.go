// Copyright 2023 The envd Authors
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
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd-server/sshname"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/ssh"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/syncthing"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
	"github.com/tensorchord/envd/pkg/util/netutil"
)

var CommandCreate = &cli.Command{
	Name:        "run",
	Category:    CategoryBasic,
	Aliases:     []string{"c"},
	Usage:       "Run the envd environment from the existing image",
	Hidden:      false,
	Description: `run is only supported in envd-server runner currently`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "image",
			Usage:       "image name",
			DefaultText: "PROJECT:dev",
			Required:    true,
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "environment name",
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "Timeout of container creation",
			Value: time.Second * 1800,
		},
		&cli.BoolFlag{
			Name:  "detach",
			Usage: "Detach from the container",
			Value: false,
		},
		&cli.PathFlag{
			Name:    "private-key",
			Usage:   "Path to the private key",
			Aliases: []string{"k"},
			Value:   sshconfig.GetPrivateKeyOrPanic(),
			Hidden:  true,
		},
		&cli.PathFlag{
			Name:    "path",
			Usage:   "Working directory path to be used as project root",
			Aliases: []string{"p"},
			Value:   ".",
		},
		&cli.StringFlag{
			Name:  "host",
			Usage: "Assign the host address for the environment SSH access server listening",
			Value: envd.Localhost,
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
			Name:  "memory",
			Usage: "Request Memory, such as 512M, 2G",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "gpu",
			Usage: "Request GPU resources (number of gpus), such as 1, 2",
			Value: "",
		},
		&cli.BoolFlag{
			Name:  "sync",
			Usage: "Sync the local directory with the remote container",
			Value: false,
		},
		&cli.StringSliceFlag{
			Name:    "volume",
			Usage:   "Mount host directory into container",
			Aliases: []string{"v"},
		},
	},
	Action: run,
}

func run(clicontext *cli.Context) error {
	c, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get current context")
	}
	telemetry.GetReporter().Telemetry(
		"run", telemetry.AddField("runner", c.Runner))

	engine, err := envd.New(clicontext.Context, envd.Options{
		Context: c,
	})
	if err != nil {
		return err
	}

	name := clicontext.String("name")
	if name == "" {
		name = strings.ToLower(randomdata.SillyName())
	}
	opt := envd.StartOptions{
		SshdHost:        clicontext.String("host"),
		Image:           clicontext.String("image"),
		Timeout:         clicontext.Duration("timeout"),
		NumMem:          clicontext.String("memory"),
		NumCPU:          clicontext.String("cpus"),
		NumGPU:          clicontext.Int("gpu"),
		ShmSize:         clicontext.Int("shm-size"),
		EnvironmentName: name,
	}
	switch c.Runner {
	case types.RunnerTypeEnvdServer:
		opt.EnvdServerSource = &envd.EnvdServerSource{
			Sync: clicontext.Bool("sync"),
		}
		if len(clicontext.StringSlice("volume")) > 0 {
			return errors.New("volume is not supported for envd-server runner")
		}
	case types.RunnerTypeDocker:
		opt.DockerSource = &envd.DockerSource{
			MountOptions: clicontext.StringSlice("volume"),
		}

		buildContext, err := filepath.Abs(clicontext.Path("path"))
		if err != nil {
			return errors.Wrap(err, "failed to get absolute path of the build context")
		}
		opt.BuildContext = buildContext
	}

	res, err := engine.StartEnvd(clicontext.Context, opt)
	if err != nil {
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		"cmd":          "run",
		"StartOptions": opt,
		"StartResult":  res,
	})

	logger.Debugf("container %s is running", res.Name)

	logger.Debugf("add entry %s to SSH config.", res.Name)
	hostname, err := c.GetSSHHostname(opt.SshdHost)
	if err != nil {
		return errors.Wrap(err, "failed to get the ssh hostname")
	}

	ac, err := home.GetManager().AuthGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the auth information")
	}
	username, err := sshname.Username(ac.Name, res.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get the username")
	}

	eo := sshconfig.EntryOptions{
		Name:               res.Name,
		IFace:              hostname,
		Port:               res.SSHPort,
		PrivateKeyPath:     clicontext.Path("private-key"),
		EnableHostKeyCheck: false,
		EnableAgentForward: false,
		User:               username,
	}
	if err = sshconfig.AddEntry(eo); err != nil {
		logger.WithError(err).
			Infof("failed to add entry %s to your SSH config file", res.Name)
		return errors.Wrap(err, "failed to add entry to your SSH config file")
	}

	// TODO(gaocegege): Test why it fails.
	if !clicontext.Bool("detach") {
		outputChannel := make(chan error)
		opt := ssh.DefaultOptions()
		opt.PrivateKeyPath = clicontext.Path("private-key")
		opt.Port = res.SSHPort
		opt.AgentForwarding = false
		opt.User = username
		opt.Server = hostname

		sshClient, err := ssh.NewClient(opt)
		if err != nil {
			outputChannel <- errors.Wrap(err, "failed to create the ssh client")
		}

		ports := res.Ports

		for _, p := range ports {
			if p.Port == 2222 {
				continue
			}

			// TODO(gaocegege): Use one remote port.
			localPort, err := netutil.GetFreePort()
			if err != nil {
				return errors.Wrap(err, "failed to get a free port")
			}
			localAddress := fmt.Sprintf("%s:%d", "localhost", localPort)
			remoteAddress := fmt.Sprintf("%s:%d", "localhost", p.Port)
			logger.Infof(`service "%s" is listening at %s\n`, p.Name, localAddress)
			go func() {
				if err := sshClient.LocalForward(localAddress, remoteAddress); err != nil {
					outputChannel <- errors.Wrap(err, "failed to forward to local port")
				}
			}()
		}

		if clicontext.Bool("sync") {
			go func() {
				if err := sshClient.LocalForward(syncthing.DefaultRemoteAPIAddress, syncthing.DefaultRemoteAPIAddress); err != nil {
					outputChannel <- errors.Wrap(err, "failed to forward to remote api port")
				}
			}()

			go func() {
				syncthingRemoteAddr := fmt.Sprintf("127.0.0.1:%s", syncthing.ParsePortFromAddress(syncthing.DefaultRemoteDeviceAddress))
				if err := sshClient.LocalForward(syncthingRemoteAddr, syncthingRemoteAddr); err != nil {
					outputChannel <- errors.Wrap(err, "failed to forward to remote port")
				}
			}()

			go func() {
				syncthingLocalAddr := fmt.Sprintf("127.0.0.1:%s", syncthing.ParsePortFromAddress(syncthing.DefaultLocalDeviceAddress))
				if err := sshClient.RemoteForward(syncthingLocalAddr, syncthingLocalAddr); err != nil {
					outputChannel <- errors.Wrap(err, "failed to forward to local port")
				}
			}()

			localSyncthing, _, err := startSyncthing(res.Name)
			if err != nil {
				return errors.Wrap(err, "failed to start syncthing")
			}
			defer localSyncthing.StopLocalSyncthing()

		}

		go func() {
			// TODO(gaocegege): Avoid the hard code.
			if err := sshClient.Attach(); err != nil {
				outputChannel <- errors.Wrap(err, "failed to attach to the container")
			}
			outputChannel <- nil
		}()

		if err := <-outputChannel; err != nil {
			return err
		}
	}
	return nil
}

func startSyncthing(name string) (*syncthing.Syncthing, *syncthing.Syncthing, error) {
	cwd, err := fileutil.CWD()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get current working directory")
	}
	projectName := filepath.Base(cwd)

	logger := logrus.WithFields(logrus.Fields{
		"cwd":         cwd,
		"projectName": projectName,
	})

	localSyncthing, err := syncthing.InitializeLocalSyncthing(name)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to initialize local syncthing")
	}

	remoteSyncthing, err := syncthing.InitializeRemoteSyncthing()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to initialize remote syncthing")
	}
	logger.Debug("Remote syncthing initialized")

	err = syncthing.ConnectDevices(localSyncthing, remoteSyncthing)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to connect devices")
	}
	logger.Debug("Syncthing devices connected")

	err = syncthing.SyncFolder(localSyncthing, remoteSyncthing, cwd, fmt.Sprintf("%s/%s", fileutil.EnvdHomeDir(), projectName))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to sync folders")
	}

	return localSyncthing, remoteSyncthing, nil
}
