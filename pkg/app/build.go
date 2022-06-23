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
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/builder"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/home"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandBuild = &cli.Command{
	Name:    "build",
	Aliases: []string{"b"},
	Usage:   "build envd environment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "tag",
			Usage:       "Name and optionally a tag in the 'name:tag' format",
			Aliases:     []string{"t"},
			DefaultText: "PROJECT:dev",
		},
		&cli.PathFlag{
			Name:    "from",
			Usage:   "Function to execute, format `file:func`",
			Aliases: []string{"f"},
			Value:   "build.envd:build",
		},
		&cli.PathFlag{
			Name:    "path",
			Usage:   "Path to the directory containing the build.envd",
			Aliases: []string{"p"},
			Value:   ".",
		},
		&cli.PathFlag{
			Name:    "public-key",
			Usage:   "Path to the public key",
			Aliases: []string{"pubk"},
			Value:   sshconfig.GetPublicKey(),
			Hidden:  true,
		},
		&cli.PathFlag{
			Name:    "output",
			Usage:   "Output destination (format: type=tar,dest=path)",
			Aliases: []string{"o"},
			Value:   "",
		},
	},

	Action: build,
}

func build(clicontext *cli.Context) error {
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build context")
	}

	filename, funcname, err := builder.ParseFromStr(clicontext.String("from"))
	if err != nil {
		return err
	}
	manifest, err := filepath.Abs(filepath.Join(buildContext, filename))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build file")
	}
	if manifest == "" {
		return errors.New("file does not exist")
	}

	config := home.GetManager().ConfigFile()

	tag := clicontext.String("tag")
	if tag == "" {
		logrus.Debug("tag not specified, using default")
		tag = fileutil.Base(buildContext)
	}

	logger := logrus.WithFields(logrus.Fields{
		"build-context":             buildContext,
		"build-file":                manifest,
		"config":                    config,
		"tag":                       tag,
		flag.FlagBuildkitdImage:     viper.GetString(flag.FlagBuildkitdImage),
		flag.FlagBuildkitdContainer: viper.GetString(flag.FlagBuildkitdContainer),
	})
	logger.Debug("starting build command")
	debug := clicontext.Bool("debug")
	builder, err := builder.New(clicontext.Context, config, manifest, funcname, buildContext, tag, clicontext.Path("output"), debug)
	if err != nil {
		return errors.Wrap(err, "failed to create the builder")
	}
	return builder.Build(clicontext.Context, clicontext.Path("public-key"))
}
