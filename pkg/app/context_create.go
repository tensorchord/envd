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
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandContextCreate = &cli.Command{
	Name:  "create",
	Usage: "Create envd context",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "Name of the context",
			Value:    "",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "builder",
			Usage: "Builder to use (docker-container, kube-pod, tcp, unix, moby-worker, nerdctl-container)",
			Value: string(types.BuilderTypeDocker),
		},
		&cli.StringFlag{
			Name:  "builder-address",
			Usage: "Builder address",
			Value: "envd_buildkitd",
		},
		&cli.StringFlag{
			Name:  "runner",
			Usage: "Runner to use(docker, envd-server)",
			Value: string(types.RunnerTypeDocker),
		},
		&cli.StringFlag{
			Name:  "runner-address",
			Usage: "Runner address",
		},
		&cli.BoolFlag{
			Name:  "use",
			Usage: "Use the context",
		},
	},
	Action: contextCreate,
}

func contextCreate(clicontext *cli.Context) error {
	name := clicontext.String("name")
	builder := clicontext.String("builder")
	builderAddress := clicontext.String("builder-address")
	runner := clicontext.String("runner")
	runnerAddress := clicontext.String("runner-address")
	use := clicontext.Bool("use")

	logger := logrus.WithFields(logrus.Fields{
		"cmd":            "context create",
		"name":           name,
		"builder":        builder,
		"builderAddress": builderAddress,
		"runner":         runner,
		"runnerAddress":  runnerAddress,
		"use":            use,
	})

	c := types.Context{
		Name:           name,
		Builder:        types.BuilderType(builder),
		BuilderAddress: builderAddress,
		Runner:         types.RunnerType(runner),
	}
	if runnerAddress != "" {
		c.RunnerAddress = &runnerAddress
	}

	err := home.GetManager().ContextCreate(c, use)
	if err != nil {
		return errors.Wrap(err, "failed to create context")
	}
	logger.Infof("Context %s is created", name)
	if use {
		logger.Infof(`Current context is "%s"`, name)
	}
	return nil
}
