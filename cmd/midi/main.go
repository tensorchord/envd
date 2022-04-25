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
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/MIDI/pkg/flag"
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
			Name:  flag.FlagCacheDir,
			Usage: "cache directory",
			Value: "~/.midi/cache",
		},
		&cli.PathFlag{
			Name:  flag.FlagConfig,
			Usage: "path to config file",
			Value: "~/.midi/config.MIDI",
		},
	}

	app.Commands = []*cli.Command{
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

		// Setup the cache directory.
		cacheDir := context.Path(flag.FlagCacheDir)
		if strings.HasPrefix(cacheDir, "~/") {
			usr, _ := user.Current()
			dir := usr.HomeDir
			cacheDir = filepath.Join(dir, cacheDir[2:])
		}
		viper.Set(flag.FlagCacheDir, cacheDir)
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return err
		}

		// Get the config file.
		configFile := context.Path(flag.FlagConfig)
		if strings.HasPrefix(configFile, "~/") {
			usr, _ := user.Current()
			dir := usr.HomeDir
			configFile = filepath.Join(dir, configFile[2:])
		}
		viper.Set(flag.FlagConfig, configFile)
		if _, err := os.Stat(configFile); err != nil {
			if os.IsNotExist(err) {
				if _, err := os.Create(configFile); err != nil {
					return err
				}
			}
		} else {
			return err
		}

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
