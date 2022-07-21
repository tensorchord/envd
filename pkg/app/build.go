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
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/builder"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/home"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandBuild = &cli.Command{
	Name:    "build",
	Aliases: []string{"b"},
	Usage:   "Build the envd environment",
	Description: `
To build an image using build.envd:
	$ envd build
To build and push the image to a registry:
	$ envd build --output type=image,name=docker.io/username/image,push=true
`,
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
		&cli.StringFlag{
			Name:    "output",
			Usage:   "Output destination (e.g. type=tar,dest=path,push=true)",
			Aliases: []string{"o"},
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force rebuild the image",
			Value: false,
		},
	},

	Action: build,
}

func build(clicontext *cli.Context) error {
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build context")
	}

	fileName, funcName, err := builder.ParseFromStr(clicontext.String("from"))
	if err != nil {
		return err
	}
	manifest, err := filepath.Abs(filepath.Join(buildContext, fileName))
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path of the build file")
	}
	if manifest == "" {
		return errors.New("file does not exist")
	}

	cfg := home.GetManager().ConfigFile()

	tag := clicontext.String("tag")
	if tag == "" {
		logrus.Debug("tag not specified, using default")
		tag = fmt.Sprintf("%s:%s", fileutil.Base(buildContext), "dev")
	}
	tag, err = docker.NormalizeNamed(tag)
	if err != nil {
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		"build-context":         buildContext,
		"build-file":            manifest,
		"config":                cfg,
		"tag":                   tag,
		flag.FlagBuildkitdImage: viper.GetString(flag.FlagBuildkitdImage),
	})
	debug := clicontext.Bool("debug")
	output := clicontext.String("output")
	force := clicontext.Bool("force")

	opt := builder.Options{
		ManifestFilePath: manifest,
		ConfigFilePath:   cfg,
		BuildFuncName:    funcName,
		BuildContextDir:  buildContext,
		Tag:              tag,
		OutputOpts:       output,
		PubKeyPath:       clicontext.Path("public-key"),
		ProgressMode:     "auto",
	}
	if debug {
		opt.ProgressMode = "plain"
	}

	logger.WithFields(logrus.Fields{
		"builder-options": opt,
	}).Debug("starting build command")

	builder, err := builder.New(clicontext.Context, opt)
	if err != nil {
		return errors.Wrap(err, "failed to create the builder")
	}
	return builder.Build(clicontext.Context, force)
}
