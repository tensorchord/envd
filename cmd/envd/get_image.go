package main

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/go-units"
	"github.com/olekukonko/tablewriter"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandGetImage = &cli.Command{
	Name:    "images",
	Aliases: []string{"image", "i"},
	Usage:   "List envd images",

	Subcommands: []*cli.Command{
		CommandGetImageDependency,
	},

	Action: getImage,
}

func getImage(clicontext *cli.Context) error {
	envdEngine, err := envd.New(clicontext.Context)
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
	table.SetHeader([]string{"Name", "GPU", "CUDA", "CUDNN", "Image ID", "Created", "Size"})

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
		envRow := make([]string, 7)
		envRow[0] = types.GetImageName(img)
		envRow[1] = strconv.FormatBool(img.GPU)
		envRow[2] = cudaString(img.CUDA)
		envRow[3] = cudnnString(img.CUDNN)
		envRow[4] = stringid.TruncateID(img.ImageSummary.ID)
		envRow[5] = createdSinceString(img.ImageSummary.Created)
		envRow[6] = units.HumanSizeWithPrecision(float64(img.ImageSummary.Size), 3)
		table.Append(envRow)
	}
	table.Render()
}

func cudaString(cuda string) string {
	if cuda == "" {
		return "<none>"
	}
	return cuda
}

func cudnnString(cudnn string) string {
	if cudnn == "" {
		return "<none>"
	}
	return cudnn
}

func createdSinceString(created int64) string {
	createdAt := time.Unix(created, 0)

	if createdAt.IsZero() {
		return ""
	}

	return units.HumanDuration(time.Now().UTC().Sub(createdAt)) + " ago"
}
