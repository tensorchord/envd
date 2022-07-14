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
			Usage: "Builder to use (docker-container, kube-pod, tcp)",
			Value: string(types.BuilderTypeDocker),
		},
		&cli.StringFlag{
			Name:  "builder-socket",
			Usage: "Builder socket",
			Value: "envd_buildkitd",
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
	builderSocket := clicontext.String("builder-socket")
	use := clicontext.Bool("use")

	err := home.GetManager().ContextCreate(name,
		types.BuilderType(builder), builderSocket, use)
	if err != nil {
		return errors.Wrap(err, "failed to create context")
	}
	logrus.Infof("Context %s is created", name)
	if use {
		logrus.Infof("Current context is now \"%s\"", name)
	}
	return nil
}
