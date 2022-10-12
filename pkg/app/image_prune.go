package app

import (
	"fmt"
	"io"
	"os"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/go-units"
	"github.com/olekukonko/tablewriter"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/urfave/cli/v2"
)

var CommandPruneImages = &cli.Command{
	Name:   "prune",
	Usage:  "Remove unused images",
	Action: pruneImages,
}

func pruneImages(clicontext *cli.Context) error {

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
