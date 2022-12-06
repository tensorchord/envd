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
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

//go:embed template
var templatef embed.FS

var CommandInit = &cli.Command{
	Name:     "init",
	Category: CategoryManagement,
	Aliases:  []string{"i"},
	Usage:    "Initializes the current directory with the build.envd file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "lang",
			Usage:    "language usage. Support Python, R, Julia",
			Aliases:  []string{"l"},
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "force",
			Usage:    "overwrite the build.envd if existed",
			Aliases:  []string{"f"},
			Required: false,
		},
		&cli.PathFlag{
			Name:    "path",
			Usage:   "Path to the directory containing the build.envd",
			Aliases: []string{"p"},
			Value:   ".",
		},
	},
	Action: initCommand,
}

func isValidLang(lang string) bool {
	switch lang {
	case
		"python",
		"r",
		"julia":
		return true
	}
	return false
}

type pythonEnv struct {
	pythonVersion string
	requirements  string
	condaEnv      string
	indent        string
	notebook      bool
}

func InitPythonEnv(dir string) error {
	requirements := ""
	condaEnv := ""
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if d.Name() == "requirements.txt" && len(requirements) <= 0 {
			requirements = relPath
			return nil
		}
		if isCondaEnvFile(d.Name()) && len(condaEnv) <= 0 {
			selectionMap[LabelCondaEnv] = []string{relPath}
		}
		return nil
	})
	if err != nil {
		return err
	}

	selectionMap[LabelPythonRequirement] = []string{requirements}
	if len(requirements) == 0 {
		startQuestion(PythonPackageChoice)
	}
	startQuestion(JupyterChoice)
	return nil
}

func (pe *pythonEnv) generate() []byte {
	var buf bytes.Buffer
	buf.WriteString("def build():\n")
	buf.WriteString(fmt.Sprintf("%sbase(os=\"ubuntu20.04\", language=\"%s\")\n", pe.indent, pe.pythonVersion))
	if len(pe.requirements) > 0 {
		buf.WriteString(fmt.Sprintf("%sinstall.python_packages(requirements=\"%s\")\n", pe.indent, pe.requirements))
	} else {
		buf.WriteString(generatePackagesStr("python", selectionMap[LabelPythonPackage]))
	}
	if len(pe.condaEnv) > 0 {
		buf.WriteString(fmt.Sprintf("%sinstall.conda_packages(env_file=\"%s\")\n", pe.indent, pe.condaEnv))
	}
	if pe.notebook {
		buf.WriteString(fmt.Sprintf("%sconfig.jupyter()\n", pe.indent))
	}
	return buf.Bytes()
}

// naive check
func isCondaEnvFile(file string) bool {
	switch file {
	case
		"environment.yml",
		"environment.yaml",
		"env.yml",
		"env.yaml":
		return true
	}
	return false
}

func initCommand(clicontext *cli.Context) error {
	lang := strings.ToLower(clicontext.String("lang"))
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	force := clicontext.Bool("force")
	if err != nil {
		return err
	}

	if !isValidLang(lang) {
		startQuestion(LanguageChoice)
		lang = selectionMap[LabelLanguage][0]
	}

	filePath := filepath.Join(buildContext, "build.envd")
	exists, err := fileutil.FileExists(filePath)
	if err != nil {
		return err
	}
	if exists && !force {
		return errors.Errorf("build.envd already exists, use --force to overwrite it")
	}

	if lang == "python" {
		err = InitPythonEnv(buildContext)
		if err != nil {
			return err
		}
	} else if lang == "r" {
		startQuestion(RPackageChoice)
	}

	startQuestion(CudaChoice)
	if selectionMap[LabelCudaChoice][0] == "Yes" {
		startQuestion(CudaVersionChoice)
	}

	generateFile(clicontext)

	return nil
}
