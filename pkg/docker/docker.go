// Copyright 2022 The MIDI Authors
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

package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type Client interface {
	// Load loads the image from the reader to the docker host.
	Load(ctx context.Context, r io.ReadCloser, quiet bool) error
	// Start creates the container for the given tag and container name.
	Start(ctx context.Context, tag, name string, gpuEnabled bool) (string, string, error)
}

type generalClient struct {
	*client.Client
}

func NewClient() (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return generalClient{cli}, nil
}

// Start creates the container for the given tag and container name.
func (g generalClient) Start(
	ctx context.Context, tag, name string, gpuEnabled bool) (string, string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":       tag,
		"container": name,
		"gpu":       gpuEnabled,
	})
	config := &container.Config{
		Image: tag,
		Entrypoint: []string{
			"/var/midi/bin/midi-ssh",
			"--no-auth",
		},
	}
	hostConfig := &container.HostConfig{}
	if gpuEnabled {
		logger.Debug("GPU is enabled.")
		// enable all gpus with -1
		hostConfig.DeviceRequests = deviceRequests(-1)
	}
	resp, err := g.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", "", err
	}

	for _, w := range resp.Warnings {
		logger.Warnf("run with warnings: %s", w)
	}

	if err := g.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", "", err
	}

	container, err := g.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", "", err
	}

	return container.Name, container.NetworkSettings.IPAddress, nil
}

// Load loads the docker image from the reader into the docker host.
// It's up to the caller to close the io.ReadCloser.
func (g generalClient) Load(ctx context.Context, r io.ReadCloser, quiet bool) error {
	resp, err := g.ImageLoad(ctx, r, quiet)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func deviceRequests(count int) []container.DeviceRequest {
	return []container.DeviceRequest{
		{
			Driver: "nvidia",
			Capabilities: [][]string{
				{"gpu"},
				{"nvidia"},
				{"compute"},
				{"compat32"},
				{"graphics"},
				{"utility"},
				{"video"},
				{"display"},
			},
			Count: count,
		},
	}
}
