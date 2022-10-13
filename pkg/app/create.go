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
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/ssh"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
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
		&cli.PathFlag{
			Name:    "private-key",
			Usage:   "Path to the private key",
			Aliases: []string{"k"},
			Value:   sshconfig.GetPrivateKeyOrPanic(),
			Hidden:  true,
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

	opt := envd.StartOptions{
		Image:   clicontext.String("image"),
		Timeout: clicontext.Duration("timeout"),
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
	hostname, err := c.GetSSHHostname()
	if err != nil {
		return errors.Wrap(err, "failed to get the ssh hostname")
	}

	eo := sshconfig.EntryOptions{
		Name:               res.Name,
		IFace:              hostname,
		Port:               res.SSHPort,
		PrivateKeyPath:     clicontext.Path("private-key"),
		EnableHostKeyCheck: false,
		EnableAgentForward: false,
		User:               res.Name,
	}
	if err = sshconfig.AddEntry(eo); err != nil {
		logrus.Infof("failed to add entry %s to your SSH config file: %s", res.Name, err)
		return errors.Wrap(err, "failed to add entry to your SSH config file")
	}

	if !clicontext.Bool("detach") {
		opt := ssh.DefaultOptions()
		opt.PrivateKeyPath = clicontext.Path("private-key")
		opt.Port = res.SSHPort
		opt.AgentForwarding = false
		opt.User = res.Name
		opt.Server = hostname

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
