package envd

import (
	"context"

	"github.com/docker/docker/client"
)

func New(ctx context.Context) (Engine, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &dockerEngine{
		Client: cli,
	}, nil
}
