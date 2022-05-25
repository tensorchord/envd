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
