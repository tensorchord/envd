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
	"os"

	"github.com/cockroachdb/errors"
	"github.com/tensorchord/envd/pkg/envd"
	cli "github.com/urfave/cli/v2"
)

var CommandGetImageDependency = &cli.Command{
	Name:    "deps",
	Aliases: []string{"dep", "d"},
	Usage:   "List all dependencies in the image",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "image",
			Usage:   "Specify the image to use",
			Aliases: []string{"i"},
		},
	},
	Action: getImageDependency,
}

func getImageDependency(clicontext *cli.Context) error {
	envName := clicontext.String("image")
	if envName == "" {
		return errors.New("image is required")
	}
	envdEngine, err := envd.New(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to create envd engine")
	}
	dep, err := envdEngine.ListImageDependency(clicontext.Context, envName)
	if err != nil {
		return errors.Wrap(err, "failed to list dependencies")
	}
	renderDependencies(dep, os.Stdout)
	return nil
}
