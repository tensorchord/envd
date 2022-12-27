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
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/app/formatter/json"
	"github.com/tensorchord/envd/pkg/app/formatter/table"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandDescribeImage = &cli.Command{
	Name:    "describe",
	Aliases: []string{"d"},
	Usage:   "Show details about image, including dependencies",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "image",
			Usage:   "Specify the image to use",
			Aliases: []string{"i"},
		},
		&formatter.FormatFlag,
	},
	Action: getImageDependency,
}

func getImageDependency(clicontext *cli.Context) error {
	envName := clicontext.String("image")
	if envName == "" {
		return errors.New("image is required")
	}
	context, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}
	opt := envd.Options{
		Context: context,
	}
	envdEngine, err := envd.New(clicontext.Context, opt)
	if err != nil {
		return errors.Wrap(err, "failed to create envd engine")
	}
	dep, err := envdEngine.ListImageDependency(clicontext.Context, envName)
	if err != nil {
		return errors.Wrap(err, "failed to list dependencies")
	}
	format := clicontext.String("format")
	switch format {
	case "table":
		table.RenderDependencies(os.Stdout, dep)
	case "json":
		return json.PrintEnvironmentDescriptions(dep, []types.PortBinding{})
	}
	return nil
}
