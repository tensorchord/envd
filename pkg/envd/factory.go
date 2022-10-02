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

	"github.com/docker/docker/client"
	envdclient "github.com/tensorchord/envd-server/client"
)

func New(ctx context.Context, backend string) (Engine, error) {
	if backend == "docker" {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}
		return &dockerEngine{
			Client: cli,
		}, nil
	} else {
		cli, err := envdclient.NewClientWithOpts(envdclient.FromEnv)
		if err != nil {
			return nil, err
		}
		return &envdServerEngine{
			Client: cli,
		}, nil
	}
}
