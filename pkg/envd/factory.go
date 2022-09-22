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
