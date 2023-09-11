// Copyright 2023 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"
)

var CommandReference = &cli.Command{
	Name:     "reference",
	Category: CategoryOther,
	Usage:    "Print envd reference documentation",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "output",
			Usage: "Output file, if not specified, print to stderr",
			Value: "",
		},
	},
	Action: outputReference,
}

const referenceHeader = `
# envd CLI Reference

This is a reference for the CLI commands of envd.

::: tip
The documentation is auto-generated from [envd app](https://github.com/tensorchord/envd/blob/main/pkg/app/app.go), please do not edit it manually.
:::

`

func outputReference(clicontext *cli.Context) error {
	doc, err := clicontext.App.ToMarkdown()
	if err != nil {
		return errors.Wrap(err, "failed to generate the markdown document")
	}
	content := referenceHeader + doc
	output := clicontext.String("output")
	if len(output) > 0 {
		return os.WriteFile(output, []byte(content), 0644)
	}
	fmt.Println(content)
	return nil
}
