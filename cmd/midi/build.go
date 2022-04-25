// Copyright 2022 The MIDI Authors
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
	"github.com/spf13/viper"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/MIDI/pkg/builder"
	"github.com/tensorchord/MIDI/pkg/flag"
)

var CommandBuild = &cli.Command{
	Name:      "build",
	Aliases:   []string{"b"},
	Usage:     "build MIDI environment",
	UsageText: `TODO`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "tag",
			Usage:   "Name and optionally a tag in the 'name:tag' format",
			Aliases: []string{"t"},
			Value:   "midi:dev",
		},
		&cli.PathFlag{
			Name:    "file",
			Usage:   "Name of the build.MIDI (Default is 'PATH/build.MIDI')",
			Aliases: []string{"f"},
			Value:   "./build.MIDI",
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

	config, err := filepath.Abs(viper.GetString(flag.FlagConfig))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the config file")
	}
	if config == "" {
		return errors.New("file does not exist")
	}

	tag := clicontext.String("tag")

	builder := builder.New("unix:///run/buildkit/buildkitd.sock", config, path, tag)
	return builder.Build(clicontext.Context)
}
