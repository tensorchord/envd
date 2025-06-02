// Copyright 2023 The envd Authors
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

	"github.com/cockroachdb/errors"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/types"
)

func RenderEnvironments(w io.Writer, envs []types.EnvdEnvironment) error {
	table := CreateTable(w)
	table.Header([]string{
		"Name", "Endpoint", "SSH Target", "Image",
		"GPU", "CUDA", "CUDNN", "Status",
	})

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
		err := table.Append(envRow)
		if err != nil {
			return errors.Wrapf(err, "failed to append row for environment %s", env.Name)
		}
	}
	return errors.Wrap(table.Render(), "failed to render environment table")
}

func RenderPortBindings(w io.Writer, ports []types.PortBinding) error {
	if ports == nil {
		return nil
	}
	table := CreateTable(w)
	table.Header([]string{"Name", "Container Port", "Protocol", "Host IP", "Host Port"})
	for _, port := range ports {
		row := make([]string, 5)
		row[0] = port.Name
		row[1] = port.Port
		row[2] = port.Protocol
		row[3] = port.HostIP
		row[4] = port.HostPort
		err := table.Append(row)
		if err != nil {
			return errors.Wrapf(err, "failed to append row for port binding %s", port.Name)
		}
	}
	return errors.Wrap(table.Render(), "failed to render port bindings table")
}

func RenderDependencies(w io.Writer, dep *types.Dependency) error {
	if dep == nil {
		return nil
	}
	table := CreateTable(w)
	table.Header([]string{"Dependencies", "Type"})
	for _, p := range dep.PyPIPackages {
		envRow := make([]string, 2)
		envRow[0] = p
		envRow[1] = "Python"
		err := table.Append(envRow)
		if err != nil {
			return errors.Wrapf(err, "failed to append row for Python package %s", p)
		}
	}
	for _, p := range dep.APTPackages {
		envRow := make([]string, 2)
		envRow[0] = p
		envRow[1] = "APT"
		err := table.Append(envRow)
		if err != nil {
			return errors.Wrapf(err, "failed to append row for APT package %s", p)
		}
	}
	return errors.Wrap(table.Render(), "failed to render dependencies table")
}

func CreateTable(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewTable(
		w,
		tablewriter.WithRowAutoWrap(tw.WrapNone),
		tablewriter.WithHeaderAutoFormat(tw.On),
		tablewriter.WithRenderer(renderer.NewBlueprint(
			tw.Rendition{
				Borders: tw.BorderNone,
				Symbols: tw.NewSymbols(tw.StyleNone),
				Settings: tw.Settings{
					Separators: tw.Separators{
						BetweenRows:    tw.Off,
						BetweenColumns: tw.Off,
					},
					Lines: tw.Lines{
						ShowHeaderLine: tw.Off,
					},
				},
			},
		)),
		tablewriter.WithConfig(
			tablewriter.Config{
				Header: tw.CellConfig{
					Alignment: tw.CellAlignment{
						Global: tw.AlignLeft,
					},
				},
				Row: tw.CellConfig{
					Alignment: tw.CellAlignment{
						Global: tw.AlignLeft,
					},
				},
			},
		),
	)

	return table
}
