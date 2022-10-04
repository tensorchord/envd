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

	"github.com/cockroachdb/errors"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandContextList = &cli.Command{
	Name:   "ls",
	Usage:  "List envd contexts",
	Action: contextList,
}

func contextList(clicontext *cli.Context) error {
	contexts, err := home.GetManager().ContextList()
	if err != nil {
		return errors.Wrap(err, "failed to list context")
	}
	renderContext(contexts, os.Stdout)
	return nil
}

func renderContext(contexts types.EnvdContext, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"context", "builder", "builder addr", "runner", "runner addr"})

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

	for _, p := range contexts.Contexts {
		envRow := make([]string, 5)
		if p.Name == contexts.Current {
			envRow[0] = fmt.Sprintf("%s (current)", p.Name)
		} else {
			envRow[0] = p.Name
		}
		envRow[1] = string(p.Builder)
		envRow[2] = fmt.Sprintf("%s://%s", p.Builder, p.BuilderAddress)
		envRow[3] = string(p.Runner)
		if p.RunnerAddress != nil {
			envRow[4] = stringOrNone(*p.RunnerAddress)
		}
		table.Append(envRow)
	}
	table.Render()
}
