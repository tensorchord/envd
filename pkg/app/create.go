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
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd-server/sshname"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/ssh"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/netutil"
)

var CommandCreate = &cli.Command{
	Name:        "create",
	Category:    CategoryBasic,
	Aliases:     []string{"c"},
	Usage:       "Create the envd environment from the existing image",
	Hidden:      true,
	Description: ``,
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
		&cli.StringFlag{
			Name:  "host",
			Usage: "Assign the host address for environment ssh acesss server listening",
			Value: envd.Localhost,
		},
	},
	Action: create,
}

func create(clicontext *cli.Context) error {
	c, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return err
	}

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
		EnvironmentName: name,
	}
	if c.Runner == types.RunnerTypeEnvdServer {
		opt.EnvdServerSource = &envd.EnvdServerSource{}
	}
	res, err := engine.StartEnvd(clicontext.Context, opt)
	if err != nil {
		return err
	}

	logrus.Debugf("container %s is running", res.Name)

	logrus.Debugf("add entry %s to SSH config.", res.Name)
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
		logrus.Infof("failed to add entry %s to your SSH config file: %s", res.Name, err)
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
			logrus.Infof("service \"%s\" is listening at %s\n", p.Name, localAddress)
			go func() {
				if err := sshClient.LocalForward(localAddress, remoteAddress); err != nil {
					outputChannel <- errors.Wrap(err, "failed to forward to local port")
				}
			}()
		}

		go func() {
			if err := sshClient.Attach(ir.DefaultGraph.Shell,
				ir.DefaultGraph.EnvironmentName); err != nil {
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
