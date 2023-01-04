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

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/go-units"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/driver/docker"
)

var CommandPruneImages = &cli.Command{
	Name:   "prune",
	Usage:  "Remove unused images",
	Action: pruneImages,
}

func pruneImages(clicontext *cli.Context) error {
	telemetry.GetReporter().Telemetry("image_prune", telemetry.AddField("runner", "docker"))

	cli, err := docker.NewClient(clicontext.Context)
	if err != nil {
		return err
	}
	report, err := cli.PruneImage(clicontext.Context)
	if err != nil {
		return err
	}
	if len(report.ImagesDeleted) > 0 {
		renderPruneReport(os.Stdout, report)
	}

	return nil
}

func renderPruneReport(w io.Writer, report dockertypes.ImagesPruneReport) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Type", "Image"})

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

	for _, img := range report.ImagesDeleted {
		envRow := make([]string, 2)
		if img.Untagged != "" {
			envRow[0] = "Untagged"
			envRow[1] = img.Untagged
		} else {
			envRow[0] = "Deleted"
			envRow[1] = img.Deleted
		}
		table.Append(envRow)
	}
	table.Render()
	fmt.Fprintln(w, "Total reclaimed space:", units.HumanSize(float64(report.SpaceReclaimed)))
}
