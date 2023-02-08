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
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer"
	_ "github.com/moby/buildkit/client/connhelper/kubepod"
	_ "github.com/moby/buildkit/client/connhelper/nerdctlcontainer"
	_ "github.com/moby/buildkit/client/connhelper/podmancontainer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/version"
)

type EnvdApp struct {
	cli.App
}

func New() EnvdApp {
	internalApp := cli.NewApp()
	internalApp.EnableBashCompletion = true
	internalApp.Name = "envd"
	internalApp.Usage = "Development environment for data science and AI/ML teams"
	internalApp.HideHelpCommand = true
	internalApp.HideVersion = true
	internalApp.Version = version.GetVersion().String()
	internalApp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		&cli.BoolFlag{
			Name:    flag.FlagAnalytics,
			Usage:   "enable analytics",
			Value:   true,
			EnvVars: []string{"ENVD_ANALYTICS"},
		},
		&cli.StringFlag{
			Name:    "analytics-token",
			Value:   "Kdj4DfmpMoFh9qia2gz9wJq7TdrfVsP7",
			EnvVars: []string{"ENVD_ANALYTICS_TOKEN"},
			Hidden:  true,
		},
		&cli.StringFlag{
			Name:   flag.FlagBuildkitdImage,
			Usage:  "docker image to use for buildkitd",
			Value:  "docker.io/moby/buildkit:v0.10.3",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   flag.FlagDockerOrganization,
			Usage:  "docker organization to use",
			Value:  "tensorchord",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:    flag.FlagBuildOwner,
			Usage:   "owner (uid:gid) in the built image",
			Value:   "",
			EnvVars: []string{"ENVD_BUILD_OWNER"},
			Hidden:  true,
		},
	}

	internalApp.Commands = []*cli.Command{
		CommandBootstrap,
		CommandCreate,
		CommandCompletion,
		CommandContext,
		CommandBuild,
		CommandDestroy,
		CommandEnvironment,
		CommandImage,
		CommandInit,
		CommandLogin,
		CommandPause,
		CommandPrune,
		CommandRun,
		CommandResume,
		CommandUp,
		CommandDebug,
		CommandVersion,
		CommandTop,
	}

	internalApp.CustomAppHelpTemplate = ` envd - Development environment for data science and AI/ML teams

 Usage:
    envd up --path <path>

    envd build --path <path> --output type=image,name=<image>,push=true{{if .VisibleCommands}}

 Build and launch envd environments. Get more information at: https://envd.tensorchord.ai/.
 To get started with using envd, check out the getting started guide: https://envd.tensorchord.ai/guide/getting-started.html.
 {{range .VisibleCategories}}{{if .Name}}
 {{.Name}}:{{range .VisibleCommands}}
	{{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
 {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{end}}{{if .VisibleFlagCategories}}

 Global Options:{{range .VisibleFlagCategories}}
	{{if .Name}}{{.Name}}
	{{end}}{{range .Flags}}{{.}}
	{{end}}{{end}}{{else}}{{if .VisibleFlags}}

 Global Options:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{wrap $option.String 6}}{{end}}{{end}}{{end}}
`

	// Deal with debug flag.
	var debugEnabled bool

	internalApp.Before = func(context *cli.Context) error {
		debugEnabled = context.Bool("debug")

		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		if debugEnabled {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if err := home.Initialize(); err != nil {
			return errors.Wrap(err, "failed to initialize home manager")
		}

		analytics := context.Bool(flag.FlagAnalytics)
		if err := telemetry.Initialize(analytics,
			context.String("analytics-token")); err != nil {
			return errors.Wrap(err, "failed to create telemetry client")
		}

		// TODO(gaocegege): Add a config struct to keep them.
		viper.Set(flag.FlagBuildkitdImage, context.String(flag.FlagBuildkitdImage))
		viper.Set(flag.FlagBuildOwner, context.String(flag.FlagBuildOwner))
		viper.Set(flag.FlagDebug, debugEnabled)
		viper.Set(flag.FlagAnalytics, analytics)
		viper.Set(flag.FlagDockerOrganization,
			context.String(flag.FlagDockerOrganization))
		return nil
	}

	return EnvdApp{
		App: *internalApp,
	}
}
