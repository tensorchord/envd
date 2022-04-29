package buildkitd

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/tensorchord/MIDI/pkg/docker"
	"github.com/tensorchord/MIDI/pkg/flag"
)

var (
	interval = time.Second * 1
)

// Client is a client for the buildkitd daemon.
// It's up to the caller to close the client.
type Client interface {
	BuildkitdAddr() string
	Bootstrap(ctx context.Context) (string, error)
	// Solve calls Solve on the controller.
	Solve(ctx context.Context, def *llb.Definition, opt client.SolveOpt, statusChan chan *client.SolveStatus) (*client.SolveResponse, error)
	Close() error
}

type generalClient struct {
	containerName string
	image         string

	*client.Client
	logger *logrus.Entry
}

func NewClient(ctx context.Context) (Client, error) {
	c := &generalClient{
		containerName: viper.GetString(flag.FlagBuildkitdContainer),
		image:         viper.GetString(flag.FlagBuildkitdImage),
	}
	c.logger = logrus.WithFields(logrus.Fields{
		"container": c.containerName,
		"image":     c.image,
	})

	cli, err := client.New(ctx, c.BuildkitdAddr(), client.WithFailFast())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the client")
	}
	c.Client = cli

	if _, err := c.Bootstrap(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to bootstrap the buildkitd")
	}
	return c, nil
}

func (c *generalClient) Bootstrap(ctx context.Context) (string, error) {
	address, err := c.maybeStart(ctx, time.Second*100, time.Second*100)
	if err != nil {
		return "", err
	}
	return address, nil
}

func (c generalClient) Close() error {
	return c.Client.Close()
}

// maybeStart ensures that the buildkitd daemon is started. It returns the URL
// that can be used to connect to it.
func (c *generalClient) maybeStart(ctx context.Context,
	runningTimeout, connectingTimeout time.Duration) (string, error) {
	dockerClient, err := docker.NewClient()
	if err != nil {
		return "", err
	}

	created, err := dockerClient.IsCreated(ctx, c.containerName)
	if err != nil {
		return "", err
	}

	if !created {
		c.logger.Debug("container not created, creating...")
		if _, err := dockerClient.StartBuildkitd(ctx, c.image, c.containerName); err != nil {
			return "", err
		}
		if err := dockerClient.WaitUntilRunning(ctx, c.containerName, runningTimeout); err != nil {
			return "", err
		}
	}
	c.logger.Debug("container is running, check if it's ready...")
	cli, err := client.New(ctx, c.BuildkitdAddr(), client.WithFailFast())
	if err != nil {
		return "", errors.Wrap(err, "failed to create the buildkit client")
	}
	c.Client = cli

	if err := c.waitUntilConnected(ctx, runningTimeout); err != nil {
		return "", errors.Wrap(err, "failed to connect to buildkitd")
	}

	return c.BuildkitdAddr(), nil
}

func (c generalClient) waitUntilConnected(
	ctx context.Context, timeout time.Duration) error {
	c.logger.Debug("waiting to connect to buildkitd")

	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-time.After(interval):
			connected, err := c.connected(ctxTimeout)
			if err != nil {
				// Has not yet started. Keep waiting.
				return errors.Wrap(err, "failed to connect to buildkitd")
			}
			if !connected {
				continue
			}
			if connected {
				c.logger.Debug("connected to buildkitd")
				return nil
			}

		case <-ctxTimeout.Done():
			return errors.Errorf("timeout %s: cannot connect to buildkitd", timeout)
		}
	}
}

func (c generalClient) connected(ctx context.Context) (bool, error) {
	if _, err := c.ListWorkers(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func (c generalClient) BuildkitdAddr() string {
	return fmt.Sprintf("docker-container://%s", c.containerName)
}
