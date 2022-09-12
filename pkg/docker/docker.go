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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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

	envdconfig "github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/util/fileutil"
	"github.com/tensorchord/envd/pkg/util/netutil"
)

const (
	localhost = "127.0.0.1"
)

var (
	interval                 = 1 * time.Second
	anchoredIdentifierRegexp = regexp.MustCompile(`^([a-f0-9]{64})$`)
)

type Client interface {
	// Load loads the image from the reader to the docker host.
	Load(ctx context.Context, r io.ReadCloser, quiet bool) error
	// StartEnvd creates the container for the given tag and container name.
	StartEnvd(ctx context.Context, tag, name, buildContext string,
		gpuEnabled bool, numGPUs int, sshPort int, g ir.Graph, timeout time.Duration,
		mountOptionsStr []string) (string, string, error)
	StartBuildkitd(ctx context.Context, tag, name, mirror string) (string, error)
	CleanEnvdIfExists(ctx context.Context, name string, force bool) error

	IsRunning(ctx context.Context, name string) (bool, error)
	Exists(ctx context.Context, name string) (bool, error)
	WaitUntilRunning(ctx context.Context, name string, timeout time.Duration) error

	Exec(ctx context.Context, cname string, cmd []string) error
	Destroy(ctx context.Context, name string) (string, error)

	ListContainer(ctx context.Context) ([]types.Container, error)
	GetContainer(ctx context.Context, cname string) (types.ContainerJSON, error)
	PauseContainer(ctx context.Context, name string) (string, error)
	ResumeContainer(ctx context.Context, name string) (string, error)

	ListImage(ctx context.Context) ([]types.ImageSummary, error)
	GetImage(ctx context.Context, image string) (types.ImageSummary, error)
	GetImageWithCacheHashLabel(ctx context.Context, image string, hash string) (types.ImageSummary, error)
	RemoveImage(ctx context.Context, image string) error

	GetInfo(ctx context.Context) (types.Info, error)
	Stats(ctx context.Context, cname string, statChan chan<- *Stats, done <-chan bool) error

	// GPUEnabled returns true if nvidia container runtime exists in docker daemon.
	GPUEnabled(ctx context.Context) (bool, error)
}

type generalClient struct {
	*client.Client
}

func NewClient(ctx context.Context) (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	_, err = cli.Ping(ctx)
	if err != nil {
		// Special note needed to give users
		if strings.Contains(err.Error(), "permission denied") {
			err = errors.New(`It seems that current user have no access to docker daemon,
please visit https://docs.docker.com/engine/install/linux-postinstall/ for more info.`)
		}
		return nil, err
	}
	return generalClient{cli}, nil
}

// Normalize the name accord the spec of docker, It may support normalize imagea and container in the future.
func NormalizeNamed(s string) (string, error) {
	if ok := anchoredIdentifierRegexp.MatchString(s); ok {
		return "", errors.Newf("invalid repository name (%s), cannot specify 64-byte hexadecimal strings, please rename it", s)
	}
	var remoteName string
	var tagSep int
	if tagSep = strings.IndexRune(s, ':'); tagSep > -1 {
		remoteName = s[:tagSep]
	} else {
		remoteName = s
	}
	if strings.ToLower(remoteName) != remoteName {
		remoteName = strings.ToLower(remoteName)
		if tagSep > -1 {
			s = remoteName + s[tagSep:]
		} else {
			s = remoteName
		}
		logrus.Warnf("The working direcotry's name is not lowercased: %s, the image built will be lowercased to %s", remoteName, s)
	}
	return s, nil

}

func (c generalClient) GPUEnabled(ctx context.Context) (bool, error) {
	info, err := c.GetInfo(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get docker info")
	}
	logrus.WithField("info", info).Debug("docker info")
	nv := info.Runtimes["nvidia"]
	return nv.Path != "", nil
}

func (c generalClient) WaitUntilRunning(ctx context.Context,
	name string, timeout time.Duration) error {
	logger := logrus.WithField("container", name)
	logger.Debug("waiting to start")

	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-time.After(interval):
			isRunning, err := c.IsRunning(ctxTimeout, name)
			if err != nil {
				// Has not yet started. Keep waiting.
				return errors.Wrap(err, "failed to check if container is running")
			}
			if isRunning {
				logger.Debug("the container is running")
				return nil
			}

		case <-ctxTimeout.Done():
			container, err := c.GetContainer(ctx, name)
			if err != nil {
				logger.Debugf("failed to inspect container %s", name)
			}
			state, err := json.Marshal(container.State)
			if err != nil {
				logger.Debug("failed to marshal container state")
			}
			logger.Debugf("container state: %s", state)
			return errors.Errorf("timeout %s: container did not start", timeout)
		}
	}
}

func (c generalClient) ListImage(ctx context.Context) ([]types.ImageSummary, error) {
	images, err := c.ImageList(ctx, types.ImageListOptions{
		Filters: dockerFilters(false),
	})
	return images, err
}

func (c generalClient) RemoveImage(ctx context.Context, image string) error {
	_, err := c.ImageRemove(ctx, image, types.ImageRemoveOptions{})
	if err != nil {
		logrus.WithError(err).Errorf("failed to remove image %s", image)
		return err
	}
	return nil
}

func (c generalClient) GetImage(ctx context.Context, image string) (types.ImageSummary, error) {
	images, err := c.ImageList(ctx, types.ImageListOptions{
		Filters: dockerFiltersWithName(image),
	})
	if err != nil {
		return types.ImageSummary{}, err
	}
	if len(images) == 0 {
		return types.ImageSummary{}, errors.Errorf("image %s not found", image)
	}
	return images[0], nil
}

func (c generalClient) GetImageWithCacheHashLabel(ctx context.Context, image string, hash string) (types.ImageSummary, error) {
	images, err := c.ImageList(ctx, types.ImageListOptions{
		Filters: dockerFiltersWithCacheLabel(image, hash),
	})
	if err != nil {
		return types.ImageSummary{}, err
	}
	if len(images) == 0 {
		return types.ImageSummary{}, errors.Errorf("image with hash %s not found", hash)
	}
	return images[0], nil
}

func (c generalClient) ListContainer(ctx context.Context) ([]types.Container, error) {
	return c.ContainerList(ctx, types.ContainerListOptions{
		Filters: dockerFilters(false),
	})
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

func (c generalClient) GetInfo(ctx context.Context) (types.Info, error) {
	return c.Info(ctx)
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
			// If the container is not found, it is already destroyed or the name is wrong.
			logger.Infof("cannot find container %s, maybe it's already destroyed or the name is wrong", name)
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

func (c generalClient) StartBuildkitd(ctx context.Context, tag, name, mirror string) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":       tag,
		"container": name,
		"mirror":    mirror,
	})
	logger.Debug("starting buildkitd")
	if _, _, err := c.ImageInspectWithRaw(ctx, tag); err != nil {
		if !client.IsErrNotFound(err) {
			return "", errors.Wrap(err, "failed to inspect image")
		}

		// Pull the image.
		logger.Debug("pulling image")
		body, err := c.ImagePull(ctx, tag, types.ImagePullOptions{})
		if err != nil {
			return "", errors.Wrap(err, "failed to pull image")
		}
		defer body.Close()
		termFd, isTerm := term.GetFdInfo(os.Stdout)
		err = jsonmessage.DisplayJSONMessagesStream(body, os.Stdout, termFd, isTerm, nil)
		if err != nil {
			logger.WithError(err).Warningln("failed to display image pull output")
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
			"/bin/sh",
			"-c",
			fmt.Sprintf("mkdir /etc/buildkit && echo '%s' > /etc/buildkit/buildkitd.toml && buildkitd", cfg),
		}
		logger.Debugf("setting buildkit config: %s", cfg)
	}
	hostConfig := &container.HostConfig{
		Privileged: true,
	}
	created, _ := c.Exists(ctx, name)
	if created {
		err := c.ContainerStart(ctx, name, types.ContainerStartOptions{})
		if err != nil {
			return name, err
		}
		return name, nil
	}
	resp, err := c.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", errors.Wrap(err, "failed to create container")
	}

	for _, w := range resp.Warnings {
		logger.Warnf("run with warnings: %s", w)
	}

	if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "failed to start container")
	}

	container, err := c.GetContainer(ctx, resp.ID)
	if err != nil {
		return "", errors.Wrap(err, "failed to inspect container")
	}

	return container.Name, nil
}

func (c generalClient) CleanEnvdIfExists(ctx context.Context, name string, force bool) error {
	created, err := c.Exists(ctx, name)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}

	// force delete the container no matter it is running or not.
	if force {
		return c.ContainerRemove(ctx, name, types.ContainerRemoveOptions{
			Force: true,
		})
	}

	running, _ := c.IsRunning(ctx, name)
	if err != nil {
		return err
	}
	if running {
		logrus.Errorf("container %s is running, cannot clean envd, please save your data and stop the running container if you need to envd up again.", name)
		return errors.Newf("\"%s\" is stil running, please run `envd destroy --name %s` stop it first", name, name)
	}
	return c.ContainerRemove(ctx, name, types.ContainerRemoveOptions{})
}

// StartEnvd creates the container for the given tag and container name.
func (c generalClient) StartEnvd(ctx context.Context, tag, name, buildContext string,
	gpuEnabled bool, numGPUs int, sshPortInHost int, g ir.Graph, timeout time.Duration,
	mountOptionsStr []string) (string, string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":           tag,
		"container":     name,
		"gpu":           gpuEnabled,
		"numGPUs":       numGPUs,
		"build-context": buildContext,
	})
	config := &container.Config{
		Image:        tag,
		User:         "envd",
		ExposedPorts: nat.PortSet{},
	}
	base := fileutil.EnvdHomeDir(filepath.Base(buildContext))
	config.WorkingDir = base

	mountOption := make([]mount.Mount, 0, len(mountOptionsStr)+len(g.Mount)+1)
	for _, option := range mountOptionsStr {
		mStr := strings.Split(option, ":")
		if len(mStr) != 2 {
			return "", "", errors.Newf("Invalid mount options %s", option)
		}

		logger.WithFields(logrus.Fields{
			"mount-path":     mStr[0],
			"container-path": mStr[1],
		}).Debug("setting up container working directory")
		mountOption = append(mountOption, mount.Mount{
			Type:   mount.TypeBind,
			Source: mStr[0],
			Target: mStr[1],
		})
	}
	for _, m := range g.Mount {
		logger.WithFields(logrus.Fields{
			"mount-path":     m.Source,
			"container-path": m.Destination,
		}).Debug("setting up declared mount directory")
		mountOption = append(mountOption, mount.Mount{
			Type:   mount.TypeBind,
			Source: m.Source,
			Target: m.Destination,
		})
	}

	mountOption = append(mountOption, mount.Mount{
		Type:   mount.TypeBind,
		Source: buildContext,
		Target: base,
	})

	logger.WithFields(logrus.Fields{
		"mount-path":  buildContext,
		"working-dir": base,
	}).Debug("setting up container working directory")

	rp := container.RestartPolicy{
		Name: "always",
	}
	hostConfig := &container.HostConfig{
		PortBindings:  nat.PortMap{},
		Mounts:        mountOption,
		RestartPolicy: rp,
	}

	// Configure ssh port.
	natPort := nat.Port(fmt.Sprintf("%d/tcp", envdconfig.SSHPortInContainer))
	hostConfig.PortBindings[natPort] = []nat.PortBinding{
		{
			HostIP:   localhost,
			HostPort: strconv.Itoa(sshPortInHost),
		},
	}

	var jupyterPortInHost int
	// TODO(gaocegege): Avoid specific logic to set the port.
	if g.JupyterConfig != nil {
		if g.JupyterConfig.Port != 0 {
			jupyterPortInHost = int(g.JupyterConfig.Port)
		} else {
			var err error
			jupyterPortInHost, err = netutil.GetFreePort()
			if err != nil {
				return "", "", errors.Wrap(err, "failed to get a free port")
			}
		}
		natPort := nat.Port(fmt.Sprintf("%d/tcp", envdconfig.JupyterPortInContainer))
		hostConfig.PortBindings[natPort] = []nat.PortBinding{
			{
				HostIP:   localhost,
				HostPort: strconv.Itoa(jupyterPortInHost),
			},
		}
		config.ExposedPorts[natPort] = struct{}{}
	}
	var rStudioPortInHost int
	if g.RStudioServerConfig != nil {
		var err error
		rStudioPortInHost, err = netutil.GetFreePort()
		if err != nil {
			return "", "", errors.Wrap(err, "failed to get a free port")
		}
		natPort := nat.Port(fmt.Sprintf("%d/tcp", envdconfig.RStudioServerPortInContainer))
		hostConfig.PortBindings[natPort] = []nat.PortBinding{
			{
				HostIP:   localhost,
				HostPort: strconv.Itoa(rStudioPortInHost),
			},
		}
		config.ExposedPorts[natPort] = struct{}{}
	}

	if len(g.RuntimeExpose) > 0 {
		for _, item := range g.RuntimeExpose {
			var err error
			if item.HostPort == 0 {
				item.HostPort, err = netutil.GetFreePort()
				if err != nil {
					return "", "", errors.Wrap(err, "failed to get a free port")
				}
			}
			natPort := nat.Port(fmt.Sprintf("%d/tcp", item.EnvdPort))
			hostConfig.PortBindings[natPort] = []nat.PortBinding{
				{
					HostIP:   localhost,
					HostPort: strconv.Itoa(item.HostPort),
				},
			}
			config.ExposedPorts[natPort] = struct{}{}
		}
	}

	if gpuEnabled {
		logger.Debug("GPU is enabled.")
		hostConfig.DeviceRequests = deviceRequests(numGPUs)
	}

	config.Labels = labels(name, g,
		sshPortInHost, jupyterPortInHost, rStudioPortInHost)

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

	if err := c.ContainerStart(
		ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		errCause := errors.UnwrapAll(err)
		// Hack to check if the port is already allocated.
		if strings.Contains(errCause.Error(), "port is already allocated") {
			logrus.Debugf("failed to allocate the port: %s", err)
			return "", "", errors.New("port is already allocated in the host")
		}
		return "", "", errors.Wrap(err, "failed to run the container")
	}

	container, err := c.GetContainer(ctx, resp.ID)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to inspect the container")
	}

	if err := c.WaitUntilRunning(
		ctx, container.Name, timeout); err != nil {
		return "", "", errors.Wrap(err, "failed to wait until the container is running")
	}

	return container.Name, container.NetworkSettings.IPAddress, nil
}

func (c generalClient) Exists(ctx context.Context, cname string) (bool, error) {
	_, err := c.GetContainer(ctx, cname)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c generalClient) IsRunning(ctx context.Context, cname string) (bool, error) {
	container, err := c.GetContainer(ctx, cname)
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
func (c generalClient) Load(ctx context.Context, r io.ReadCloser, quiet bool) error {
	resp, err := c.ImageLoad(ctx, r, quiet)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (c generalClient) Exec(ctx context.Context, cname string, cmd []string) error {
	execConfig := types.ExecConfig{
		Cmd:    cmd,
		Detach: true,
	}
	resp, err := c.ContainerExecCreate(ctx, cname, execConfig)
	if err != nil {
		return err
	}
	execID := resp.ID
	return c.ContainerExecStart(ctx, execID, types.ExecStartCheck{
		Detach: true,
	})
}

func (c generalClient) Stats(ctx context.Context, cname string, statChan chan<- *Stats, done <-chan bool) (retErr error) {
	errC := make(chan error, 1)
	containerStats, err := c.ContainerStats(ctx, cname, true)
	readCloser := containerStats.Body
	quit := make(chan struct{})
	defer func() {
		close(statChan)
		close(quit)

		if err := <-errC; err != nil && retErr == nil {
			retErr = err
		}

		if err := readCloser.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()

	go func() {
		// block here waiting for the signal to stop function
		select {
		case <-done:
			readCloser.Close()
		case <-quit:
			return
		}
	}()

	if err != nil {
		return err
	}
	decoder := json.NewDecoder(readCloser)
	stats := new(Stats)
	for err := decoder.Decode(stats); !errors.Is(err, io.EOF); err = decoder.Decode(stats) {
		if err != nil {
			return err
		}
		statChan <- stats
		stats = new(Stats)
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
