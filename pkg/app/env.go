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
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/pkg/stringid"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandEnvironment = &cli.Command{
	Name:     "envs",
	Category: CategoryBasic,
	Aliases:  []string{"env", "e"},
	Usage:    "Manage envd environments",

	Subcommands: []*cli.Command{
		CommandDescribeEnvironment,
		CommandListEnv,
	},
}

var CommandListEnv = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls", "l"},
	Usage:   "List envd environments",
	Action:  getEnvironment,
}

func getEnvironment(clicontext *cli.Context) error {
	envdEngine, err := envd.New(clicontext.Context)
	if err != nil {
		return err
	}
	envs, err := envdEngine.ListEnvironment(clicontext.Context)
	if err != nil {
		return err
	}
	renderEnvironments(envs, os.Stdout)
	return nil
}

func renderEnvironments(envs []types.EnvdEnvironment, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Name", "Endpoint", "SSH Target", "Image",
		"GPU", "CUDA", "CUDNN", "Status", "Container ID",
	})

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
		envRow := make([]string, 9)
		envRow[0] = env.Name
		envRow[1] = endpointOrNone(env)
		envRow[2] = fmt.Sprintf("%s.envd", env.Name)
		envRow[3] = env.Container.Image
		envRow[4] = strconv.FormatBool(env.GPU)
		envRow[5] = stringOrNone(env.CUDA)
		envRow[6] = stringOrNone(env.CUDNN)
		envRow[7] = env.Status
		envRow[8] = stringid.TruncateID(env.Container.ID)
		table.Append(envRow)
	}
	table.Render()
}

func endpointOrNone(env types.EnvdEnvironment) string {
	var res strings.Builder
	if env.JupyterAddr != nil {
		res.WriteString(fmt.Sprintf("jupyter: %s", *env.JupyterAddr))
	}
	if env.RStudioServerAddr != nil {
		res.WriteString(fmt.Sprintf("rstudio: %s", *env.RStudioServerAddr))
	}
	return res.String()
}
