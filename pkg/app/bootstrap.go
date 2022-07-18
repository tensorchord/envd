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
	"io/ioutil"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	ac "github.com/tensorchord/envd/pkg/autocomplete"
	"github.com/tensorchord/envd/pkg/buildkitd"
	"github.com/tensorchord/envd/pkg/home"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandBootstrap = &cli.Command{
	Name:  "bootstrap",
	Usage: "Bootstrap the envd installation including shell autocompletion and buildkit image download",
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
			Name:    "ssh-keypair",
			Usage:   fmt.Sprintf("Manually specify ssh key pair as `publicKey,privateKey`. Envd will generate a keypair at %s and %s if not specified", sshconfig.GetPublicKey(), sshconfig.GetPrivateKey()),
			Aliases: []string{"k"},
		},
	},

	Action: bootstrap,
}

func bootstrap(clicontext *cli.Context) error {
	// If not specified, generate only if key doesn't exist
	sshKeyPair := clicontext.StringSlice("ssh-keypair")
	if len(sshKeyPair) != 0 && len(sshKeyPair) != 2 {
		return errors.Errorf("Invliad ssh-keypair flag")
	}
	if len(sshKeyPair) == 0 {
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
	} else if len(sshKeyPair) == 2 {
		// Specified new pairs. Only change to specified pairs when default key doesn't exist
		// Raise error if already exists
		keyExists, err := sshconfig.DefaultKeyExists()
		if err != nil {
			return errors.Wrap(err, "Cannot get default key status")
		}
		if keyExists {
			var exists bool
			var newPrivateKeyName string
			for ok := true; ok; ok = exists {
				newPrivateKeyName = filepath.Join(filepath.Dir(sshconfig.GetPrivateKey()), fmt.Sprintf("envd_%s.pk", namesgenerator.GetRandomName(0)))
				exists, err = fileutil.FileExists(newPrivateKeyName)
				if err != nil {
					return err
				}
			}
			logrus.Debugf("New key name: %s", newPrivateKeyName)
			err := sshconfig.ReplaceKeyManagedByEnvd(sshconfig.GetPrivateKey(), newPrivateKeyName)
			if err != nil {
				return err
			}
		}
		pub, pri := sshKeyPair[0], sshKeyPair[1]
		pubKey, err := ioutil.ReadFile(pub)
		if err != nil {
			return errors.Wrap(err, "Cannot open public key")
		}
		err = ioutil.WriteFile(sshconfig.GetPublicKey(), pubKey, 0644)
		if err != nil {
			return errors.Wrap(err, "Cannot write public key")
		}

		priKey, err := ioutil.ReadFile(pri)
		if err != nil {
			return errors.Wrap(err, "Cannot open private key")
		}
		err = ioutil.WriteFile(sshconfig.GetPrivateKey(), priKey, 0600)

		if err != nil {
			return errors.Wrap(err, "Cannot write private key")
		}

	}

	autocomplete := clicontext.Bool("with-autocomplete")
	if autocomplete {
		// Because this requires sudo, it should warn and not fail the rest of it.
		err := ac.InsertBashCompleteEntry()
		if err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
			err = nil
		}
		err = ac.InsertZSHCompleteEntry()
		if err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
			err = nil
		}

		logrus.Info("You may have to restart your shell for autocomplete to get initialized (e.g. run \"exec $SHELL\")\n")
	}

	buildkit := clicontext.Bool("buildkit")

	if buildkit {
		currentDriver, currentSocket, err := home.GetManager().ContextGetCurrent()
		if err != nil {
			return errors.Wrap(err, "failed to get the current context")
		}
		logrus.Debug("bootstrap the buildkitd container")
		bkClient, err := buildkitd.NewClient(clicontext.Context,
			currentDriver, currentSocket, clicontext.String("dockerhub-mirror"))
		if err != nil {
			return errors.Wrap(err, "failed to create buildkit client")
		}
		defer bkClient.Close()
		logrus.Infof("The buildkit is running at %s", bkClient.BuildkitdAddr())
	}
	return nil
}
