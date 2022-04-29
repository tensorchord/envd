// Copyright 2022 The MIDI Authors
// Copyright 2022 The midi Authors
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
package main

import (
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	ac "github.com/tensorchord/MIDI/pkg/autocomplete"
	"github.com/tensorchord/MIDI/pkg/buildkitd"
)

var CommandBootstrap = &cli.Command{
	Name:  "bootstrap",
	Usage: "Bootstraps midi installation including shell autocompletion and buildkit image download",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "buildkit",
			Usage:   "Download the image and bootstrap buildkit",
			Aliases: []string{"b"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:  "with-autocomplete",
			Usage: "Add midi autocompletions",
			Value: true,
		},
	},

	Action: bootstrap,
}

func bootstrap(clicontext *cli.Context) error {
	autocomplete := clicontext.Bool("with-autocomplete")
	if autocomplete {
		// Because this requires sudo, it should warn and not fail the rest of it.
		err := ac.InsertBashCompleteEntry()
		if err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
			err = nil
		}
		err = ac.InsertZSHCompleteEntry()
		if err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
			err = nil
		}

		logrus.Info("You may have to restart your shell for autocomplete to get initialized (e.g. run \"exec $SHELL\")\n")
	}

	buildkit := clicontext.Bool("buildkit")

	if buildkit {
		logrus.Debug("bootstrap the buildkitd container")
		bkClient, err := buildkitd.NewClient(clicontext.Context)
		if err != nil {
			return errors.Wrap(err, "failed to create buildkit client")
		}
		defer bkClient.Close()
		logrus.Infof("The buildkit is running at %s", bkClient.BuildkitdAddr())
	}
	return nil
}
