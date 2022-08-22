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
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandDescribeEnvironment = &cli.Command{
	Name:    "describe",
	Aliases: []string{"d"},
	Usage:   "Show details about environments, including dependencies and port binding",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "env",
			Usage:    "Specify the envd environment to use",
			Aliases:  []string{"e"},
			Required: true,
		},
	},
	Action: getEnvironmentDescriptions,
}

func getEnvironmentDescriptions(clicontext *cli.Context) error {
	envName := clicontext.String("env")
	envdEngine, err := envd.New(clicontext.Context)
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

	renderDependencies(os.Stdout, dep)
	renderPortBindings(os.Stdout, ports)
	return nil
}

func createTable(w io.Writer, headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetHeader(headers)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(true)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	return table
}

func renderPortBindings(w io.Writer, ports []types.PortBinding) {
	if ports == nil {
		return
	}
	table := createTable(w, []string{"Container Port", "Protocol", "Host IP", "Host Port"})
	for _, port := range ports {
		row := make([]string, 4)
		row[0] = port.Port
		row[1] = port.Protocol
		row[2] = port.HostIP
		row[3] = port.HostPort
		table.Append(row)
	}
	table.Render()
}

func renderDependencies(w io.Writer, dep *types.Dependency) {
	if dep == nil {
		return
	}
	table := createTable(w, []string{"Dependencies", "Type"})
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
