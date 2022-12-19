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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cockroachdb/errors"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
)

var selectionMap = make(map[string][]string)
var itemStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"})
var selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Dark: "14", Light: "6"})
var indentation = "    "

const (
	SINGLE_SELECT   string = "single select"
	MULTIPLE_SELECT string = "multiple select"
	INPUT           string = "input"
)

const (
	LabelLanguage          string = "Language"
	LabelCuda              string = "Cuda"
	LabelCudaChoice        string = "Cuda Choice"
	LabelJupyterChoice     string = "Jupyter Choice"
	LabelPythonPackage     string = "Python Package"
	LabelRPackage          string = "R Package"
	LabelCondaEnv          string = "Conda Env"
	LabelPythonRequirement string = "Python Requirement"
)

type model struct {
	step     int
	cursor   int
	selected map[int]struct{}
	input    input
	exit     bool
}

type input struct {
	prompt    string
	inputType string
	label     string
	options   []string
}

var LanguageChoice = input{
	prompt:    "Choose a programming language",
	inputType: SINGLE_SELECT,
	label:     LabelLanguage,
	options: []string{
		"python",
		"r",
		"julia",
	},
}

var PythonPackageChoice = input{
	prompt:    "Choose your python packages",
	inputType: MULTIPLE_SELECT,
	label:     LabelPythonPackage,
	options: []string{
		"numpy",
		"tensorflow",
	},
}

var RPackageChoice = input{
	prompt:    "Choose your r packages",
	inputType: MULTIPLE_SELECT,
	label:     LabelRPackage,
	options: []string{
		"remotes",
		"rlang",
	},
}

var CudaChoice = input{
	prompt:    "Include Cuda?",
	inputType: SINGLE_SELECT,
	label:     LabelCudaChoice,
	options: []string{
		"Yes",
		"No",
	},
}

var CudaVersionChoice = input{
	prompt:    "Choose a cuda version",
	inputType: SINGLE_SELECT,
	label:     LabelCuda,
	options: []string{
		"11.6.2",
		"11.3.1",
		"11.2.2",
	},
}

var JupyterChoice = input{
	prompt:    "Include Jupyter?",
	inputType: SINGLE_SELECT,
	label:     LabelJupyterChoice,
	options: []string{
		"Yes",
		"No",
	},
}

func InitModel(input input) model {
	return model{
		input:    input,
		step:     0,
		selected: make(map[int]struct{}),
		exit:     false,
	}
}

func (m model) View() string {
	s := m.input.prompt

	switch m.input.inputType {
	case SINGLE_SELECT, MULTIPLE_SELECT:
		s += m.renderMultipleChoice()

	case INPUT:
		// TODO: implement input if needed

	}

	s += "\nPress q to quit. "
	if m.input.inputType == MULTIPLE_SELECT {
		s += "Press space to select"
	}

	return s
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	choices := m.input.options

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(choices)-1 {
				m.cursor++
			}
		case " ":
			if m.input.inputType == MULTIPLE_SELECT {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		case "enter":
			if m.input.inputType == SINGLE_SELECT {
				m.selected[m.cursor] = struct{}{}
			}
			selectionMap[m.input.label] = []string{choices[m.cursor]}
			m.addSelection()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) addSelection() {
	selectionMap[m.input.label] = []string{}
	for i := range m.selected {
		selectionMap[m.input.label] = append(selectionMap[m.input.label], m.input.options[i])
	}
}

func startQuestion(input input) {
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	if !isTerminal {
		return
	}

	p := tea.NewProgram(InitModel(input))
	m, err := p.Run()
	if m.(model).exit {
		os.Exit(0)
	}
	if err != nil {
		fmt.Printf("There was an error generating build.envd: %v", err)
		os.Exit(1)
	}
}
func generateFile(clicontext *cli.Context) error {
	var buf bytes.Buffer
	buf.WriteString("def build():\n")
	buf.WriteString(fmt.Sprintf("%sbase(os=\"ubuntu20.04\", language=\"%s\")\n", indentation, selectionMap[LabelLanguage][0]))
	buf.WriteString(generatePackagesStr("python", selectionMap[LabelPythonPackage]))
	buf.WriteString(generatePackagesStr("r", selectionMap[LabelRPackage]))
	if len(selectionMap[LabelPythonRequirement]) > 0 {
		buf.WriteString(fmt.Sprintf("%sinstall.python_packages(requirements=\"%s\")\n", indentation, selectionMap[LabelPythonRequirement][0]))
	}
	if len(selectionMap[LabelCondaEnv]) > 0 {
		buf.WriteString(fmt.Sprintf("%sinstall.conda_packages(env_file=\"%s\")\n", indentation, selectionMap[LabelCondaEnv][0]))
	}
	if len(selectionMap[LabelCudaChoice]) > 0 && selectionMap[LabelCudaChoice][0] == "Yes" {
		buf.WriteString(fmt.Sprintf("%scuda(version=\"%s\", cudann=\"8\")\n", indentation, selectionMap[LabelCuda][0]))
	}
	if len(selectionMap[LabelJupyterChoice]) > 0 && selectionMap[LabelJupyterChoice][0] == "Yes" {
		buf.WriteString(fmt.Sprintf("%sconfig.jupyter()\n", indentation))
	}

	buildEnvdContent := buf.Bytes()
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return err
	}
	filePath := filepath.Join(buildContext, "build.envd")
	err = os.WriteFile(filePath, buildEnvdContent, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to create build.envd")
	}
	return nil
}

func generatePackagesStr(name string, packages []string) string {
	if len(packages) == 0 {
		return ""
	}
	s := fmt.Sprintf("%sinstall.%s_packages(name = [\n", indentation, name)
	for i, p := range packages {
		s += fmt.Sprintf("%s\"%s\"", strings.Repeat(indentation, 2), p)
		if i != len(packages)-1 {
			s += ", "
		}
		s += "\n"
	}
	s += fmt.Sprintf("%s])\n", indentation)
	return s
}

func (m model) renderMultipleChoice() string {
	s := "\n\n"
	for i, choice := range m.input.options {
		cursor := " "
		style := itemStyle
		if m.cursor == i {
			cursor = ">"
			style = selectedItemStyle
		}

		if m.input.inputType == MULTIPLE_SELECT {
			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			s += style.Render((fmt.Sprintf("%s [%s] %s", cursor, checked, choice))) + "\n"
		}

		if m.input.inputType == SINGLE_SELECT {
			s += style.Render(fmt.Sprintf("%s %s", cursor, choice)) + "\n"
		}
	}
	return s
}
