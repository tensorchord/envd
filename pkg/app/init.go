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
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandInit = &cli.Command{
	Name:     "init",
	Category: CategoryBasic,
	Aliases:  []string{"i"},
	Usage:    "Automatically generate the build.envd",
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

	if len(requirements) == 0 {
		startQuestion(PythonPackageChoice)
	} else {
		selectionMap[LabelPythonRequirement] = []string{requirements}
	}
	startQuestion(JupyterChoice)
	return nil
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
		if len(selectionMap[LabelLanguage]) > 0 {
			lang = selectionMap[LabelLanguage][0]
		} else {
			lang = "python"
			selectionMap[LabelLanguage] = []string{lang}
		}
	} else {
		selectionMap[LabelLanguage] = []string{lang}
	}
	defer func(start time.Time) {
		telemetry.GetReporter().Telemetry(
			"init", telemetry.AddField("duration", time.Since(start).Seconds()))
	}(time.Now())

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
	if len(selectionMap[LabelCudaChoice]) > 0 && selectionMap[LabelCudaChoice][0] == "Yes" {
		startQuestion(CudaVersionChoice)
	}

	err = generateFile(clicontext)
	if err != nil {
		return errors.Wrapf(err, "Failed to create build.envd")
	}

	return nil
}
