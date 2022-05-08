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
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/MIDI/pkg/editor/jupyter"
	"github.com/tensorchord/MIDI/pkg/lang/ir"
	"github.com/tensorchord/MIDI/pkg/util/fileutil"
)

var (
	interval = 1 * time.Second
)

type Client interface {
	// Load loads the image from the reader to the docker host.
	Load(ctx context.Context, r io.ReadCloser, quiet bool) error
	// Start creates the container for the given tag and container name.
	StartMIDI(ctx context.Context, tag, name string,
		gpuEnabled bool, g ir.Graph, timeout time.Duration, mountOptionsStr []string) (string, string, error)
	StartBuildkitd(ctx context.Context, tag, name string) (string, error)
	IsRunning(ctx context.Context, name string) (bool, error)
	IsCreated(ctx context.Context, name string) (bool, error)
	WaitUntilRunning(ctx context.Context, name string, timeout time.Duration) error
	Exec(ctx context.Context, cname string, cmd []string) error
	Destroy(ctx context.Context, name string) error
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

func (g generalClient) WaitUntilRunning(ctx context.Context,
	name string, timeout time.Duration) error {
	logger := logrus.WithField("container", name)
	logger.Debug("waiting to start")

	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-time.After(interval):
			isRunning, err := g.IsRunning(ctxTimeout, name)
			if err != nil {
				// Has not yet started. Keep waiting.
				return errors.Wrap(err, "failed to check if container is running")
			}
			if !isRunning {
				continue
			}
			if isRunning {
				logger.Debug("the container is running")
				return nil
			}

		case <-ctxTimeout.Done():
			return errors.Errorf("timeout %s: buildkitd container did not start", timeout)
		}
	}
}

func (c generalClient) Destroy(ctx context.Context, name string) error {
	// Refer to https://docs.docker.com/engine/reference/commandline/container_kill/
	if err := c.ContainerKill(ctx, name, "KILL"); err != nil {
		return errors.Wrap(err, "failed to kill the container")
	}
	if err := c.ContainerRemove(ctx, name, types.ContainerRemoveOptions{}); err != nil {
		return errors.Wrap(err, "failed to remove the container")
	}
	return nil
}

func (g generalClient) StartBuildkitd(ctx context.Context,
	tag, name string) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":       tag,
		"container": name,
	})
	logger.Debug("starting buildkitd")
	if _, _, err := g.ImageInspectWithRaw(ctx, tag); err != nil {
		if client.IsErrNotFound(err) {
			// Pull the image.
			logger.Debug("pulling image")
			body, err := g.ImagePull(ctx, tag, types.ImagePullOptions{})
			if err != nil {
				return "", errors.Wrap(err, "failed to pull image")
			}
			_, err = io.Copy(os.Stdout, body)
			if err != nil {
				logger.WithError(err).Warningln("failed to copy image pull output")
			}
			defer body.Close()
		} else {
			return "", errors.Wrap(err, "failed to inspect image")
		}
	}
	config := &container.Config{
		Image: tag,
	}
	hostConfig := &container.HostConfig{
		Privileged: true,
	}
	resp, err := g.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", errors.Wrap(err, "failed to create container")
	}

	for _, w := range resp.Warnings {
		logger.Warnf("run with warnings: %s", w)
	}

	if err := g.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "failed to start container")
	}

	container, err := g.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", errors.Wrap(err, "failed to inspect container")
	}

	return container.Name, nil
}

// Start creates the container for the given tag and container name.
func (c generalClient) StartMIDI(ctx context.Context, tag, name string,
	gpuEnabled bool, g ir.Graph, timeout time.Duration, mountOptionsStr []string) (string, string, error) {
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
		ExposedPorts: nat.PortSet{},
	}
	path, err := fileutil.CWD()
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get current working directory")
	}
	base, err := fileutil.RootDir()
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get root directory")
	}
	base = fmt.Sprintf("/root/%s", base)
	config.WorkingDir = base

	mountOption := make([]mount.Mount, len(mountOptionsStr)+1)
	for i, option := range mountOptionsStr {
		mStr := strings.Split(option, ":")
		if len(mStr) != 2 {
			return "", "", errors.Wrap(err, fmt.Sprintf("Invalid mount options %s", option))
		}

		logger.WithFields(logrus.Fields{
			"mount-path":     mStr[0],
			"container-path": mStr[1],
		}).Debug("setting up container working directory")
		mountOption[i] = mount.Mount{
			Type:   mount.TypeBind,
			Source: mStr[0],
			Target: mStr[1],
		}
	}
	mountOption[len(mountOptionsStr)] = mount.Mount{
		Type:   mount.TypeBind,
		Source: path,
		Target: base,
	}

	logger.WithFields(logrus.Fields{
		"mount-path":  path,
		"working-dir": base,
	}).Debug("setting up container working directory")

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{},
		Mounts:       mountOption,
	}
	// TODO(gaocegege): Avoid specific logic to set the port.
	if g.JupyterConfig != nil {
		natPort := nat.Port(fmt.Sprintf("%d/tcp", g.JupyterConfig.Port))
		hostConfig.PortBindings[natPort] = []nat.PortBinding{
			{
				HostIP:   "localhost",
				HostPort: strconv.Itoa(int(g.JupyterConfig.Port)),
			},
		}
		config.ExposedPorts[natPort] = struct{}{}
	}
	if gpuEnabled {
		logger.Debug("GPU is enabled.")
		// enable all gpus with -1
		hostConfig.DeviceRequests = deviceRequests(-1)
	}
	resp, err := c.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", "", err
	}

	for _, w := range resp.Warnings {
		logger.Warnf("run with warnings: %s", w)
	}

	if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", "", err
	}

	container, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", "", err
	}

	if err := c.WaitUntilRunning(
		ctx, container.Name, timeout); err != nil {
		return "", "", errors.Wrap(err, "failed to wait until the container is running")
	}

	if g.JupyterConfig != nil {
		cmd, err := jupyter.GenerateCommand(*ir.DefaultGraph, config.WorkingDir)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to generate jupyter command")
		}
		logger.WithField("cmd", cmd).Debug("configuring jupyter")
		err = c.Exec(ctx, container.Name, cmd)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to exec the command")
		}
	}
	return container.Name, container.NetworkSettings.IPAddress, nil
}

func (g generalClient) IsCreated(ctx context.Context, cname string) (bool, error) {
	_, err := g.ContainerInspect(ctx, cname)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (g generalClient) IsRunning(ctx context.Context, cname string) (bool, error) {
	container, err := g.ContainerInspect(ctx, cname)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return container.State.Running, nil
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

func (g generalClient) Exec(ctx context.Context, cname string, cmd []string) error {
	execConfig := types.ExecConfig{
		Cmd:    cmd,
		Detach: true,
	}
	resp, err := g.ContainerExecCreate(ctx, cname, execConfig)
	if err != nil {
		return err
	}
	execID := resp.ID
	err = g.ContainerExecStart(ctx, execID, types.ExecStartCheck{
		Detach: true,
	})
	if err != nil {
		return err
	}
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
