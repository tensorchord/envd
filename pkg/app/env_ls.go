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
	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
)

var CommandListEnv = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls", "l"},
	Usage:   "List envd environments",
	Flags: []cli.Flag{
		&formatter.FormatFlag,
	},
	Action: getEnvironment,
}

func getEnvironment(clicontext *cli.Context) error {
	context, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}

	telemetry.GetReporter().Telemetry("environment_list", telemetry.AddField("runner", context.Runner))
	opt := envd.Options{
		Context: context,
	}
	envdEngine, err := envd.New(clicontext.Context, opt)
	if err != nil {
		return err
	}
	envs, err := envdEngine.ListEnvironment(clicontext.Context)
	if err != nil {
		return err
	}
	format := clicontext.String("format")
	switch format {
	case "table":
		table.RenderEnvironments(os.Stdout, envs)
	case "json":
		return json.PrintEnvironments(envs)
	}
	return nil
}
