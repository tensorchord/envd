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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/driver"
)

var (
	anchoredIdentifierRegexp = regexp.MustCompile(`^([a-f0-9]{64})$`)
	waitingInterval          = 1 * time.Second
)

type dockerClient struct {
	*client.Client
}

func NewClient(ctx context.Context) (driver.Client, error) {
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
	return dockerClient{cli}, nil
}

// Normalize the name accord the spec of docker, It may support normalize imagea and container in the future.
func NormalizeName(s string) (string, error) {
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
		logrus.Warnf("The working directory's name is not lowercased: %s, the image built will be lowercased to %s", remoteName, s)
	}
	// remove the spaces
	s = strings.ReplaceAll(s, " ", "")
	return s, nil
}

func (c dockerClient) ListImage(ctx context.Context) ([]types.ImageSummary, error) {
	images, err := c.ImageList(ctx, types.ImageListOptions{
		Filters: dockerFilters(false),
	})
	return images, err
}

func (c dockerClient) RemoveImage(ctx context.Context, image string) error {
	_, err := c.ImageRemove(ctx, image, types.ImageRemoveOptions{})
	if err != nil {
		logrus.WithError(err).Errorf("failed to remove image %s", image)
		return err
	}
	return nil
}

func (c dockerClient) GetImage(ctx context.Context, image string) (types.ImageSummary, error) {
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

func (c dockerClient) GetImageWithCacheHashLabel(ctx context.Context, image string, hash string) (types.ImageSummary, error) {
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

func (c dockerClient) PauseContainer(ctx context.Context, name string) (string, error) {
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

func (c dockerClient) ResumeContainer(ctx context.Context, name string) (string, error) {
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

func (c dockerClient) StartBuildkitd(ctx context.Context,
	tag, name, mirror string, timeout time.Duration) (string, error) {
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

	container, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", errors.Wrap(err, "failed to inspect container")
	}

	err = c.waitUntilRunning(ctx, container.Name, timeout)
	if err != nil {
		return "", err
	}

	return container.Name, nil
}

func (c dockerClient) Exists(ctx context.Context, cname string) (bool, error) {
	_, err := c.ContainerInspect(ctx, cname)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c dockerClient) IsRunning(ctx context.Context, cname string) (bool, error) {
	container, err := c.ContainerInspect(ctx, cname)
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
func (c dockerClient) Load(ctx context.Context, r io.ReadCloser, quiet bool) error {
	resp, err := c.ImageLoad(ctx, r, quiet)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (c dockerClient) Exec(ctx context.Context, cname string, cmd []string) error {
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

func (c dockerClient) PruneImage(ctx context.Context) (types.ImagesPruneReport, error) {
	pruneReport, err := c.ImagesPrune(ctx, filters.Args{})
	if err != nil {
		return types.ImagesPruneReport{}, errors.Wrap(err, "failed to prune images")
	}
	return pruneReport, nil
}

func (c dockerClient) Stats(ctx context.Context, cname string, statChan chan<- *driver.Stats, done <-chan bool) (retErr error) {
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
	stats := new(driver.Stats)
	for err := decoder.Decode(stats); !errors.Is(err, io.EOF); err = decoder.Decode(stats) {
		if err != nil {
			return err
		}
		statChan <- stats
		stats = new(driver.Stats)
	}
	return nil
}

func (c dockerClient) waitUntilRunning(ctx context.Context,
	name string, timeout time.Duration) error {
	logger := logrus.WithField("container", name)
	logger.Debug("waiting to start")

	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-time.After(waitingInterval):
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
			container, err := c.ContainerInspect(ctx, name)
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

func GetDockerVersion() (int, error) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return -1, err
	}
	defer cli.Close()

	info, err := cli.Info(ctx)
	if err != nil {
		return -1, err
	}
	version, err := strconv.Atoi(strings.Split(info.ServerVersion, ".")[0])
	if err != nil {
		return -1, err
	}
	return version, nil
}
