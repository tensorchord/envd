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
	"fmt"
	"os"

	_ "github.com/moby/buildkit/client/connhelper/dockercontainer"
	_ "github.com/moby/buildkit/client/connhelper/kubepod"
	_ "github.com/moby/buildkit/client/connhelper/podmancontainer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/MIDI/pkg/flag"
	"github.com/tensorchord/MIDI/pkg/home"
	"github.com/tensorchord/MIDI/pkg/version"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Name, version.Package, c.App.Version, version.Revision)
	}

	// TODO(gaocegege): Enclose the app, maybe create the struct MIDIApp.
	app := cli.NewApp()
	app.Name = "midi"
	app.Usage = "Build tools for data scientists"
	app.Version = version.Version
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		&cli.PathFlag{
			Name:  flag.FlagConfig,
			Usage: "path to config file",
			Value: "~/.midi/config.MIDI",
		},
		&cli.PathFlag{
			Name:  flag.FlagHomeDir,
			Usage: "path to midi home",
			Value: "~/.midi",
		},
		&cli.StringFlag{
			Name:  flag.FlagBuildkitdImage,
			Usage: "docker image to use for buildkitd",
			Value: "docker.io/moby/buildkit:v0.10.1",
		},
		&cli.StringFlag{
			Name:  flag.FlagBuildkitdContainer,
			Usage: "buildkitd container to use for buildkitd",
			Value: "midi_buildkitd",
		},
	}

	app.Commands = []*cli.Command{
		CommandBootstrap,
		CommandBuild,
		CommandUp,
	}

	// Deal with debug flag.
	var debugEnabled bool

	app.Before = func(context *cli.Context) error {
		debugEnabled = context.Bool("debug")

		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		if debugEnabled {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// Get the config file.
		configFile := context.Path(flag.FlagConfig)

		homeDir := context.Path(flag.FlagHomeDir)

		if err := home.Intialize(homeDir, configFile); err != nil {
			return errors.Wrap(err, "failed to initialize home manager")
		}

		// TODO(gaocegege): Add a config struct to keep them.
		viper.Set(flag.FlagBuildkitdContainer, context.String(flag.FlagBuildkitdContainer))
		viper.Set(flag.FlagBuildkitdImage, context.String(flag.FlagBuildkitdImage))

		return nil
	}
	handleErr(debugEnabled, app.Run(os.Args))
}

func handleErr(debug bool, err error) {
	if err == nil {
		return
	}
	// TODO(gaocegege): Print the error with starlark stacks.
	if debug {
		// TODO(gaocegege): Add debug info.
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	os.Exit(1)
}
