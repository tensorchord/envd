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
	"os"
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
			Value:    "python",
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

func NewPythonEnv(dir string) (*pythonEnv, error) {
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
			condaEnv = relPath
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pythonEnv{
		pythonVersion: "python", // use the default one
		requirements:  requirements,
		condaEnv:      condaEnv,
		indent:        "    ",
		notebook:      false,
	}, nil
}

func (pe *pythonEnv) generate() []byte {
	var buf bytes.Buffer
	buf.WriteString("def build():\n")
	buf.WriteString(fmt.Sprintf("%sbase(os=\"ubuntu20.04\", language=\"%s\")\n", pe.indent, pe.pythonVersion))
	if len(pe.requirements) > 0 {
		buf.WriteString(fmt.Sprintf("%sinstall.python_packages(requirements=\"%s\")\n", pe.indent, pe.requirements))
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

func initPythonEnv(dir string) ([]byte, error) {
	env, err := NewPythonEnv(dir)
	if err != nil {
		return nil, err
	}
	return env.generate(), nil
}

func initCommand(clicontext *cli.Context) error {
	lang := strings.ToLower(clicontext.String("lang"))
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	force := clicontext.Bool("force")
	if err != nil {
		return err
	}
	if !isValidLang(lang) {
		return errors.Errorf("invalid language (%s)", lang)
	}

	filePath := filepath.Join(buildContext, "build.envd")
	exists, err := fileutil.FileExists(filePath)
	if err != nil {
		return err
	}
	if exists && !force {
		return errors.Errorf("build.envd already exists, use --force to overwrite it")
	}

	var buildEnvdContent []byte
	if lang == "python" {
		buildEnvdContent, err = initPythonEnv(buildContext)
	} else {
		buildEnvdContent, err = templatef.ReadFile("template/" + lang + ".envd")
	}
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, buildEnvdContent, 0644)
	if err != nil {
		return errors.Wrapf(err, "Failed to create build.envd")
	}
	return nil
}
