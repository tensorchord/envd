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

package main

import (
	"io"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/envd"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
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
			Value:   defaultEnvName(),
		},
		&cli.PathFlag{
			Name:    "private-key",
			Usage:   "Path to the private key",
			Aliases: []string{"k"},
			Value:   sshconfig.GetPrivateKey(),
		},
		&cli.BoolFlag{
			Name:    "full",
			Usage:   "Show full dependency information",
			Aliases: []string{"f"},
		},
	},
	Action: getEnvironmentDependency,
}

func defaultEnvName() string {
	name, err := fileutil.RootDir()
	// TODO(gaocegege): https://github.com/tensorchord/envd/issues/210
	// remove the panic.
	if err != nil {
		panic(err)
	}
	return name
}

func getEnvironmentDependency(clicontext *cli.Context) error {
	envName := clicontext.String("env")
	envdEngine, err := envd.New(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to create envd engine")
	}
	full := clicontext.Bool("full")
	if full {
		output, err := envdEngine.ListEnvFullDependency(clicontext.Context, envName, clicontext.Path("private-key"))
		if err != nil {
			return errors.Wrap(err, "failed to list dependencies")
		}
		logrus.Infof("%s", output)
	} else {
		dep, err := envdEngine.ListEnvDependency(clicontext.Context, envName)
		if err != nil {
			return errors.Wrap(err, "failed to list dependencies")
		}
		renderDependencies(dep, os.Stdout)
	}
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
