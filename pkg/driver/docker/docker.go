// Copyright 2023 The envd Authors
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
	"bytes"
	"context"
	"encoding/base64"
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
	"github.com/containerd/errdefs"
	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/pkg/docker/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockerimage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	imagespec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/driver"
	"github.com/tensorchord/envd/pkg/envd"
	containerType "github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/buildkitutil"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

const buildkitdConfigPath = "/etc/registry"

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

// Normalize the name accord the spec of docker, It may support normalize image and container in the future.
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
	name, err := reference.Parse(s)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse the name '%s', please provide a valid image name", s)
	}
	return name.String(), nil
}

func (c dockerClient) ListImage(ctx context.Context) ([]dockerimage.Summary, error) {
	images, err := c.ImageList(ctx, dockerimage.ListOptions{
		Filters: dockerFilters(false),
	})
	return images, err
}

func (c dockerClient) RemoveImage(ctx context.Context, image string) error {
	_, err := c.ImageRemove(ctx, image, dockerimage.RemoveOptions{})
	if err != nil {
		logrus.WithError(err).Errorf("failed to remove image %s", image)
		return err
	}
	return nil
}

func (c dockerClient) PushImage(ctx context.Context, image string, platform string) error {
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return errors.Wrap(err, "failed to normalize the image name")
	}
	auth, err := config.GetCredentialsForRef(nil, ref)
	if err != nil {
		return errors.Wrap(err, "failed to get credentials for image")
	}
	buf, err := json.Marshal(auth)
	if err != nil {
		return errors.Wrap(err, "failed to marshal auth struct")
	}
	platformInfo := strings.Split(platform, "/")
	if len(platformInfo) != 2 {
		return errors.New("invalid platform format, should be <architecture>/<os>")
	}
	reader, err := c.ImagePush(ctx, image, dockerimage.PushOptions{
		RegistryAuth: base64.URLEncoding.EncodeToString(buf),
		Platform: &imagespec.Platform{
			Architecture: platformInfo[0],
			OS:           platformInfo[1],
		},
	})
	if err != nil {
		logrus.WithError(err).Errorf("failed to push image %s", image)
		return err
	}

	bar := envd.InitProgressBar(0)

	defer func() {
		reader.Close()
		bar.Finish()
	}()

	decoder := json.NewDecoder(reader)
	stats := new(jsonmessage.JSONMessage)
	for err := decoder.Decode(stats); !errors.Is(err, io.EOF); err = decoder.Decode(stats) {
		if err != nil {
			return err
		}
		if stats.Error != nil {
			return stats.Error
		}

		if stats.Status != "" {
			if stats.ID == "" {
				bar.UpdateTitle(stats.Status)
			} else {
				bar.UpdateTitle(fmt.Sprintf("Pushing image => [%s] %s %s", stats.ID, stats.Status, stats.Progress))
			}
		}

		stats = new(jsonmessage.JSONMessage)
	}
	return nil
}

func (c dockerClient) GetImage(ctx context.Context, image string) (dockerimage.Summary, error) {
	images, err := c.ImageList(ctx, dockerimage.ListOptions{
		Filters: dockerFiltersWithName(image),
	})
	if err != nil {
		return dockerimage.Summary{}, err
	}
	if len(images) == 0 {
		return dockerimage.Summary{}, errors.Errorf("image %s not found", image)
	}
	return images[0], nil
}

func (c dockerClient) GetImageWithCacheHashLabel(ctx context.Context, image string, hash string) (dockerimage.Summary, error) {
	images, err := c.ImageList(ctx, dockerimage.ListOptions{
		Filters: dockerFiltersWithCacheLabel(image, hash),
	})
	if err != nil {
		return dockerimage.Summary{}, err
	}
	if len(images) == 0 {
		return dockerimage.Summary{}, errors.Errorf("image with hash %s not found", hash)
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

func (c dockerClient) RemoveContainer(ctx context.Context, name string) (string, error) {
	logger := logrus.WithField("container", name)
	err := c.ContainerRemove(ctx, name, container.RemoveOptions{})
	if err != nil {
		errCause := errors.UnwrapAll(err).Error()
		switch {
		case strings.Contains(errCause, "No such container"):
			logger.Debug("container is not found, there is no need to remove it")
			return "", errors.New("container not found")
		default:
			return "", errors.Wrap(err, "failed to remove container")
		}
	}
	return name, nil
}

func (c dockerClient) StartBuildkitd(ctx context.Context, tag, name string, bc *buildkitutil.BuildkitConfig, timeout time.Duration) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":             tag,
		"container":       name,
		"buildkit-config": bc,
	})
	logger.Debug("starting buildkitd")
	var buf bytes.Buffer
	if _, err := c.ImageInspect(ctx, tag, client.ImageInspectWithRawResponse(&buf)); err != nil {
		if !errdefs.IsNotFound(err) {
			return "", errors.Wrap(err, "failed to inspect image")
		}

		// Pull the image.
		logger.Debug("pulling image")
		body, err := c.ImagePull(ctx, tag, dockerimage.PullOptions{})
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
	hostConfig := &container.HostConfig{
		Privileged: true,
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fileutil.DefaultConfigDir,
				Target: buildkitdConfigPath,
			},
		},
	}

	err := bc.Save()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate buildkit config")
	}
	config.Entrypoint = []string{
		"buildkitd", "--config", filepath.Join(buildkitdConfigPath, "buildkitd.toml"),
	}
	created, _ := c.Exists(ctx, name)
	if created {
		status, err := c.GetStatus(ctx, name)
		if err != nil {
			return name, errors.Wrap(err, "failed to get container status")
		}

		err = c.handleContainerCreated(ctx, name, status, timeout)
		if err != nil {
			return name, errors.Wrap(err, "failed to handle container created condition")
		}

		// When status is StatusDead/StatusRemoving, we need to create and start the container later(not to return directly).
		if status != containerType.StatusDead && status != containerType.StatusRemoving {
			return name, nil
		}
	}
	resp, err := c.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", errors.Wrap(err, "failed to create container")
	}

	for _, w := range resp.Warnings {
		logger.Warnf("run with warnings: %s", w)
	}

	if err := c.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
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
		if errdefs.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c dockerClient) IsRunning(ctx context.Context, cname string) (bool, error) {
	container, err := c.ContainerInspect(ctx, cname)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return container.State.Running, nil
}

func (c dockerClient) GetStatus(ctx context.Context, cname string) (containerType.ContainerStatus, error) {
	container, err := c.ContainerInspect(ctx, cname)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return "", nil
		}
		return "", err
	}
	return containerType.ContainerStatus(container.State.Status), nil
}

// Load loads the docker image from the reader into the docker host.
// It's up to the caller to close the io.ReadCloser.
func (c dockerClient) Load(ctx context.Context, r io.ReadCloser, quiet bool) error {
	resp, err := c.ImageLoad(ctx, r, client.ImageLoadWithQuiet(quiet))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (c dockerClient) Exec(ctx context.Context, cname string, cmd []string) error {
	execConfig := container.ExecOptions{
		Cmd:    cmd,
		Detach: true,
	}
	resp, err := c.ContainerExecCreate(ctx, cname, execConfig)
	if err != nil {
		return err
	}
	execID := resp.ID
	return c.ContainerExecStart(ctx, execID, container.ExecStartOptions{
		Detach: true,
	})
}

func (c dockerClient) PruneImage(ctx context.Context) (dockerimage.PruneReport, error) {
	pruneReport, err := c.ImagesPrune(ctx, filters.Args{})
	if err != nil {
		return dockerimage.PruneReport{}, errors.Wrap(err, "failed to prune images")
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

func (c dockerClient) waitUntilRemoved(ctx context.Context,
	name string, timeout time.Duration) error {
	logger := logrus.WithField("container", name)
	logger.Debug("waiting to be removed")

	// Wait for the container to be removed
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-time.After(waitingInterval):
			exist, err := c.Exists(ctxTimeout, name)
			if err != nil {
				return errors.Wrap(err, "failed to check if container has been removed")
			}
			if !exist {
				logger.Debug("the container has been removed")
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
			return errors.Errorf("timeout %s: container can't be removed", timeout)
		}
	}
}

func (c dockerClient) handleContainerCreated(ctx context.Context,
	cname string, status containerType.ContainerStatus, timeout time.Duration) error {
	logger := logrus.WithFields(logrus.Fields{
		"container": cname,
		"status":    status,
	})

	switch status {
	case containerType.StatusPaused:
		logger.Info("container was paused, unpause it now...")
		if _, err := c.ResumeContainer(ctx, cname); err != nil {
			logger.WithError(err).Error("can not run buildkitd")
			return errors.Wrap(err, "failed to unpause container")
		}
	case containerType.StatusExited:
		logger.Info("container exited, try to start it...")
		if err := c.ContainerStart(ctx, cname, container.StartOptions{}); err != nil {
			logger.WithError(err).Error("can not run buildkitd")
			return errors.Wrap(err, "failed to start exited container")
		}
	case containerType.StatusDead:
		logger.Info("container is dead, try to remove it...")
		if err := c.ContainerRemove(ctx, cname, container.RemoveOptions{}); err != nil {
			logger.WithError(err).Error("can not run buildkitd")
			return errors.Wrap(err, "failed to remove container")
		}
	case containerType.StatusCreated:
		logger.Info("container is being created")
		if err := c.waitUntilRunning(ctx, cname, timeout); err != nil {
			logger.WithError(err).Error("can not run buildkitd")
			return errors.Wrap(err, "failed to start container")
		}
	case containerType.StatusRemoving:
		logger.Info("container is being removed.")
		if err := c.waitUntilRemoved(ctx, cname, timeout); err != nil {
			logger.WithError(err).Error("can not run buildkitd")
			return errors.Wrap(err, "failed to remove container")
		}
	}
	// No process for StatusRunning

	return nil
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
