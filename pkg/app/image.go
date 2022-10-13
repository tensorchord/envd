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
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/go-units"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandImage = &cli.Command{
	Name:     "images",
	Category: CategoryBasic,
	Aliases:  []string{"image"},
	Usage:    "Manage envd images",

	Subcommands: []*cli.Command{
		CommandDescribeImage,
		CommandListImage,
		CommandRemoveImage,
	},
}

var CommandListImage = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls", "l"},
	Usage:   "List envd images",
	Action:  getImage,
}

func getImage(clicontext *cli.Context) error {
	context, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}
	opt := envd.Options{
		Context: context,
	}
	envdEngine, err := envd.New(clicontext.Context, opt)
	if err != nil {
		return err
	}
	envs, err := envdEngine.ListImage(clicontext.Context)
	if err != nil {
		return err
	}
	renderImages(envs, os.Stdout)
	return nil
}

func renderImages(imgs []types.EnvdImage, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Context", "GPU", "CUDA", "CUDNN", "Image ID", "Created", "Size"})

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

	for _, img := range imgs {
		envRow := make([]string, 8)
		envRow[0] = types.GetImageName(img)
		envRow[1] = stringOrNone(img.BuildContext)
		envRow[2] = strconv.FormatBool(img.GPU)
		envRow[3] = stringOrNone(img.CUDA)
		envRow[4] = stringOrNone(img.CUDNN)
		envRow[5] = stringid.TruncateID(img.ImageSummary.ID)
		envRow[6] = createdSinceString(img.ImageSummary.Created)
		envRow[7] = units.HumanSizeWithPrecision(float64(img.ImageSummary.Size), 3)
		table.Append(envRow)
	}
	table.Render()
}

func stringOrNone(cuda string) string {
	if cuda == "" {
		return "<none>"
	}
	return cuda
}

func createdSinceString(created int64) string {
	createdAt := time.Unix(created, 0)

	if createdAt.IsZero() {
		return ""
	}

	return units.HumanDuration(time.Now().UTC().Sub(createdAt)) + " ago"
}
