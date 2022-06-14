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

package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/moby/term"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

const (
	localhost = "127.0.0.1"
)

var (
	interval = 1 * time.Second
)

type Client interface {
	// Load loads the image from the reader to the docker host.
	Load(ctx context.Context, r io.ReadCloser, quiet bool) error
	// Start creates the container for the given tag and container name.
	StartEnvd(ctx context.Context, tag, name, buildContext string,
		gpuEnabled bool, numGPUs int, sshPort int, g ir.Graph, timeout time.Duration,
		mountOptionsStr []string) (string, string, error)
	StartBuildkitd(ctx context.Context, tag, name, mirror string) (string, error)

	IsRunning(ctx context.Context, name string) (bool, error)
	IsCreated(ctx context.Context, name string) (bool, error)
	WaitUntilRunning(ctx context.Context, name string, timeout time.Duration) error

	Exec(ctx context.Context, cname string, cmd []string) error
	Destroy(ctx context.Context, name string) (string, error)

	ListContainer(ctx context.Context) ([]types.Container, error)
	GetContainer(ctx context.Context, cname string) (types.ContainerJSON, error)
	PauseContainer(ctx context.Context, name string) (string, error)
	ResumeContainer(ctx context.Context, name string) (string, error)

	ListImage(ctx context.Context) ([]types.ImageSummary, error)
	GetImage(ctx context.Context, image string) (types.ImageSummary, error)

	// GPUEnabled returns true if nvidia container runtime exists in docker daemon.
	GPUEnabled(ctx context.Context) (bool, error)
}

type generalClient struct {
	*client.Client
}

func NewClient(ctx context.Context) (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)
	return generalClient{cli}, nil
}

func (g generalClient) GPUEnabled(ctx context.Context) (bool, error) {
	info, err := g.Info(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get docker info")
	}
	logrus.WithField("info", info).Debug("docker info")
	if nv, ok := info.Runtimes["nvidia"]; ok {
		return nv.Path != "", nil
	}
	return false, nil
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
			return errors.Errorf("timeout %s: container did not start", timeout)
		}
	}
}

func (c generalClient) ListImage(ctx context.Context) ([]types.ImageSummary, error) {
	images, err := c.ImageList(ctx, types.ImageListOptions{
		Filters: dockerfilters(false),
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (c generalClient) GetImage(
	ctx context.Context, image string) (types.ImageSummary, error) {
	images, err := c.ImageList(ctx, types.ImageListOptions{
		Filters: dockerfiltersWithName(image),
	})
	if err != nil {
		return types.ImageSummary{}, err
	}
	if len(images) == 0 {
		return types.ImageSummary{}, errors.Errorf("image %s not found", image)
	}
	return images[0], nil
}

func (c generalClient) ListContainer(ctx context.Context) ([]types.Container, error) {
	ctrs, err := c.ContainerList(ctx, types.ContainerListOptions{
		Filters: dockerfilters(false),
	})
	if err != nil {
		return nil, err
	}
	return ctrs, nil
}

func (c generalClient) PauseContainer(ctx context.Context, name string) (string, error) {
	logger := logrus.WithField("container", name)
	err := c.ContainerPause(ctx, name)
	if err != nil {
		errCause := errors.UnwrapAll(err).Error()
		switch {
		case strings.Contains(errCause, "is already paused"):
			logger.Debug("container is already paused, there is no need to pause it again")
			return "", nil
		case strings.Contains(errCause, "No such container"):
			logger.Debug("container is not found, there is no need to pause it")
			return "", errors.New("container not found")
		default:
			return "", errors.Wrap(err, "failed to pause container")
		}
	}
	return name, nil
}

func (c generalClient) ResumeContainer(ctx context.Context, name string) (string, error) {
	logger := logrus.WithField("container", name)
	err := c.ContainerUnpause(ctx, name)
	if err != nil {
		errCause := errors.UnwrapAll(err).Error()
		switch {
		case strings.Contains(errCause, "is not paused"):
			logger.Debug("container is not paused, there is no need to resume")
			return "", nil
		case strings.Contains(errCause, "No such container"):
			logger.Debug("container is not found, there is no need to resume it")
			return "", errors.New("container not found")
		default:
			return "", errors.Wrap(err, "failed to resume container")
		}
	}
	return name, nil
}

func (c generalClient) GetContainer(ctx context.Context, cname string) (types.ContainerJSON, error) {
	return c.ContainerInspect(ctx, cname)
}

func (c generalClient) Destroy(ctx context.Context, name string) (string, error) {
	logger := logrus.WithField("container", name)
	// Refer to https://docs.docker.com/engine/reference/commandline/container_kill/
	if err := c.ContainerKill(ctx, name, "KILL"); err != nil {
		errCause := errors.UnwrapAll(err).Error()
		switch {
		case strings.Contains(errCause, "is not running"):
			// If the container is not running, there is no need to kill it.
			logger.Debug("container is not running, there is no need to kill it")
		case strings.Contains(errCause, "No such container"):
			// If the container is not found, it is already destroyed.
			logger.Debug("container is not found, there is no need to destroy it")
			return "", nil
		default:
			return "", errors.Wrap(err, "failed to kill the container")
		}
	}

	if err := c.ContainerRemove(ctx, name, types.ContainerRemoveOptions{}); err != nil {
		return "", errors.Wrap(err, "failed to remove the container")
	}
	return name, nil
}

func (g generalClient) StartBuildkitd(ctx context.Context,
	tag, name, mirror string) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":       tag,
		"container": name,
		"mirror":    mirror,
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
			defer body.Close()
			termFd, isTerm := term.GetFdInfo(os.Stdout)
			err = jsonmessage.DisplayJSONMessagesStream(body, os.Stdout, termFd, isTerm, nil)
			if err != nil {
				logger.WithError(err).Warningln("failed to display image pull output")
			}
		} else {
			return "", errors.Wrap(err, "failed to inspect image")
		}
	}
	config := &container.Config{
		Image: tag,
	}
	if mirror != "" {
		cfg := fmt.Sprintf(`
[registry."docker.io"]
	mirrors = ["%s"]`, mirror)
		config.Entrypoint = []string{
			"/bin/sh", "-c", fmt.Sprintf("mkdir /etc/buildkit && echo '%s' > /etc/buildkit/buildkitd.toml && buildkitd", cfg),
		}
		logger.Debugf("setting buildkit config: %s", cfg)
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
func (c generalClient) StartEnvd(ctx context.Context, tag, name, buildContext string,
	gpuEnabled bool, numGPUs int, sshPort int, g ir.Graph, timeout time.Duration, mountOptionsStr []string) (string, string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":           tag,
		"container":     name,
		"gpu":           gpuEnabled,
		"numGPUs":       numGPUs,
		"build-context": buildContext,
	})
	config := &container.Config{
		Image: tag,
		User:  "envd",
		Entrypoint: []string{
			"tini",
			"--",
			"bash",
			"-c",
		},
		ExposedPorts: nat.PortSet{},
	}
	base := fileutil.Base(buildContext)
	base = filepath.Join("/home/envd", base)
	config.WorkingDir = base
	config.Entrypoint = append(config.Entrypoint,
		entrypointSH(g, config.WorkingDir, sshPort))

	mountOption := make([]mount.Mount, len(mountOptionsStr)+1)
	for i, option := range mountOptionsStr {
		mStr := strings.Split(option, ":")
		if len(mStr) != 2 {
			return "", "", errors.Newf("Invalid mount options %s", option)
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
		Source: buildContext,
		Target: base,
	}

	logger.WithFields(logrus.Fields{
		"mount-path":  buildContext,
		"working-dir": base,
	}).Debug("setting up container working directory")

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{},
		Mounts:       mountOption,
	}

	// Configure ssh port.
	natPort := nat.Port(fmt.Sprintf("%d/tcp", sshPort))
	hostConfig.PortBindings[natPort] = []nat.PortBinding{
		{
			HostIP:   localhost,
			HostPort: strconv.Itoa(sshPort),
		},
	}
	config.ExposedPorts[natPort] = struct{}{}

	// TODO(gaocegege): Avoid specific logic to set the port.
	if g.JupyterConfig != nil {
		natPort := nat.Port(fmt.Sprintf("%d/tcp", g.JupyterConfig.Port))
		hostConfig.PortBindings[natPort] = []nat.PortBinding{
			{
				HostIP:   localhost,
				HostPort: strconv.Itoa(int(g.JupyterConfig.Port)),
			},
		}
		config.ExposedPorts[natPort] = struct{}{}
	}

	if gpuEnabled {
		logger.Debug("GPU is enabled.")
		hostConfig.DeviceRequests = deviceRequests(numGPUs)
	}

	config.Labels = labels(name, g.JupyterConfig, sshPort)

	logger = logger.WithFields(logrus.Fields{
		"entrypoint":  config.Entrypoint,
		"working-dir": config.WorkingDir,
	})
	logger.Debugf("starting %s container", name)

	resp, err := c.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create the container")
	}

	for _, w := range resp.Warnings {
		logger.Warnf("run with warnings: %s", w)
	}

	if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		errCause := errors.UnwrapAll(err)
		// Hack to check if the port is already allocated.
		if strings.Contains(errCause.Error(), "port is already allocated") {
			logrus.Debugf("failed to allocate the port: %s", err)
			return "", "", errors.New("jupyter port is already allocated in the host")
		}
		return "", "", errors.Wrap(err, "failed to run the container")
	}

	container, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to inpsect the container")
	}

	if err := c.WaitUntilRunning(
		ctx, container.Name, timeout); err != nil {
		return "", "", errors.Wrap(err, "failed to wait until the container is running")
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
