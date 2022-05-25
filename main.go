package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func main() {
	data := [][]string{
		[]string{"mnist", "http://localhost:8888", "34c5e99ec315", "97%", "Running"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "jupyter", "GPU util", "Status"})

	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}
