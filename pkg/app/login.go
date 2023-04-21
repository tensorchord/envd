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
	"os/user"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	servertypes "github.com/tensorchord/envd-server/api/types"
	"github.com/tensorchord/envd-server/client"
	"github.com/tensorchord/envd-server/errdefs"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandLogin = &cli.Command{
	Name:     "login",
	Category: CategoryManagement,
	Hidden:   false,
	Usage:    "Login to the envd server defined in the current context",
	Action:   login,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "username",
			Usage:    "the login name in envd server",
			Aliases:  []string{"u"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "password",
			Usage:    "password",
			Aliases:  []string{"p"},
			Required: false,
		},
	},
}

func login(clicontext *cli.Context) error {
	c, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get current context")
	}
	if c.Runner == types.RunnerTypeDocker {
		logrus.Warn("login is not needed for docker runner, skipping")
		return nil
	}
	hostAddr := c.RunnerAddress

	telemetry.GetReporter().Telemetry("auth", telemetry.AddField("runner", c.Runner))

	publicKeyPath, err := fileutil.ConfigFile(config.PublicKeyFile)
	if err != nil {
		return errors.Wrap(err, "failed to get the public key path")
	}
	key, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to read the public key path")
	}

	stringK := string(key)
	loginName := clicontext.String("username")
	pwd := clicontext.String("password")

	auth := true
	if pwd == "" {
		auth = false
		logrus.Warn("The password is nil, skip the authentication. Please make sure that the server is running in no-auth mode")
		if loginName == "" {
			loginName, err = generateLoginName()
			if err != nil {
				return errors.Wrap(err, "failed to generate the login name")
			}
			logrus.Warnf("The login name is nil, use `%s` as the login name", loginName)
		}
	}
	req := servertypes.AuthNRequest{
		LoginName: loginName,
		Password:  pwd,
	}

	logger := logrus.WithFields(logrus.Fields{
		"login_name":   loginName,
		"auth-enabled": auth,
	})

	opts := []client.Opt{
		client.FromEnv,
	}
	if hostAddr != nil {
		opts = append(opts, client.WithHost(*hostAddr))
	}
	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return errors.Wrap(err, "failed to create the envd-server client")
	}

	var resp servertypes.AuthNResponse

	if auth {
		logger.Debug("login request")
		resp, err = cli.Login(clicontext.Context, req)
		if err != nil {
			return errors.Wrap(err, "failed to get the response from envd-server client")
		}
	} else {
		logger.Debug("register request")
		resp, err = cli.Register(clicontext.Context, req)
		if err != nil {
			if !errdefs.IsConflict(err) {
				return errors.Wrap(err, "failed to get the response from envd-server client")
			}
		}

		logger.Debug("login request after register")
		resp, err = cli.Login(clicontext.Context, req)
		if err != nil {
			return errors.Wrap(err, "failed to get the response from envd-server client")
		}
	}

	// Recreate the cli with the login user.
	opts = []client.Opt{
		client.FromEnv,
		client.WithJWTToken(resp.LoginName, resp.IdentityToken),
	}
	if hostAddr != nil {
		opts = append(opts, client.WithHost(*hostAddr))
	}
	cli, err = client.NewClientWithOpts(opts...)
	if err != nil {
		return errors.Wrap(err, "failed to create the envd-server client")
	}
	keyResp, err := cli.KeyCreate(clicontext.Context, servertypes.KeyCreateRequest{
		Name:      "default",
		PublicKey: stringK,
	})
	if err != nil {
		if !errdefs.IsConflict(err) {
			return errors.Wrap(err, "failed to generate the key")
		}
	}

	logrus.WithField("key", keyResp.Name).Debug("key is added successfully")
	if err := home.GetManager().AuthCreate(types.AuthConfig{
		Name:     resp.LoginName,
		JWTToken: resp.IdentityToken,
	}, true); err != nil {
		return errors.Wrap(err, "failed to create the auth config")
	}
	fmt.Println(resp.Status)
	return nil
}

func generateLoginName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", errors.Wrap(err, "failed to get the hostname")
	}

	username, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "failed to get the user")
	}

	return fmt.Sprintf("%s-%s", username.Username, hostname), nil
}
