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

package envd

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/client"
	envdclient "github.com/tensorchord/envd-server/client"

	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

type Options struct {
	Context *types.Context
}

func New(ctx context.Context, opt Options) (Engine, error) {
	if opt.Context == nil {
		return nil, errors.New("failed to get the context")
	}
	if opt.Context.Runner == types.RunnerTypeEnvdServer {
		ac, err := home.GetManager().AuthGetCurrent()
		if err != nil {
			return nil, err
		}

		cli, err := envdclient.NewClientWithOpts(envdclient.FromEnv)
		if err != nil {
			return nil, err
		}
		return &envdServerEngine{
			Client:        cli,
			IdentityToken: ac.IdentityToken,
		}, nil
	} else {
		cli, err := client.NewClientWithOpts(
			client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}
		return &dockerEngine{
			Client: cli,
		}, nil
	}
}
