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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cockroachdb/errors"
	cli "github.com/urfave/cli/v2"
)

var CommandDemo = &cli.Command{
	Name:     "demo",
	Category: CategoryManagement,
	Usage:    "Interactively initializes the current directory with the build.envd file",
	Action:   demoCommand,
}

func demoCommand(clicontext *cli.Context) error {

	startQuestion(LanguageInput())
	if selectionMap[LabelLanguage][0] == "python" {
		startQuestion(PythonPackageInput())
	} else if selectionMap[LabelLanguage][0] == "r" {
		startQuestion(RPackageChoice())
	}
	startQuestion(CudaChoice())
	if selectionMap["Cuda"][0] == "Yes" {
		startQuestion(CudaVersionChoice())
	}

	err := generateFile(clicontext)
	if err != nil {
		return errors.Wrap(err, "error generating build.envd")
	}

	fmt.Println("Successfully generated build.envd file!")
	return nil
}

func startQuestion(input input) tea.Model {
	p := tea.NewProgram(InitModel(input))
	m, err := p.Run()
	if err != nil {
		fmt.Printf("There was an error generating build.envd: %v", err)
		os.Exit(1)
	}
	return m
}

var selectionMap = make(map[string][]string)
var itemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))
var selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4dff4d"))

const (
	SINGLE_SELECT   string = "single select"
	MULTIPLE_SELECT string = "multiple select"
	INPUT           string = "input"
)

const (
	LabelLanguage      string = "Language"
	LabelCuda          string = "Cuda"
	LabelPythonPackage string = "Python Package"
	LabelRPackage      string = "R Package"
)

type model struct {
	step     int
	cursor   int
	selected map[int]struct{}
	input    input
}

type input struct {
	prompt    string
	inputType string
	label     string
	options   []string
}

func LanguageInput() input {
	return input{
		prompt:    "Choose a programming language",
		inputType: SINGLE_SELECT,
		label:     LabelLanguage,
		options: []string{
			"python",
			"r",
			"julia",
		},
	}
}

func PythonPackageInput() input {
	return input{
		prompt:    "Choose your python packages",
		inputType: MULTIPLE_SELECT,
		label:     LabelPythonPackage,
		options: []string{
			"numpy",
			"tensorflow",
		},
	}
}

func RPackageChoice() input {
	return input{
		prompt:    "Choose your r packages",
		inputType: MULTIPLE_SELECT,
		label:     LabelRPackage,
		options: []string{
			"remotes",
			"rlang",
		},
	}
}

func CudaChoice() input {
	return input{
		prompt:    "Include Cuda?",
		inputType: SINGLE_SELECT,
		label:     "Cuda",
		options: []string{
			"Yes",
			"No",
		},
	}
}

func CudaVersionChoice() input {
	return input{
		prompt:    "Choose a cuda version",
		inputType: SINGLE_SELECT,
		label:     LabelCuda,
		options: []string{
			"11.0",
			"10.2",
			"10.1",
		},
	}
}

func InitModel(input input) model {
	return model{
		input:    input,
		step:     0,
		selected: make(map[int]struct{}),
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

	return s + "\n"
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
			os.Exit(0)
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

func generateFile(clicontext *cli.Context) error {
	var buf bytes.Buffer
	buf.WriteString("def build():\n")
	buf.WriteString(fmt.Sprintf("    base(os=\"ubuntu20.04\", language=\"%s\")\n", selectionMap[LabelLanguage][0]))
	buf.WriteString(generatePackagesStr("python", selectionMap[LabelPythonPackage]))
	buf.WriteString(generatePackagesStr("r", selectionMap[LabelRPackage]))
	if selectionMap[LabelCuda][0] == "Yes" {
		buf.WriteString(fmt.Sprintf("    cuda(version=\"%s\", cudann=\"8\")\n", selectionMap[LabelCuda][0]))
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

func generatePackagesStr(packageName string, packages []string) string {
	if len(packages) == 0 {
		return ""
	}
	s := fmt.Sprintf("    install.%s_packages(name = [\n", packageName)
	for i, p := range packages {
		s += fmt.Sprintf("        \"%s\"", p)
		if i != len(packages)-1 {
			s += ", "
		}
		s += "\n"
	}
	s += "    ])\n"
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
