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
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/types"
)

type Engine interface {
	ListImage(ctx context.Context) ([]types.EnvdImage, error)
	ListImageDependency(ctx context.Context, image string) (*types.Dependency, error)

	PauseEnvironment(ctx context.Context, env string) (string, error)
	ResumeEnvironment(ctx context.Context, env string) (string, error)
	ListEnvironment(ctx context.Context) ([]types.EnvdEnvironment, error)
	ListEnvDependency(ctx context.Context, env string) (*types.Dependency, error)
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

func (e generalEngine) ListImage(ctx context.Context) ([]types.EnvdImage, error) {
	imgs, err := e.dockerCli.ListImage(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list images")
	}
	envdImgs := make([]types.EnvdImage, 0)
	for _, img := range imgs {
		envdImg, err := types.NewImage(img)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create envd image from the docker image")
		}
		envdImgs = append(envdImgs, *envdImg)
	}
	return envdImgs, nil
}

func (e generalEngine) ListEnvironment(
	ctx context.Context) ([]types.EnvdEnvironment, error) {
	ctrs, err := e.dockerCli.ListContainer(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list containers")
	}

	envs := make([]types.EnvdEnvironment, 0)
	for _, ctr := range ctrs {
		env, err := types.NewEnvironment(ctr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create env from the container")
		}
		envs = append(envs, *env)
	}
	return envs, nil
}

func (e generalEngine) PauseEnvironment(ctx context.Context, env string) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"env": env,
	})
	logger.Debug("pausing environment")
	name, err := e.dockerCli.PauseContainer(ctx, env)
	if err != nil {
		return "", errors.Wrap(err, "failed to pause the environment")
	}
	return name, nil
}

func (e generalEngine) ResumeEnvironment(ctx context.Context, env string) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"env": env,
	})
	logger.Debug("resuming environment")
	name, err := e.dockerCli.ResumeContainer(ctx, env)
	if err != nil {
		return "", errors.Wrap(err, "failed to resume the environment")
	}
	return name, nil
}

// ListEnvDependency gets the dependencies of the given environment.
func (e generalEngine) ListImageDependency(
	ctx context.Context, image string) (*types.Dependency, error) {
	logger := logrus.WithFields(logrus.Fields{
		"image": image,
	})
	logger.Debug("getting dependencies")
	img, err := e.dockerCli.GetImage(ctx, image)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get container")
	}
	dep, err := types.NewDependencyFromImage(img)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dependency from image")
	}
	return dep, nil
}

// ListEnvDependency gets the dependencies of the given environment.
func (e generalEngine) ListEnvDependency(
	ctx context.Context, env string) (*types.Dependency, error) {
	logger := logrus.WithFields(logrus.Fields{
		"env": env,
	})
	logger.Debug("getting dependencies")
	ctr, err := e.dockerCli.GetContainer(ctx, env)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get container")
	}
	dep, err := types.NewDependencyFromContainerJSON(ctr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dependency from the container")
	}
	return dep, nil
}
