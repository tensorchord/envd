// Copyright 2022 The envd Authors
// Copyright 2022 The okteto remote Authors
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

// ssh is the CLI running in the container as the sshd.
package main

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	rawssh "golang.org/x/crypto/ssh"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/remote/sshd"
	"github.com/tensorchord/envd/pkg/version"
)

const (
	flagDebug   = "debug"
	flagAuthKey = "authorized-keys"
	flagNoAuth  = "no-auth"
	flagPort    = "port"
	flagShell   = "shell"
	flagHostKey = "hostkey"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Name, version.Package, version.GetVersion().String())
	}

	app := cli.NewApp()
	app.Name = "envd-sshd"
	app.Usage = "ssh server for envd"
	app.Version = version.GetVersion().String()
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  flagDebug,
			Usage: "enable debug output in logs",
		},
		&cli.StringFlag{
			Name:    flagAuthKey,
			Usage:   "path to authorized keys file, defaults to " + config.ContainerAuthorizedKeysPath,
			Value:   config.ContainerAuthorizedKeysPath,
			EnvVars: []string{"ENVD_AUTHORIZED_KEYS_PATH"},
			Aliases: []string{"a"},
		},
		&cli.StringFlag{
			Name:    flagHostKey,
			Usage:   "path to the host key",
			EnvVars: []string{"ENVD_HOST_KEY"},
		},
		&cli.BoolFlag{
			Name:  flagNoAuth,
			Usage: "disable authentication",
			Value: false,
		},
		&cli.IntFlag{
			Name:    flagPort,
			Usage:   "port to listen on",
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:  flagShell,
			Usage: "shell to use",
			Value: "bash",
		},
	}

	// Deal with debug flag.
	var debugEnabled bool

	app.Before = func(context *cli.Context) error {
		debugEnabled = context.Bool(flagDebug)

		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		if debugEnabled {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	app.Action = sshServer
	handleErr(debugEnabled, app.Run(os.Args))
}

func sshServer(c *cli.Context) error {
	err := sshd.GetShell(c.String(flagShell))
	if err != nil {
		logrus.Fatal(err.Error())
	}
	shell := c.String(flagShell)

	port := c.Int(flagPort)
	if port == 0 {
		return errors.New("port must be set")
	} else if port <= 1024 {
		return errors.New("failed to parse port: port is reserved")
	}

	noAuth := c.Bool(flagNoAuth)
	var keys []ssh.PublicKey
	if !noAuth {
		var err error
		path := c.String(flagAuthKey)
		keys, err = sshd.LoadAuthorizedKeys(path)
		if err != nil {
			return errors.Wrapf(err, "failed to load authorized keys at %s", path)
		}

		if keys == nil {
			return errors.Errorf("failed to load authorized keys: no keys found at %s", path)
		}

		logrus.Debugf("loaded %d authorized keys from %s", len(keys), path)
	} else {
		logrus.Warn("no authentication enabled")
	}

	var hostKey ssh.Signer = nil
	if c.String(flagHostKey) != "" {
		// read private key file
		pemBytes, err := os.ReadFile(c.String(flagHostKey))
		if err != nil {
			return errors.Wrapf(
				err, "reading private key %s failed", c.String(flagHostKey))
		}
		if privateKey, err := rawssh.ParsePrivateKey(pemBytes); err != nil {
			return err
		} else {
			logrus.Debugf("load host key from %s", c.String(flagHostKey))
			hostKey = privateKey
		}
	}

	srv := sshd.Server{
		Port:           port,
		Shell:          shell,
		AuthorizedKeys: keys,
		Hostkey:        hostKey,
	}

	logrus.Infof("ssh server %s started in 0.0.0.0:%d", version.GetVersion().String(), srv.Port)
	return srv.ListenAndServe()
}

func handleErr(debug bool, err error) {
	if err == nil {
		return
	}
	if debug {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	os.Exit(1)
}
