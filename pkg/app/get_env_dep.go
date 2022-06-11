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
	"io"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/olekukonko/tablewriter"
	"github.com/tensorchord/envd/pkg/envd"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
	cli "github.com/urfave/cli/v2"
)

var CommandGetEnvironmentDependency = &cli.Command{
	Name:    "deps",
	Aliases: []string{"dep", "d"},
	Usage:   "List all dependencies",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "env",
			Usage:   "Specify the envd environment to use",
			Aliases: []string{"e"},
		},
		&cli.PathFlag{
			Name:    "private-key",
			Usage:   "Path to the private key",
			Aliases: []string{"k"},
			Value:   sshconfig.GetPrivateKey(),
		},
	},
	Action: getEnvironmentDependency,
}

func getEnvironmentDependency(clicontext *cli.Context) error {
	envName := clicontext.String("env")
	if envName == "" {
		return errors.New("env is required")
	}
	envdEngine, err := envd.New(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to create envd engine")
	}

	dep, err := envdEngine.ListEnvDependency(clicontext.Context, envName)
	if err != nil {
		return errors.Wrap(err, "failed to list dependencies")
	}
	renderDependencies(dep, os.Stdout)
	return nil
}

func renderDependencies(dep *types.Dependency, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Dependency", "Type"})

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	if dep == nil {
		return
	}
	for _, p := range dep.PyPIPackages {
		envRow := make([]string, 2)
		envRow[0] = p
		envRow[1] = "Python"
		table.Append(envRow)
	}
	for _, p := range dep.APTPackages {
		envRow := make([]string, 2)
		envRow[0] = p
		envRow[1] = "APT"
		table.Append(envRow)
	}
	table.Render()
}
