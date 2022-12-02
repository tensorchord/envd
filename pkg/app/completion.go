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
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	ac "github.com/tensorchord/envd/pkg/autocomplete"
)

var CommandCompletion = &cli.Command{
	Name:     "completion",
	Category: CategoryManagement,
	Usage:    "Install shell completion scripts for envd",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "shell",
			Usage:   "Shell type to install completion",
			Aliases: []string{"s"},
		},
	},

	Action: completion,
}

func completion(clicontext *cli.Context) error {
	shellList := clicontext.StringSlice("shell")

	n := len(shellList)
	if n == 0 {
		return errors.Errorf("at least one specified shell type")
	}

	for i := 0; i < n; i++ {
		logrus.Infof("[%d/%d] Add completion %s", i+1, n, shellList[i])
		switch shellList[i] {
		case "zsh":
			if err := ac.InsertZSHCompleteEntry(); err != nil {
				return err
			}
		case "bash":
			if err := ac.InsertBashCompleteEntry(); err != nil {
				return err
			}
		default:
			return errors.Errorf("unknown shell type %s (support type: {bash|zsh})", shellList[i])
		}
	}

	return nil
}
