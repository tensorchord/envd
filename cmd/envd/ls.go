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
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/docker/docker/pkg/stringid"
	"github.com/olekukonko/tablewriter"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandLS = &cli.Command{
	Name:    "ls",
	Aliases: []string{"l"},
	Usage:   "list envd environments",

	Action: list,
}

func list(clicontext *cli.Context) error {
	envdEngine, err := envd.New(clicontext.Context)
	if err != nil {
		return err
	}
	envs, err := envdEngine.List(clicontext.Context)
	if err != nil {
		return err
	}
	render(envs, os.Stdout)
	return nil
}

func render(envs []types.EnvdEnvironment, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "jupyter", "SSH Target", "GPU", "Status", "Container ID"})

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

	for _, env := range envs {
		envRow := make([]string, 6)
		envRow[0] = env.Name
		envRow[1] = env.JupyterAddr
		envRow[2] = fmt.Sprintf("%s.envd", env.Name)
		envRow[3] = strconv.FormatBool(env.GPU)
		envRow[4] = env.State
		envRow[5] = stringid.TruncateID(env.Container.ID)
		table.Append(envRow)
	}
	table.Render()
}
