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
	"embed"
	"io/ioutil"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/tensorchord/envd/pkg/util/fileutil"
	cli "github.com/urfave/cli/v2"
)

//go:embed template
var templatef embed.FS

var CommandInit = &cli.Command{
	Name:    "init",
	Aliases: []string{"i"},
	Usage:   "Initializes the current directory with the build.envd file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "lang",
			Usage:    "language usage. Support Python, R",
			Aliases:  []string{"l"},
			Required: true,
		},
	},
	Action: initCommand,
}

func isValidLang(lang string) bool {
	switch lang {
	case
		"python",
		"r":
		return true
	}
	return false
}

func initCommand(clicontext *cli.Context) error {
	lang := strings.ToLower(clicontext.String("lang"))
	if !isValidLang(lang) {
		return errors.Errorf("invalid language %s", lang)
	}

	exists, err := fileutil.FileExists("build.envd")
	if err != nil {
		return err
	}
	if exists {
		return errors.Errorf("build.envd already exists")
	}

	buildEnvdContent, err := templatef.ReadFile("template/" + lang + ".envd")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("build.envd", buildEnvdContent, 0644)
	if err != nil {
		return errors.Wrapf(err, "Failed to create build.envd")
	}
	return nil
}
