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
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/app/formatter/json"
	"github.com/tensorchord/envd/pkg/app/formatter/table"
)

var CommandVersion = &cli.Command{
	Name:     "version",
	Category: CategoryOther,
	Aliases:  []string{"v"},
	Usage:    "Print envd version information",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "short",
			Usage:   "Only print the version number",
			Value:   false,
			Aliases: []string{"s"},
		},
		&cli.BoolFlag{
			Name:    "detail",
			Usage:   "Print details about the envd environment",
			Value:   false,
			Aliases: []string{"d"},
		},
		&formatter.FormatFlag,
	},
	Action: outputVersion,
}

func outputVersion(clicontext *cli.Context) error {
	format := clicontext.String("format")
	switch format {
	case "table":
		return table.PrintVersion(clicontext)
	case "json":
		return json.PrintVersion(clicontext)
	}
	return nil
}
