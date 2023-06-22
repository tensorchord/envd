// Copyright 2023 The envd Authors
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
	"github.com/tensorchord/envd/pkg/util/buildkitutil"
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
		&cli.BoolFlag{
			Name:  "use-http",
			Usage: "Use HTTP instead of HTTPS for the registry",
			Value: false,
		},
		&cli.StringFlag{
			Name:    "registry",
			Usage:   "Specify private registries and their respective ca/key/cert file paths. Each group is separated by a space. Format: --registry 'registry=my-registry,ca=/etc/config/ca.pem,key=/etc/config/key.pem,cert=/etc/config/cert.pem registry=my-registry2,ca=/etc/config/ca2.pem,key=/etc/config/key2.pem,cert=/etc/config/cert2.pem'",
			Aliases: []string{"r"},
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
		"registry",
		registryCA,
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

func registryCA(clicontext *cli.Context) error {
	// Loop over all registry strings
	if clicontext.IsSet("registry") {
		for _, reg := range strings.Split(clicontext.String("registry"), ";") {
			// Split into parts
			parts := strings.Split(reg, ",")
			if len(parts) != 4 {
				return fmt.Errorf("invalid registry string: %s", reg)
			}

			// Parse registry and ca/key/cert parts
			names := []string{"registry", "ca", "key", "cert"}
			for _, part := range parts {
				kv := strings.SplitN(part, "=", 2)
				if len(kv) != 2 {
					return fmt.Errorf("invalid part: %s", part)
				}

				// Check for valid key
				index := -1
				for i, name := range names {
					if name == kv[0] {
						index = i
						break
					}
				}
				if index == -1 {
					return errors.Newf("parse error: `%s` is not a valid ca/key/cert key or it's duplicated")
				}

				// Read file if not registry part
				if kv[0] != "registry" {
					exist, err := fileutil.FileExists(kv[1])
					if err != nil {
						return errors.Wrap(err, fmt.Sprintf("failed to parse file path %s", part))
					}
					if !exist {
						return fmt.Errorf("file %s doesn't exist", kv[1])
					}

					// Generate file path
					path, err := fileutil.ConfigFile(fmt.Sprintf("registry_%s.pem", kv[0]))
					if err != nil {
						return errors.Wrap(err, "failed to get the envd config file path")
					}

					// Read and write file content
					content, err := os.ReadFile(kv[1])
					if err != nil {
						return errors.Wrap(err, "failed to read the file")
					}
					if err = os.WriteFile(path, content, 0644); err != nil {
						return errors.Wrap(err, "failed to store the CA file")
					}
				}

				// Remove key from names
				names = append(names[:index], names[index+1:]...)
			}

			// Check if any parts were not provided
			if len(names) != 0 {
				return fmt.Errorf("%s parts are not provided", names)
			}
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

		privatePath, err := sshconfig.GetPrivateKey()
		if err != nil {
			return errors.Wrap(err, "Cannot get private key path")
		}

		if keyExists {
			var exists bool
			var newPrivateKeyName string

			for ok := true; ok; ok = exists {
				newPrivateKeyName = filepath.Join(filepath.Dir(privatePath),
					fmt.Sprintf("envd_%s.pk", namesgenerator.GetRandomName(0)))
				exists, err = fileutil.FileExists(newPrivateKeyName)
				if err != nil {
					return err
				}
			}
			logrus.Debugf("New key name: %s", newPrivateKeyName)
			if err := sshconfig.ReplaceKeyManagedByEnvd(
				privatePath, newPrivateKeyName); err != nil {
				return err
			}
		}
		pub, pri := sshKeyPair[0], sshKeyPair[1]
		pubKey, err := os.ReadFile(pub)
		if err != nil {
			return errors.Wrap(err, "Cannot open public key")
		}
		publicPath, err := sshconfig.GetPublicKey()
		if err != nil {
			return errors.Wrap(err, "Cannot get the public key path")
		}
		err = os.WriteFile(publicPath, pubKey, 0644)
		if err != nil {
			return errors.Wrap(err, "Cannot write public key")
		}

		priKey, err := os.ReadFile(pri)
		if err != nil {
			return errors.Wrap(err, "Cannot open private key")
		}

		err = os.WriteFile(privatePath, priKey, 0600)
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
		if err := ac.InsertZSHCompleteEntry(clicontext); err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
		}
	} else if strings.Contains(shell, "bash") {
		logrus.Infof("Install bash autocompletion")
		if err := ac.InsertBashCompleteEntry(clicontext); err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
		}
	} else if strings.Contains(shell, "fish") {
		logrus.Infof("Install fish autocompletion")
		if err := ac.InsertFishCompleteEntry(clicontext); err != nil {
			logrus.Warnf("Warning: %s\n", err.Error())
		}
	} else {
		logrus.Infof("Install bash autocompletion (fallback from \"%s\")", shell)
		if err := ac.InsertBashCompleteEntry(clicontext); err != nil {
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

	config := buildkitutil.BuildkitConfig{
		Mirror:  clicontext.String("dockerhub-mirror"),
		UseHTTP: clicontext.Bool("use-http"),
	}

	// Update the config with the registry info if needed
	if clicontext.IsSet("registry") {
		err = updateConfigWithRegistry(clicontext, &config)
		if err != nil {
			return err
		}
	}

	var bkClient buildkitd.Client
	if c.Builder == types.BuilderTypeMoby {
		bkClient, err = buildkitd.NewMobyClient(clicontext.Context,
			c.Builder, c.BuilderAddress, &config)
		if err != nil {
			return errors.Wrap(err, "failed to create moby buildkit client")
		}
	} else {
		bkClient, err = buildkitd.NewClient(clicontext.Context,
			c.Builder, c.BuilderAddress, &config)
		if err != nil {
			return errors.Wrap(err, "failed to create buildkit client")
		}
	}
	defer bkClient.Close()
	logrus.Infof("The buildkit is running at %s", bkClient.BuildkitdAddr())

	return nil
}

func updateConfigWithRegistry(clicontext *cli.Context, config *buildkitutil.BuildkitConfig) error {
	// Parse registry strings
	dirSet := make(map[string]bool)
	for _, reg := range strings.Split(clicontext.String("registry"), ";") {
		parts := strings.Split(reg, ",")
		if len(parts) != 4 {
			return fmt.Errorf("invalid registry string: %s", reg)
		}

		config.Registry = append(config.Registry, strings.Split(parts[0], "=")[1])
		caPath := strings.Split(parts[1], "=")[1]
		config.Ca = append(config.Ca, caPath)
		config.Key = append(config.Key, strings.Split(parts[2], "=")[1])
		config.Cert = append(config.Cert, strings.Split(parts[3], "=")[1])

		// Add bindings to config
		dir := filepath.Dir(caPath)
		if !dirSet[dir] {
			dirSet[dir] = true
			config.Bindings = append(config.Bindings, fmt.Sprintf("%s:%s", dir, dir))
		}
	}

	return nil
}
