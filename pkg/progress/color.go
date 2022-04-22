// Copyright 2022 The MIDI Authors
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

package progress

import "github.com/fatih/color"

var noColor = makeNoColor()
var cachedColor = makeColor(color.FgHiGreen)
var successColor = makeColor(color.FgHiGreen)

var availablePrefixColors = []*color.Color{
	makeColor(color.FgBlue),
	makeColor(color.FgMagenta),
	makeColor(color.FgCyan),
	makeColor(color.FgRed),
	makeColor(color.FgYellow),
	makeColor(color.FgGreen),
	makeColor(color.FgHiBlue),
	makeColor(color.FgHiMagenta),
	makeColor(color.FgHiCyan),
	makeColor(color.FgHiRed),
	makeColor(color.FgHiYellow),
	makeColor(color.FgHiWhite),
}

func makeColor(att color.Attribute) *color.Color {
	c := color.New()
	c.Add(att)
	return c
}

func makeNoColor() *color.Color {
	c := color.New()
	c.DisableColor()
	return c
}
