package app

import (
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/version"
)

type EnvdApp struct {
	cli.App
	Debug bool
}

func New() EnvdApp {
	internalApp := cli.NewApp()
	internalApp.EnableBashCompletion = true
	internalApp.Name = "envd"
	internalApp.Usage = "Build tools for data scientists"
	internalApp.Version = version.GetVersion().String()
	internalApp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		&cli.StringFlag{
			Name:  flag.FlagBuildkitdImage,
			Usage: "docker image to use for buildkitd",
			Value: "docker.io/moby/buildkit:v0.10.3",
		},
		&cli.StringFlag{
			Name:  flag.FlagBuildkitdContainer,
			Usage: "buildkitd container to use for buildkitd",
			Value: "envd_buildkitd",
		},
	}

	internalApp.Commands = []*cli.Command{
		CommandBootstrap,
		CommandBuild,
		CommandDestroy,
		CommandGet,
		CommandPause,
		CommandResume,
		CommandUp,
		CommandVersion,
	}

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

		// TODO(gaocegege): Add a config struct to keep them.
		viper.Set(flag.FlagBuildkitdContainer, context.String(flag.FlagBuildkitdContainer))
		viper.Set(flag.FlagBuildkitdImage, context.String(flag.FlagBuildkitdImage))
		return nil
	}

	return EnvdApp{
		App:   *internalApp,
		Debug: debugEnabled,
	}
}
