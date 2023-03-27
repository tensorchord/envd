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

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	ac "github.com/tensorchord/envd/pkg/autocomplete"
)

var CommandCompletion = &cli.Command{
	Name:     "completion",
	Category: CategorySettings,
	Usage:    "Install shell completion scripts for envd",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "shell",
			Usage:   "Shell type to install completion",
			Aliases: []string{"s"},
		},
		&cli.BoolFlag{
			Name:  "no-install",
			Usage: "Only output the completion script and don't install it",
		},
	},

	Action: completion,
}

func handleCompletion(clicontext *cli.Context, installFunc func(*cli.Context) error, outputFunc func(*cli.Context) (string, error)) error {
	if clicontext.Bool("no-install") {
		script, err := outputFunc(clicontext)
		if err != nil {
			return err
		}
		fmt.Println(script)
	} else {
		if err := installFunc(clicontext); err != nil {
			return err
		}
	}
	return nil
}

func completion(clicontext *cli.Context) error {
	shellList := clicontext.StringSlice("shell")

	n := len(shellList)
	if n == 0 {
		defaultShell := os.Getenv("SHELL")
		if defaultShell != "" {
			shellList = append(shellList, filepath.Base(defaultShell))
			n++
		} else {
			return errors.Errorf("Can't detect the default shell, please specify at least one shell type with --shell")
		}
	}

	for i := 0; i < n; i++ {
		logrus.Infof("[%d/%d] Add completion %s", i+1, n, shellList[i])
		switch shellList[i] {
		case "zsh":
			if err := handleCompletion(clicontext, ac.InsertZSHCompleteEntry, ac.ZshCompleteEntry); err != nil {
				return err
			}
		case "bash":
			if err := handleCompletion(clicontext, ac.InsertBashCompleteEntry, ac.BashCompleteEntry); err != nil {
				return err
			}
		case "fish":
			if err := handleCompletion(clicontext, ac.InsertFishCompleteEntry, ac.FishCompleteEntry); err != nil {
				return err
			}
		default:
			return errors.Errorf("unknown shell type %s (support type: {bash|zsh})", shellList[i])
		}
	}

	return nil
}
