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
	servertypes "github.com/tensorchord/envd-server/api/types"
	"github.com/tensorchord/envd-server/client"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/home"
)

var CommandK8s = &cli.Command{
	Name:     "k8s",
	Category: CategoryBasic,
	Hidden:   true,
	Usage:    "TestK8s",
	Action:   k8s,
}

func k8s(clicontext *cli.Context) error {
	ac, err := home.GetManager().AuthGetCurrent()
	if err != nil {
		return err
	}
	it := ac.IdentityToken
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	req := servertypes.EnvironmentCreateRequest{
		IdentityToken: it,
		Image:         "gaocegege/test-envd",
	}
	_, err = c.EnvironmentCreate(clicontext.Context, req)
	if err != nil {
		return err
	}

	return nil
}
