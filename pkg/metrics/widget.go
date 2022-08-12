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

package metrics

import (
	"fmt"

	"github.com/bcicen/ctop/cwidgets"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	defaultRowHeight int = 3
)

type WidgetCol struct {
	widget ui.Drawable
	Width  int
	Height int
}

type WidgetRow struct {
	ui.Block
	Cols   []*WidgetCol
	X      int
	Y      int
	Width  int
	Height int
}

func NewWidgetRow(nRow int) *WidgetRow {
	return &WidgetRow{
		Cols:   make([]*WidgetCol, 0),
		X:      0,
		Y:      nRow * defaultRowHeight,
		Width:  0,
		Height: defaultRowHeight,
	}
}

func (row *WidgetRow) Add(col *WidgetCol) {
	nx := row.X + row.Width
	col.widget.SetRect(nx, row.Y, nx+col.Width, row.Y+row.Height)
	row.Width = row.Width + col.Width
	row.SetRect(row.X, row.Y, row.X+row.Width, row.Y+row.Height)
	row.Cols = append(row.Cols, col)
}

func (row *WidgetRow) Draw(buf *ui.Buffer) {
	for item := range row.Cols {
		row.Cols[item].widget.Draw(buf)
	}
}

func NewNameCol(name string) *WidgetCol {
	w := widgets.NewParagraph()
	w.Text = name
	w.Border = false
	w.TextStyle.Fg = ui.ColorClear
	return &WidgetCol{
		widget: w,
		Height: defaultRowHeight,
		Width:  20,
	}
}

func NewCPUCol(metChan <-chan Metrics) *WidgetCol {
	w := widgets.NewGauge()
	w.Percent = 0
	w.Border = false
	go func() {
		for {
			ms := <-metChan
			val := ms.CPUUtil
			w.BarColor = colorScale(val)
			w.Label = fmt.Sprintf("%d%%", val)
			w.LabelStyle.Fg = ui.ColorClear
			if val > 100 {
				val = 100
			}
			w.Percent = val
		}
	}()
	return &WidgetCol{
		widget: w,
		Height: defaultRowHeight,
		Width:  20,
	}
}

func NewMEMCol(metChan <-chan Metrics) *WidgetCol {
	w := widgets.NewGauge()
	w.Percent = 0
	w.Border = false
	go func() {
		for {
			ms := <-metChan
			mPercent := ms.MemPercent
			w.BarColor = colorScale(mPercent)
			w.LabelStyle.Fg = ui.ColorClear
			w.Label = fmt.Sprintf("%s / %s", cwidgets.ByteFormat64Short(ms.MemUsage), cwidgets.ByteFormat64Short(ms.MemLimit))
			w.Percent = mPercent
		}
	}()
	return &WidgetCol{
		widget: w,
		Height: defaultRowHeight,
		Width:  20,
	}
}

func colorScale(n int) ui.Color {
	if n < 50 {
		return ui.ColorGreen
	} else if n < 80 {
		return ui.ColorYellow
	} else {
		return ui.ColorRed
	}
}
