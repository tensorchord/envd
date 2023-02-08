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
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	ac "github.com/tensorchord/envd/pkg/autocomplete"
	"github.com/tensorchord/envd/pkg/buildkitd"
	"github.com/tensorchord/envd/pkg/home"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandBootstrap = &cli.Command{
	Name:     "bootstrap",
	Category: CategorySettings,
	Usage:    "Bootstrap the envd installation",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "buildkit",
			Usage:   "Download the image and bootstrap buildkit",
			Aliases: []string{"b"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:  "with-autocomplete",
			Usage: "Add envd autocompletions",
			Value: true,
		},
		&cli.StringFlag{
			Name:    "dockerhub-mirror",
			Usage:   "Dockerhub mirror to use",
			Aliases: []string{"m"},
		},
		&cli.StringSliceFlag{
			Name: "ssh-keypair",
			Usage: fmt.Sprintf("Manually specify ssh key pair as `publicKey,privateKey`. Envd will generate a keypair at %s and %s if not specified",
				sshconfig.GetPublicKeyOrPanic(), sshconfig.GetPrivateKeyOrPanic()),
			Aliases: []string{"k"},
		},
	},

	Action: bootstrap,
}

func bootstrap(clicontext *cli.Context) error {
	stages := []struct {
		Name string
		Func func(*cli.Context) error
	}{{
		"SSH Key",
		sshKey,
	}, {
		"autocomplete",
		autocomplete,
	}, {
		"buildkit",
		buildkit,
	}}

	total := len(stages)
	for i, stage := range stages {
		logrus.Infof("[%d/%d] Bootstrap %s", i+1, total, stage.Name)
		err := stage.Func(clicontext)
		if err != nil {
			return err
		}
	}

	return nil
}

func sshKey(clicontext *cli.Context) error {
	sshKeyPair := clicontext.StringSlice("ssh-keypair")

	switch len(sshKeyPair) {
	case 0:
		// If not specified, generate only if key doesn't exist
		keyExists, err := sshconfig.DefaultKeyExists()
		if err != nil {
			return errors.Wrap(err, "Cannot get default key status")
		}
		if !keyExists {
			// Generate SSH keys only if key doesn't exist
			if err := sshconfig.GenerateKeys(); err != nil {
				return errors.Wrap(err, "failed to generate ssh key")
			}
		}
		return nil
	case 2:
		// Specified new pairs. Only change to specified pairs when default key doesn't exist
		// Raise error if already exists
		keyExists, err := sshconfig.DefaultKeyExists()
		if err != nil {
			return errors.Wrap(err, "Cannot get default key status")
		}

		path, err := sshconfig.GetPrivateKey()
		if err != nil {
			return errors.Wrap(err, "Cannot get private key path")
		}

		if keyExists {
			var exists bool
			var newPrivateKeyName string

			for ok := true; ok; ok = exists {
				newPrivateKeyName = filepath.Join(filepath.Dir(path),
					fmt.Sprintf("envd_%s.pk", namesgenerator.GetRandomName(0)))
				exists, err = fileutil.FileExists(newPrivateKeyName)
				if err != nil {
					return err
				}
			}
			logrus.Debugf("New key name: %s", newPrivateKeyName)
			if err := sshconfig.ReplaceKeyManagedByEnvd(
				path, newPrivateKeyName); err != nil {
				return err
			}
		}
		pub, pri := sshKeyPair[0], sshKeyPair[1]
		pubKey, err := os.ReadFile(pub)
		if err != nil {
			return errors.Wrap(err, "Cannot open public key")
		}
		err = os.WriteFile(path, pubKey, 0644)
		if err != nil {
			return errors.Wrap(err, "Cannot write public key")
		}

		priKey, err := os.ReadFile(pri)
		if err != nil {
			return errors.Wrap(err, "Cannot open private key")
		}

		err = os.WriteFile(path, priKey, 0600)
		if err != nil {
			return errors.Wrap(err, "Cannot write private key")
		}
		return nil

	default:
		return errors.Errorf("Invalid ssh-keypair flag")
	}
}

func autocomplete(clicontext *cli.Context) error {
	autocomplete := clicontext.Bool("with-autocomplete")
	if !autocomplete {
		return nil
	}

	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		logrus.Infof("Install zsh autocompletion")
		if err := ac.InsertZSHCompleteEntry(); err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
		}
	} else if strings.Contains(shell, "bash") {
		logrus.Infof("Install bash autocompletion")
		if err := ac.InsertBashCompleteEntry(); err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
		}
	} else {
		logrus.Infof("Install bash autocompletion (fallback from \"%s\")", shell)
		if err := ac.InsertBashCompleteEntry(); err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
		}
	}

	logrus.Info("You may have to restart your shell for autocomplete to get initialized (e.g. run \"exec $SHELL\")\n")
	return nil
}

func buildkit(clicontext *cli.Context) error {
	buildkit := clicontext.Bool("buildkit")
	if !buildkit {
		return nil
	}

	c, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}

	logrus.Debug("bootstrap the buildkitd container")
	var bkClient buildkitd.Client
	if c.Builder == types.BuilderTypeMoby {
		bkClient, err = buildkitd.NewMobyClient(clicontext.Context,
			c.Builder, c.BuilderAddress, clicontext.String("dockerhub-mirror"))
		if err != nil {
			return errors.Wrap(err, "failed to create moby buildkit client")
		}
	} else {
		bkClient, err = buildkitd.NewClient(clicontext.Context,
			c.Builder, c.BuilderAddress, clicontext.String("dockerhub-mirror"))
		if err != nil {
			return errors.Wrap(err, "failed to create buildkit client")
		}
	}
	defer bkClient.Close()
	logrus.Infof("The buildkit is running at %s", bkClient.BuildkitdAddr())

	return nil
}
