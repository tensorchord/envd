package buildkitd

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"github.com/tensorchord/MIDI/pkg/docker"
	"github.com/tensorchord/MIDI/pkg/flag"
)

type Client interface {
	Bootstrap(ctx context.Context) error
}

type generalClient struct {
	containerName string
	image         string
}

func NewClient() Client {
	return &generalClient{
		containerName: viper.GetString(flag.FlagBuildkitdContainer),
		image:         viper.GetString(flag.FlagBuildkitdImage),
	}
}

func (c generalClient) Bootstrap(ctx context.Context) error {
	_, err := c.MaybeStart(ctx, time.Second*1000)
	if err != nil {
		return err
	}
	return nil
}

// MaybeStart ensures that the buildkitd daemon is started. It returns the URL
// that can be used to connect to it.
func (c generalClient) MaybeStart(
	ctx context.Context, runningTimeout time.Duration) (*string, error) {
	dockerClient, err := docker.NewClient()
	if err != nil {
		return nil, err
	}

	created, err := dockerClient.IsCreated(ctx, c.containerName)
	if err != nil {
		return nil, err
	}

	if !created {
		if _, err := dockerClient.StartBuildkitd(ctx, c.image, c.containerName); err != nil {
			return nil, err
		}
		if err := dockerClient.WaitUntilRunning(ctx, c.containerName, runningTimeout); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
