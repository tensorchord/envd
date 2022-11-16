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
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
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
	table := createTable(w, []string{"Name", "Container Port", "Protocol", "Host IP", "Host Port"})
	for _, port := range ports {
		row := make([]string, 5)
		row[0] = port.Name
		row[1] = port.Port
		row[2] = port.Protocol
		row[3] = port.HostIP
		row[4] = port.HostPort
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
