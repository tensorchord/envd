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
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/app/formatter/json"
	"github.com/tensorchord/envd/pkg/app/formatter/table"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
)

func getCurrentDirOrPanic() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Base(dir)
}

var CommandDescribeEnvironment = &cli.Command{
	Name:    "describe",
	Aliases: []string{"d"},
	Usage:   "Show details about environments, including dependencies and port binding",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "env",
			Usage:   "Specify the envd environment to use",
			Aliases: []string{"e"},
			Value:   getCurrentDirOrPanic(),
		},
		&formatter.FormatFlag,
	},
	Action: getEnvironmentDescriptions,
}

func getEnvironmentDescriptions(clicontext *cli.Context) error {
	envName := clicontext.String("env")
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

	dep, err := envdEngine.ListEnvDependency(clicontext.Context, envName)
	if err != nil {
		return errors.Wrap(err, "failed to list dependencies")
	}

	ports, err := envdEngine.ListEnvPortBinding(clicontext.Context, envName)
	if err != nil {
		return errors.Wrap(err, "failed to list port bindings")
	}
	format := clicontext.String("format")
	switch format {
	case "table":
		table.RenderDependencies(os.Stdout, dep)
		table.RenderPortBindings(os.Stdout, ports)
	case "json":
		return json.PrintEnvironmentDescriptions(dep, ports)
	}

	return nil
}
