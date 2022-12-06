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
	"github.com/sirupsen/logrus"
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
			return nil, errors.Wrap(err, "failed to get the auth information")
		}

		logger := logrus.WithFields(logrus.Fields{
			"runner":     opt.Context.Runner,
			"login-name": ac.Name,
			"jwt-key":    ac.JWTToken,
		})
		logger.Debug("Creating envd server client")
		// Get the runner host.
		opts := []envdclient.Opt{
			envdclient.WithTLSClientConfigFromEnv(),
			envdclient.WithJWTToken(ac.Name, ac.JWTToken),
		}

		if opt.Context.RunnerAddress != nil {
			opts = append(opts,
				envdclient.WithHost(*opt.Context.RunnerAddress))
		} else {
			opts = append(opts,
				envdclient.WithHostFromEnv())
		}

		cli, err := envdclient.NewClientWithOpts(opts...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create the envd-server client")
		}
		return &envdServerEngine{
			Client:    cli,
			Loginname: ac.Name,
		}, nil
	} else {
		cli, err := client.NewClientWithOpts(
			client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, errors.Wrap(err, "failed to create the docker client")
		}
		return &dockerEngine{
			Client: cli,
		}, nil
	}
}
