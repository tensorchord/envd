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
	"io"
	"strconv"

	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/go-units"
	"github.com/olekukonko/tablewriter"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/types"
)

func RenderImages(w io.Writer, imgs []types.EnvdImage) {
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
		envRow[0] = img.Name
		envRow[1] = formatter.StringOrNone(img.BuildContext)
		envRow[2] = strconv.FormatBool(img.GPU)
		envRow[3] = formatter.StringOrNone(img.CUDA)
		envRow[4] = formatter.StringOrNone(img.CUDNN)
		envRow[5] = stringid.TruncateID(img.Digest)
		envRow[6] = formatter.CreatedSinceString(img.Created)
		envRow[7] = units.HumanSizeWithPrecision(float64(img.Size), 3)
		table.Append(envRow)
	}
	table.Render()
}
