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

	p := tea.NewProgram(InitModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("There was an error generating build.envd: %v", err)
		os.Exit(1)
	}
	finalModel := m.(model)
	err = generateFile(clicontext, finalModel.selections)
	if err != nil {
		fmt.Printf("There was an error generating build.envd: %v", err)
		os.Exit(1)
	}
	fmt.Println("Generated build.envd file!")
	return nil
}

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
	step         int
	inputs       []input
	currentInput *input
	cursor       int
	selected     map[int]struct{}
	prevText     string
	selections   map[string][]string
}

type input struct {
	prompt    string
	inputType string
	label     string
	options   []inputnode
	next      *input
}

type inputnode struct {
	label string
	// value string
	next *input
}

func InitChoice() []input {
	condaChoice := input{
		prompt:    "Include Conda?",
		inputType: SINGLE_SELECT,
		label:     "Conda",
		options: []inputnode{
			{
				label: "Yes",
			},
			{
				label: "No",
			},
		},
	}

	pythonPackageChoice := input{
		prompt:    "Choose your python packages",
		inputType: MULTIPLE_SELECT,
		label:     LabelPythonPackage,
		options: []inputnode{
			{
				label: "numpy",
			},
			{
				label: "tensorflow",
			},
		},
		next: &condaChoice,
	}

	RPackageChoice := input{
		prompt:    "Choose your R packages",
		inputType: MULTIPLE_SELECT,
		label:     LabelRPackage,
		options: []inputnode{
			{
				label: "remotes",
			},
			{
				label: "rlang",
			},
		},
	}

	languageChoice := input{
		prompt:    "Choose a programming language",
		inputType: SINGLE_SELECT,
		label:     LabelLanguage,
		options: []inputnode{
			{
				label: "python",
				next:  &pythonPackageChoice,
			},
			{
				label: "r",
				next:  &RPackageChoice,
			},
			{
				label: "julia",
			},
		},
	}

	cudaVersion := input{
		prompt:    "Choose a cuda version",
		inputType: SINGLE_SELECT,
		label:     LabelCuda,
		options: []inputnode{
			{
				label: "11.0",
			},
			{
				label: "10.2",
			},
			{
				label: "10.1",
			},
		},
	}

	cudaChoice := input{
		prompt:    "Include Cuda?",
		inputType: SINGLE_SELECT,
		label:     "Cuda",
		options: []inputnode{
			{
				label: "Yes",
				next:  &cudaVersion,
			},
			{
				label: "No",
			},
		},
	}

	return []input{
		languageChoice,
		cudaChoice,
	}
}

func InitModel() model {
	choices := InitChoice()
	return model{
		inputs:       choices,
		step:         0,
		selected:     make(map[int]struct{}),
		currentInput: &choices[0],
		prevText:     "",
		selections:   make(map[string][]string),
	}
}

func (m model) View() string {
	if m.currentInput == nil {
		return ""
	}

	s := m.currentInput.prompt

	switch m.currentInput.inputType {
	case SINGLE_SELECT, MULTIPLE_SELECT:
		s += m.renderMultipleChoice()

	case INPUT:
		// TODO: implement input

	}

	s += "\nPress q to quit. "
	if m.currentInput.inputType == MULTIPLE_SELECT {
		s += "Press space to select"
	}

	return fmt.Sprintf("%s%s\n", m.prevText, s)
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	choices := m.currentInput.options

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
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
			if m.currentInput.inputType == MULTIPLE_SELECT {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		case "enter":
			selected := m.getSelectedOptions()
			m.prevText += renderSelection(m.currentInput.prompt, selected)
			m = m.iterate()
			if m.currentInput == nil {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func generateFile(clicontext *cli.Context, selections map[string][]string) error {
	var buf bytes.Buffer
	buf.WriteString("def build():\n")
	// buf.WriteString(fmt.Sprintf("    base(os=\"ubuntu20.04\", language=\"%s\")\n", selections[LabelLanguage][0]))

	buildEnvdContent := buf.Bytes()
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return err
	}
	filePath := filepath.Join(buildContext, "build.envd")
	fmt.Println("File Path", filePath)
	err = os.WriteFile(filePath, buildEnvdContent, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to create build.envd")
	}
	return nil
}

func (m model) renderMultipleChoice() string {
	s := "\n\n"
	for i, choice := range m.currentInput.options {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		if m.currentInput.inputType == MULTIPLE_SELECT {
			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.label)
		} else {
			s += fmt.Sprintf("%s  %s\n", cursor, choice.label)
		}
	}
	return s
}

func (m model) iterate() model {
	m.selections[m.currentInput.label] = m.getSelectedOptions()
	currentChoice := m.currentInput.options[m.cursor]
	m.selected = make(map[int]struct{})
	if currentChoice.next == nil {
		if m.currentInput.next != nil {
			m.currentInput = m.currentInput.next
		} else {
			m.step++
			if m.step >= len(m.inputs) {
				m.currentInput = nil
				return m
			}
			m.currentInput = &m.inputs[m.step]

		}
	} else {
		m.currentInput = currentChoice.next
	}
	return m
}

func (m model) getSelectedOptions() []string {
	var selected []string
	for i := range m.selected {
		if _, ok := m.selected[i]; ok {
			selected = append(selected, m.currentInput.options[i].label)
		}
	}
	if m.currentInput.inputType == SINGLE_SELECT {
		selected = append(selected, m.currentInput.options[m.cursor].label)
	}
	return selected
}

func renderSelection(prompt string, labels []string) string {
	s := prompt + "\n"
	for _, label := range labels {
		s += label + "\n"
	}
	return s + "\n"
}
