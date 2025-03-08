// Copyright 2025 The envd Authors
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
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var (
	//go:embed template/uv.envd
	templateUV string
	//go:embed template/conda.envd
	templateConda string
	//go:embed template/torch.envd
	templateTorch string

	templates = map[string]string{
		"uv":    templateUV,
		"conda": templateConda,
		"torch": templateTorch,
	}
)

func joinKeysToString(table map[string]string) string {
	keys := make([]string, 0, len(table))
	for k := range table {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func isDefaultTemplate(name string) bool {
	_, ok := templates[name]
	return ok
}

var CommandNew = &cli.Command{
	Name:     "new",
	Category: CategoryBasic,
	Aliases:  []string{"n"},
	Usage:    "Create a new `build.envd` file from pre-defined templates",
	Description: `The template used by this command is stored in the 
		'$HOME/.config/envd/templates' directory, we provide some pre-defined templates for you
		to use during 'envd bootstrap', you can also add your own templates to this directory.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "template",
			Usage:    fmt.Sprintf("Template name to use (`envd bootstrap` will add [%s])", joinKeysToString(templates)),
			Aliases:  []string{"t"},
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "force",
			Usage:    "Overwrite the build.envd if existed",
			Aliases:  []string{"f"},
			Required: false,
		},
		&cli.PathFlag{
			Name:    "path",
			Usage:   "Path to the directory of the build.envd",
			Aliases: []string{"p"},
			Value:   ".",
		},
	},
	Action: newCommand,
}

func newCommand(clicontext *cli.Context) error {
	workDir, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path")
	}

	force := clicontext.Bool("force")
	filePath := filepath.Join(workDir, "build.envd")
	exists, err := fileutil.FileExists(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to check file exists")
	}
	if exists && !force {
		return errors.New("build.envd already exists, use `--force` to overwrite")
	}

	template := clicontext.String("template")
	templateFile := fmt.Sprintf("%s.envd", template)
	templatePath, err := fileutil.TemplateFile(templateFile)
	if err != nil {
		return errors.Wrapf(err, "failed to get template file: `%s`", templateFile)
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		if os.IsNotExist(err) && isDefaultTemplate(template) {
			// Add default templates to the template directory if not exist
			err = addTemplates(clicontext)
			if err != nil {
				return err
			}
			content, err = os.ReadFile(templatePath)
			if err != nil {
				return err
			}
		} else {
			return errors.Wrapf(err, "failed to read the template file `%s`", templatePath)
		}
	}
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write the build.envd file")
	}
	logrus.Infof("Template `%s` is created in `%s`", template, filePath)

	return nil
}

func addTemplates(clicontext *cli.Context) error {
	for name, content := range templates {
		file, err := fileutil.TemplateFile(name + ".envd")
		if err != nil {
			return errors.Wrapf(err, "failed to get template file path: %s", name)
		}
		exist, err := fileutil.FileExists(file)
		if err != nil {
			return errors.Wrapf(err, "failed to check file exists: %s", file)
		}
		if exist {
			logrus.Debugf("Template file `%s` already exists in `%s`", name, file)
			continue
		}
		err = os.WriteFile(file, []byte(content), 0644)
		if err != nil {
			return errors.Wrapf(err, "failed to write template file: %s", name)
		}
		logrus.Debugf("Template file `%s` is added to `%s`", name, file)
	}

	return nil
}
