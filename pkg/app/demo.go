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
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if _, err := os.Stat(buildContext + "/build.envd"); !errors.Is(err, os.ErrNotExist) {
		return errors.New("build.envd already exists")
	}

	startQuestion(LanguageChoice)
	languageChoice := selectionMap[LabelLanguage][0]
	if languageChoice == "python" {
		startQuestion(PythonPackageChoice)
		startQuestion(JupyterChoice)
	} else if languageChoice == "r" {
		startQuestion(RPackageChoice)
	}
	startQuestion(CudaChoice)
	if selectionMap[LabelCudaChoice][0] == "Yes" {
		startQuestion(CudaVersionChoice)
	}

	err = generateFile(clicontext)
	if err != nil {
		return errors.Wrap(err, "error generating build.envd")
	}

	fmt.Println("Successfully generated build.envd file!")
	return nil
}

func startQuestion(input input) {
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
