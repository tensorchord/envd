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

package main

import (
	"path/filepath"

	"github.com/cockroachdb/errors"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/builder"
	"github.com/tensorchord/envd/pkg/home"
)

var CommandBuild = &cli.Command{
	Name:    "build",
	Aliases: []string{"b"},
	Usage:   "build envd environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "tag",
			Usage:   "Name and optionally a tag in the 'name:tag' format",
			Aliases: []string{"t"},
			Value:   "envd:dev",
		},
		&cli.PathFlag{
			Name:    "file",
			Usage:   "Name of the build.envd (Default is 'PATH/build.envd')",
			Aliases: []string{"f"},
			Value:   "./build.envd",
		},
	},

	Action: build,
}

func build(clicontext *cli.Context) error {
	path, err := filepath.Abs(clicontext.Path("file"))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build file")
	}
	if path == "" {
		return errors.New("file does not exist")
	}

	config := home.GetManager().ConfigFile()

	tag := clicontext.String("tag")

	builder, err := builder.New(clicontext.Context, config, path, tag)
	if err != nil {
		return errors.Wrap(err, "failed to create the builder")
	}
	return builder.Build(clicontext.Context)
}
