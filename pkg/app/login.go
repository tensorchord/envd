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

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/tensorchord/envd-server/api/types"
	"github.com/tensorchord/envd-server/client"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandLogin = &cli.Command{
	Name:     "login",
	Category: CategoryManagement,
	Aliases:  []string{"i"},
	Usage:    "Login to the envd server.",
	Flags:    []cli.Flag{},
	Action:   login,
}

func login(clicontext *cli.Context) error {
	publicKeyPath, err := fileutil.ConfigFile(config.PublicKeyFile)
	if err != nil {
		return errors.Wrap(err, "failed to get the public key path")
	}
	key, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to read the public key path")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.Wrap(err, "failed to create the envd-server client")
	}
	req := types.AuthRequest{
		PublicKey:     string(key),
		IdentityToken: uuid.NewString(),
	}
	resp, err := cli.Auth(clicontext.Context, req)
	if err != nil {
		return errors.Wrap(err, "failed to get the response from envd-server client")
	}
	fmt.Println(resp.Status)
	return nil
}
