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
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/urfave/cli/v2"
)

var CommandPause = &cli.Command{
	Name:    "pause",
	Aliases: []string{"p"},
	Usage:   "pause the envd environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "env",
			Usage:   "environment name",
			Aliases: []string{"e"},
		},
	},

	Action: pause,
}

func pause(clicontext *cli.Context) error {
	env := clicontext.String("env")
	if env == "" {
		return errors.New("env is required")
	}
	envdEngine, err := envd.New(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to create envd engine")
	}
	if name, err := envdEngine.PauseEnvironment(clicontext.Context, env); err != nil {
		return errors.Wrap(err, "failed to pause the environment")
	} else if name != "" {
		logrus.Infof("%s is paused", name)
	}
	return nil
}
