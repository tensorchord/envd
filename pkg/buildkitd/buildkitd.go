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

package buildkitd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	gateway "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tonistiigi/units"

	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/types"
)

var (
	interval          = time.Second * 1
	timeoutConnection = time.Second * 5
	timeoutRun        = time.Second * 3
)

// Client is a client for the buildkitd daemon.
// It's up to the caller to close the client.
type Client interface {
	BuildkitdAddr() string
	// Solve calls Solve on the controller.
	Solve(ctx context.Context, def *llb.Definition,
		opt client.SolveOpt, statusChan chan *client.SolveStatus) (*client.SolveResponse, error)
	Build(ctx context.Context, opt client.SolveOpt, product string,
		buildFunc gateway.BuildFunc, statusChan chan *client.SolveStatus,
	) (*client.SolveResponse, error)
	Prune(ctx context.Context, keepDuration time.Duration,
		keepStorage float64, filter []string, verbose, all bool) error
	Close() error
}

type generalClient struct {
	containerName string
	image         string
	mirror        string

	driver types.BuilderType
	socket string

	*client.Client
	logger *logrus.Entry
}

func NewClient(ctx context.Context, driver types.BuilderType,
	socket, mirror string) (Client, error) {
	c := &generalClient{
		containerName: socket,
		image:         viper.GetString(flag.FlagBuildkitdImage),
		mirror:        mirror,
	}
	c.socket = socket
	c.driver = driver
	c.logger = logrus.WithFields(logrus.Fields{
		"container": c.containerName,
		"image":     c.image,
		"socket":    c.socket,
		"driver":    c.driver,
	})

	cli, err := client.New(ctx, c.BuildkitdAddr(), client.WithFailFast())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the client")
	}
	c.Client = cli

	if _, err := c.Bootstrap(ctx, timeoutRun, timeoutConnection); err != nil {
		return nil, errors.Wrap(err, "failed to bootstrap the buildkitd")
	}
	return c, nil
}

func (c *generalClient) Bootstrap(ctx context.Context,
	runningTimeout, connectingTimeout time.Duration) (string, error) {
	address, err := c.maybeStart(ctx, runningTimeout, connectingTimeout)
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
	if c.driver == types.BuilderTypeDocker {
		dockerClient, err := docker.NewClient(ctx)
		if err != nil {
			return "", err
		}
		opt := envd.Options{
			Context: &types.Context{
				Runner: types.RunnerTypeDocker,
			},
		}
		engine, err := envd.New(ctx, opt)
		if err != nil {
			return "", err
		}

		created, err := engine.Exists(ctx, c.containerName)
		if err != nil {
			return "", err
		}

		if !created {
			c.logger.Debug("container not created, creating...")
			if _, err := dockerClient.StartBuildkitd(
				ctx, c.image, c.containerName, c.mirror); err != nil {
				return "", err
			}
			if err := engine.WaitUntilRunning(
				ctx, c.containerName, runningTimeout); err != nil {
				return "", err
			}
		}
		running, _ := engine.IsRunning(ctx, c.containerName)
		if created && !running {
			c.logger.Warnf("start the created contrainer %s", c.containerName)
			_, err := dockerClient.StartBuildkitd(ctx, c.image, c.containerName, c.mirror)
			if err != nil {
				c.logger.Warnf("please remove or restart the container %s", c.containerName)
				return "", errors.Errorf("container %s is stopped", c.containerName)
			}
		}

		c.logger.Debugf("container is running, check if it's ready at %s...", c.BuildkitdAddr())
	}

	if err := c.waitUntilConnected(ctx, connectingTimeout); err != nil {
		return "", errors.Wrapf(err,
			"failed to connect to buildkitd %s", c.BuildkitdAddr())
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
				logrus.Debugf("failed to connect to buildkitd: %s", err.Error())
				continue
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

func (c generalClient) Prune(ctx context.Context, keepDuration time.Duration,
	keepStorage float64, filter []string, verbose, all bool) error {
	opts := []client.PruneOption{
		client.WithFilter(filter),
		client.WithKeepOpt(keepDuration, int64(keepStorage*1e6)),
	}
	if all {
		opts = append(opts, client.PruneAll)
	}

	ch := make(chan client.UsageInfo)
	printed := make(chan struct{})

	tw := tabwriter.NewWriter(os.Stdout, 1, 8, 1, '\t', 0)
	first := true
	total := int64(0)

	go func() {
		defer close(printed)
		for du := range ch {
			total += du.Size
			if verbose {
				printVerbose(tw, []*client.UsageInfo{&du})
			} else {
				if first {
					printTableHeader(tw)
					first = false
				}
				printTableRow(tw, &du)
				tw.Flush()
			}
		}
	}()

	err := c.Client.Prune(ctx, ch, opts...)
	close(ch)
	<-printed
	if err != nil {
		return err
	}

	tw = tabwriter.NewWriter(os.Stdout, 1, 8, 1, '\t', 0)
	fmt.Fprintf(tw, "Total:\t%.2f\n", units.Bytes(total))
	tw.Flush()
	return nil
}

func (c generalClient) BuildkitdAddr() string {
	return fmt.Sprintf("%s://%s", c.driver, c.socket)
}
