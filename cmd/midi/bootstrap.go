// Copyright 2022 The MIDI Authors
// Copyright 2022 The midi Authors
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
	"runtime"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/containerd/containerd/log"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/MIDI/pkg/buildkitd"
	"github.com/tensorchord/MIDI/pkg/util/fileutil"
)

var CommandBootstrap = &cli.Command{
	Name:  "bootstrap",
	Usage: "Bootstraps midi installation including shell autocompletion and buildkit image download",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "buildkit",
			Usage:   "Download the image and bootstrap buildkit",
			Aliases: []string{"b"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:  "with-autocomplete",
			Usage: "Add midi autocompletions",
			Value: true,
		},
	},

	Action: bootstrap,
}

func bootstrap(clicontext *cli.Context) error {
	autocomplete := clicontext.Bool("with-autocomplete")
	if autocomplete {
		// Because this requires sudo, it should warn and not fail the rest of it.
		err := insertBashCompleteEntry()
		if err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
			err = nil
		}
		err = insertZSHCompleteEntry()
		if err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
			err = nil
		}

		logrus.Info("You may have to restart your shell for autocomplete to get initialized (e.g. run \"exec $SHELL\")\n")
	}

	buildkit := clicontext.Bool("buildkit")

	if buildkit {
		logrus.Debug("bootstrap the buildkitd container")
		bkClient := buildkitd.NewClient()
		defer bkClient.Close()
		addr, err := bkClient.Bootstrap(clicontext.Context)
		if err != nil {
			return errors.Wrap(err, "failed to bootstrap buildkit")
		}
		logrus.Infof("The buildkit is running at %s", addr)
	}
	return nil
}

// If debugging this, it might be required to run `rm ~/.zcompdump*` to remove the cache
func insertZSHCompleteEntry() error {
	// should be the same on linux and macOS
	path := "/usr/local/share/zsh/site-functions/_midi"
	dirPath := filepath.Dir(path)

	dirPathExists, err := fileutil.DirExists(dirPath)
	if err != nil {
		return errors.Wrapf(err, "failed to check if %s exists", dirPath)
	}
	if !dirPathExists {
		log.L.Warn("Warning: unable to enable zsh-completion: %s does not exist\n", dirPath)
		return nil // zsh-completion isn't available, silently fail.
	}

	pathExists, err := fileutil.FileExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed to check if %s exists", path)
	}
	if pathExists {
		return nil // file already exists, don't update it.
	}

	// create the completion file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	compEntry, err := zshCompleteEntry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable zsh-completion: %s\n", err)
		return nil // zsh-completion isn't available, silently fail.
	}

	_, err = f.Write([]byte(compEntry))
	if err != nil {
		return errors.Wrapf(err, "failed writing to %s", path)
	}

	return deleteZcompdump()
}

func zshCompleteEntry() (string, error) {
	template := `#compdef _midi midi

function _midi {
    autoload -Uz bashcompinit
    bashcompinit
    complete -o nospace -C '__midi__' midi
}
`
	return renderEntryTemplate(template)
}

func insertBashCompleteEntry() error {
	var path string
	if runtime.GOOS == "darwin" {
		path = "/usr/local/etc/bash_completion.d/midi"
	} else {
		path = "/usr/share/bash-completion/completions/midi"
	}
	dirPath := filepath.Dir(path)

	dirPathExists, err := fileutil.DirExists(dirPath)
	if err != nil {
		return errors.Wrapf(err, "failed checking if %s exists", dirPath)
	}
	if !dirPathExists {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable bash-completion: %s does not exist\n", dirPath)
		return nil // bash-completion isn't available, silently fail.
	}

	pathExists, err := fileutil.FileExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed checking if %s exists", path)
	}
	if pathExists {
		return nil // file already exists, don't update it.
	}

	// create the completion file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bashEntry, err := bashCompleteEntry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable bash-completion: %s\n", err)
		return nil // bash-completion isn't available, silently fail.
	}

	_, err = f.Write([]byte(bashEntry))
	if err != nil {
		return errors.Wrapf(err, "failed writing to %s", path)
	}
	return nil
}

func bashCompleteEntry() (string, error) {
	template := "complete -o nospace -C '__midi__' midi\n"
	return renderEntryTemplate(template)
}

func renderEntryTemplate(template string) (string, error) {
	midiPath, err := os.Executable()
	if err != nil {
		return "", errors.Wrapf(err, "failed to determine midi path: %s", err)
	}
	return strings.ReplaceAll(template, "__midi__", midiPath), nil
}

func deleteZcompdump() error {
	var homeDir string
	sudoUser, found := os.LookupEnv("SUDO_USER")
	if !found {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return errors.Wrapf(err, "failed to lookup current user home dir")
		}
	} else {
		currentUser, err := user.Lookup(sudoUser)
		if err != nil {
			return errors.Wrapf(err, "failed to lookup user %s", sudoUser)
		}
		homeDir = currentUser.HomeDir
	}
	files, err := os.ReadDir(homeDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read dir %s", homeDir)
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".zcompdump") {
			path := filepath.Join(homeDir, f.Name())
			err := os.Remove(path)
			if err != nil {
				return errors.Wrapf(err, "failed to remove %s", path)
			}
		}
	}
	return nil
}
