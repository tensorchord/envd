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
	"io/ioutil"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	ac "github.com/tensorchord/envd/pkg/autocomplete"
	"github.com/tensorchord/envd/pkg/buildkitd"

	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
)

var CommandBootstrap = &cli.Command{
	Name:  "bootstrap",
	Usage: "Bootstraps envd installation including shell autocompletion and buildkit image download",
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
			Usage:   "Manually specify ssh key pair(format: `public_key_path,private_key_path`). Envd will generate it if key doesn't exist at envd home directory.",
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
	} else if len(sshKeyPair) == 0 {

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
			return errors.Errorf("Key already exists at %s and %s. "+
				"Overwriting those keys will break access to existing environments. "+
				"Please backup those keys and manually delete them if you "+
				"want to use your own key for future projects.", sshconfig.GetPublicKey(), sshconfig.GetPrivateKey())
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
		logrus.Debug("bootstrap the buildkitd container")
		bkClient, err := buildkitd.NewClient(clicontext.Context, clicontext.String("dockerhub-mirror"))
		if err != nil {
			return errors.Wrap(err, "failed to create buildkit client")
		}
		defer bkClient.Close()
		logrus.Infof("The buildkit is running at %s", bkClient.BuildkitdAddr())
	}
	return nil
}
