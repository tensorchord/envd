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
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/types"
)

type Engine interface {
	List(ctx context.Context) ([]types.EnvdEnvironment, error)
}

type generalEngine struct {
	dockerCli docker.Client
}

func New(ctx context.Context) (Engine, error) {
	dc, err := docker.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create docker client")
	}
	return &generalEngine{
		dockerCli: dc,
	}, nil
}

func (e generalEngine) List(ctx context.Context) ([]types.EnvdEnvironment, error) {
	ctrs, err := e.dockerCli.List(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list containers")
	}

	envs := make([]types.EnvdEnvironment, 0)
	for _, ctr := range ctrs {
		envs = append(envs, types.FromContainer(ctr))
	}
	return envs, nil
}
