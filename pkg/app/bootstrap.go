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
	"encoding/json"
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
			Usage: "Add envd auto-completions",
			Value: true,
		},
		&cli.StringFlag{
			Name:    "dockerhub-mirror",
			Usage:   "DockerHub mirror to use",
			Aliases: []string{"m"},
		},
		&cli.StringFlag{
			Name:    "registry-ca-keypair",
			Usage:   "Specify the ca/key/cert file path for the private registry (format: 'ca=/etc/config/ca.pem,key=/etc/config/key.pem,cert=/etc/config/cert.pem')",
			Aliases: []string{"ca"},
		},
		&cli.StringSliceFlag{
			Name: "ssh-keypair",
			Usage: fmt.Sprintf("Manually specify ssh key pair as `publicKey,privateKey`. envd will generate a keypair at %s and %s if not specified",
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
			Usage:   "Specify the registry to pull the image from",
			Aliases: []string{"r"},
			Value:   "docker.io",
		},
		&cli.StringFlag{
			Name:  "registry-config",
			Usage: "Path to a JSON file containing registry configuration. Cannot be used with 'registry' or 'registry-ca-keypair'",
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
		"registry CA keypair",
		registryCA,
	}, {
		"registry json config",
		registryJSONConfig,
	}, {
		"autocomplete",
		autocomplete,
	}, {
		"buildkit",
		buildkit,
	}, {
		"add pre-defined templates",
		addTemplates,
	}}

	total := len(stages)
	for i, stage := range stages {
		logrus.WithField("cmd", "bootstrap").
			Infof("[%d/%d] Bootstrap %s", i+1, total, stage.Name)
		err := stage.Func(clicontext)
		if err != nil {
			return err
		}
	}

	return nil
}

func registryCA(clicontext *cli.Context) error {
	configFile := clicontext.String("registry-config")
	ca := clicontext.String("registry-ca-keypair")
	registry := clicontext.String("registry")

	if len(ca) == 0 {
		return nil
	}

	// We only need this check in registryCA because it is called before registryJSONConfig
	if len(configFile) > 0 && len(ca) > 0 {
		return errors.New("only one of `registry-config` and `registry-ca-keypair` can be used")
	}

	mirror := clicontext.String("dockerhub-mirror")
	if len(mirror) == 0 {
		return errors.New("`registry-ca-keypair` should be used with `dockerhub-mirror`")
	}

	// Parse ca/key/cert
	kvPairs := strings.Split(ca, ",")
	if len(kvPairs) != 3 {
		return errors.New("`registry-ca-keypair` requires ca/key/cert 3 part separated by ','")
	}
	names := []string{"ca", "cert", "key"}
	for _, pair := range kvPairs {
		kv := strings.SplitN(pair, "=", 2)
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
		exist, err := fileutil.FileExists(kv[1])
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse file path %s", pair))
		}
		if !exist {
			return errors.Newf("file %s doesn't exist", kv[1])
		}

		path, err := fileutil.ConfigFile(fmt.Sprintf("%s_%s.pem", registry, kv[0]))
		if err != nil {
			return errors.Wrap(err, "failed to get the envd config file path")
		}
		content, err := os.ReadFile(kv[1])
		if err != nil {
			return errors.Wrap(err, "failed to read the file")
		}
		if err = os.WriteFile(path, content, 0644); err != nil {
			return errors.Wrap(err, "failed to store the CA file")
		}
		names = append(names[:index], names[index+1:]...)
	}

	if len(names) != 0 {
		return errors.Newf("registry %s are not provided", names)
	}
	return nil
}

func registryJSONConfig(clicontext *cli.Context) error {
	configFile := clicontext.String("registry-config")
	if len(configFile) == 0 {
		return nil
	}

	// Check if file exists
	exist, err := fileutil.FileExists(configFile)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to parse file path %s", configFile))
	}
	if !exist {
		return errors.Newf("file %s doesn't exist", configFile)
	}

	config := buildkitutil.BuildkitConfig{}
	configJson, err := os.ReadFile(configFile)
	if err != nil {
		return errors.Wrap(err, "Failed to read registry config file")
	}
	if err := json.Unmarshal(configJson, &config); err != nil {
		return errors.Wrap(err, "Failed to parse registry config file")
	}

	// Check for required keys in each registry
	for i, registry := range config.Registries {
		if registry.Name == "" {
			return errors.Newf("`name` key is required in the config for registry at index %d", i)
		}

		// Check for optional keys and if they exist, ensure they point to existing files
		optionalKeys := map[string]string{"ca": registry.Ca, "cert": registry.Cert, "key": registry.Key}
		for key, value := range optionalKeys {
			if value != "" {
				exist, err := fileutil.FileExists(value)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("failed to parse file path %s", value))
				}
				if !exist {
					return errors.Newf("file %s doesn't exist", value)
				}

				// Read the file
				content, err := os.ReadFile(value)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("failed to read the %s file for registry %s", key, registry.Name))
				}

				// Write the content to the envd config directory
				envdConfigPath, err := fileutil.ConfigFile(fmt.Sprintf("%s_%s.pem", registry.Name, key))
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("failed to get the envd config file path for %s of registry %s", key, registry.Name))
				}

				if err = os.WriteFile(envdConfigPath, content, 0644); err != nil {
					return errors.Wrap(err, fmt.Sprintf("failed to store the %s file for registry %s", key, registry.Name))
				}
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

			logrus.WithFields(logrus.Fields{
				"cmd":        "bootstrap",
				"stage":      "sshKey",
				"sshKeyPair": sshKeyPair,
			}).Debugf("New key name: %s", newPrivateKeyName)

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
		return errors.Newf("Invalid ssh-keypair flag")
	}
}

func autocomplete(clicontext *cli.Context) error {
	autocomplete := clicontext.Bool("with-autocomplete")
	if !autocomplete {
		return nil
	}

	shell := os.Getenv("SHELL")

	logger := logrus.WithFields(logrus.Fields{
		"cmd":   "bootstrap",
		"stage": "autocomplete",
		"shell": shell,
	})

	if strings.Contains(shell, "zsh") {
		logger.Infof("Install zsh autocompletion")
		if err := ac.InsertZSHCompleteEntry(clicontext); err != nil {
			logger.WithError(err).Warn()
		}
	} else if strings.Contains(shell, "bash") {
		logger.Infof("Install bash autocompletion")
		if err := ac.InsertBashCompleteEntry(clicontext); err != nil {
			logger.WithError(err).Warn()
		}
	} else if strings.Contains(shell, "fish") {
		logger.Infof("Install fish autocompletion")
		if err := ac.InsertFishCompleteEntry(clicontext); err != nil {
			logger.WithError(err).Warn()
		}
	} else {
		logger.Infof(`Install bash autocompletion (fallback from "%s")`, shell)
		if err := ac.InsertBashCompleteEntry(clicontext); err != nil {
			logger.WithError(err).Warn()
		}
	}

	logger.Info(`You may have to restart your shell for autocomplete to get initialized (e.g. run "exec $SHELL")\n`)
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

	logger := logrus.WithFields(logrus.Fields{
		"cmd":   "bootstrap",
		"stage": "buildkit",
	})

	logger.Debug("bootstrap the buildkitd container")
	// Populate the BuildkitConfig struct
	config := buildkitutil.BuildkitConfig{}

	configFile := clicontext.String("registry-config")
	if len(configFile) != 0 {
		configJson, err := os.ReadFile(configFile)
		if err != nil {
			return errors.Wrap(err, "Failed to read registry config file")
		}
		if err := json.Unmarshal(configJson, &config); err != nil {
			return errors.Wrap(err, "Failed to parse registry config file")
		}
	} else if len(clicontext.String("registry-ca-keypair")) != 0 {
		// The values of Ca, Cert, and Key don't actually matter since we already copied their contents to the envd config directory and mounted to `/etc/registry`.
		// So instead of parsing registry-ca-keypair again, we'll just put the default value.
		// This is to ensure that buildkitConfigTemplate parses properly.
		config.Registries = append(config.Registries, buildkitutil.Registry{
			Name:    clicontext.String("registry"),
			Ca:      "/etc/registry",
			Cert:    "/etc/registry",
			Key:     "/etc/registry",
			UseHTTP: clicontext.Bool("use-http"),
			Mirror:  clicontext.String("dockerhub-mirror"),
		})
	} else if len(clicontext.String("dockerhub-mirror")) != 0 {
		config.Registries = append(config.Registries, buildkitutil.Registry{
			Name:    clicontext.String("registry"),
			UseHTTP: clicontext.Bool("use-http"),
			Mirror:  clicontext.String("dockerhub-mirror"),
		})
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
	logger.Infof("The buildkit is running at %s", bkClient.BuildkitdAddr())

	return nil
}
