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
	"strconv"

	"github.com/cockroachdb/errors"
	rawssh "github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/remote/sshd"
	"github.com/tensorchord/envd/pkg/version"
)

const (
	authorizedKeysPath = "/var/envd/remote/authorized_keys"
	envPort            = "envd_SSH_PORT"
	flagDebug          = "debug"
	flagAuthKey        = "authorized-keys"
	flagNoAuth         = "no-auth"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Name, version.Package, c.App.Version, version.Revision)
	}

	app := cli.NewApp()
	app.Name = "envd-ssh"
	app.Usage = "ssh server for envd"
	app.Version = version.Version
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  flagDebug,
			Usage: "enable debug output in logs",
		},
		&cli.StringFlag{
			Name:    flagAuthKey,
			Usage:   "path to authorized keys file, defaults to " + authorizedKeysPath,
			Value:   authorizedKeysPath,
			Aliases: []string{"a"},
		},
		&cli.BoolFlag{
			Name:  flagNoAuth,
			Usage: "disable authentication",
			Value: false,
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
	shell, err := sshd.GetShell()
	if err != nil {
		logrus.Fatal(err.Error())
	}

	port := 2222
	// TODO(gaocegege): Set it as a flag.
	if p, ok := os.LookupEnv(envPort); ok {
		var err error
		port, err = strconv.Atoi(p)
		if err != nil {
			return errors.Wrap(err, "failed to parse port")
		}

		if port <= 1024 {
			return errors.New("failed to parse port: port is reserved")
		}
	}

	noAuth := c.Bool(flagNoAuth)
	var keys []rawssh.PublicKey
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

	srv := sshd.Server{
		Port:           port,
		Shell:          shell,
		AuthorizedKeys: keys,
	}

	logrus.Infof("ssh server %s started in 0.0.0.0:%d", version.Version, srv.Port)
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
