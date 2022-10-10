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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	servertypes "github.com/tensorchord/envd-server/api/types"
	"github.com/tensorchord/envd-server/client"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandLogin = &cli.Command{
	Name:     "login",
	Category: CategoryManagement,
	Hidden:   true,
	Usage:    "Login to the envd server.",
	Action:   login,
}

func login(clicontext *cli.Context) error {
	publicKeyPath, err := fileutil.ConfigFile(config.PublicKeyFile)
	if err != nil {
		return errors.Wrap(err, "failed to get the public key path")
	}
	key, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to read the public key path")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.Wrap(err, "failed to create the envd-server client")
	}
	stringK := string(key)
	it := GetMD5Hash(stringK)
	req := servertypes.AuthRequest{
		PublicKey:     stringK,
		IdentityToken: it,
	}
	resp, err := cli.Auth(clicontext.Context, req)
	if err != nil {
		return errors.Wrap(err, "failed to get the response from envd-server client")
	}
	if err := home.GetManager().AuthCreate(types.AuthConfig{
		Name:          resp.IdentityToken,
		IdentityToken: resp.IdentityToken,
	}, true); err != nil {
		return errors.Wrap(err, "failed to create the auth config")
	}
	fmt.Println(resp.Status)
	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
