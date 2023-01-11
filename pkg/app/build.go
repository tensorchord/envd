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
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	buildutil "github.com/tensorchord/envd/pkg/app/build"
	"github.com/tensorchord/envd/pkg/app/telemetry"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
)

var CommandBuild = &cli.Command{
	Name:     "build",
	Category: CategoryBasic,
	Aliases:  []string{"b"},
	Usage:    "Build the envd environment",
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
		&cli.BoolFlag{
			Name:    "use-proxy",
			Usage:   "Use HTTPS_PROXY/HTTP_PROXY/NO_PROXY in the build process",
			Aliases: []string{"proxy"},
			Value:   false,
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
			Value:   sshconfig.GetPublicKeyOrPanic(),
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
		// https://github.com/urfave/cli/issues/1134#issuecomment-1191407527
		&cli.StringFlag{
			Name:    "export-cache",
			Usage:   "Export the cache (e.g. type=registry,ref=<image>)",
			Aliases: []string{"ec"},
		},
		&cli.StringFlag{
			Name:    "import-cache",
			Usage:   "Import the cache (e.g. type=registry,ref=<image>)",
			Aliases: []string{"ic"},
		},
	},
	Action: build,
}

func build(clicontext *cli.Context) error {
	opt, err := buildutil.ParseBuildOpt(clicontext)
	if err != nil {
		return err
	}
	defer func(start time.Time) {
		telemetry.GetReporter().Telemetry(
			"build", telemetry.AddField("duration", time.Since(start).Seconds()))
	}(time.Now())

	logger := logrus.WithField("builder-options", opt)
	logger.Debug("starting build command")

	builder, err := buildutil.GetBuilder(clicontext, opt)
	if err != nil {
		return err
	}
	if err = buildutil.InterpretEnvdDef(builder); err != nil {
		return err
	}
	return buildutil.BuildImage(clicontext, builder)
}
