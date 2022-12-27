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
	"os"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/app/formatter/json"
	"github.com/tensorchord/envd/pkg/app/formatter/table"
	"github.com/tensorchord/envd/pkg/home"
)

var CommandContextList = &cli.Command{
	Name:  "ls",
	Usage: "List envd contexts",
	Flags: []cli.Flag{
		&formatter.FormatFlag,
	},
	Action: contextList,
}

func contextList(clicontext *cli.Context) error {
	contexts, err := home.GetManager().ContextList()
	if err != nil {
		return errors.Wrap(err, "failed to list context")
	}
	format := clicontext.String("format")
	switch format {
	case "table":
		table.RenderContext(os.Stdout, contexts)
	case "json":
		return json.PrintContext(contexts)
	}
	return nil
}
