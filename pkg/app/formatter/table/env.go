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

package table

import (
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/types"
)

func RenderEnvironments(w io.Writer, envs []types.EnvdEnvironment) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Name", "Endpoint", "SSH Target", "Image",
		"GPU", "CUDA", "CUDNN", "Status",
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
		envRow[1] = formatter.FormatEndpoint(env)
		envRow[2] = fmt.Sprintf("%s.envd", env.Name)
		envRow[3] = env.Spec.Image
		envRow[4] = strconv.FormatBool(env.GPU)
		envRow[5] = formatter.StringOrNone(env.CUDA)
		envRow[6] = formatter.StringOrNone(env.CUDNN)
		envRow[7] = env.Status.Phase
		table.Append(envRow)
	}
	table.Render()
}

func RenderPortBindings(w io.Writer, ports []types.PortBinding) {
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

func RenderDependencies(w io.Writer, dep *types.Dependency) {
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
