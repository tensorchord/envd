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

package nerdctl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/api/types/container"
	dockerimage "github.com/docker/docker/api/types/image"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/driver"
	containerType "github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/buildkitutil"
)

type nerdctlClient struct {
	bin string
}

var (
	waitingInterval = 1 * time.Second
)

func NewClient(ctx context.Context) (driver.Client, error) {
	bin, err := exec.LookPath("nerdctl")
	if err != nil {
		// for mac lima users, check nerdctl.lima
		bin, err = exec.LookPath("nerdctl.lima")
		if err != nil {
			return nil, errors.New("can not found nerdctl(or nerdctl.lima for mac) in PATH")
		}
	}

	return &nerdctlClient{
		bin: bin,
	}, nil
}

func (nc *nerdctlClient) Load(ctx context.Context, r io.ReadCloser, quiet bool) error {
	cmd := exec.CommandContext(ctx, nc.bin, "load")
	cmd.Stdin = r
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	logrus.Debug(out.String())
	return nil
}

func (nc *nerdctlClient) StartBuildkitd(ctx context.Context, tag, name string, bc *buildkitutil.BuildkitConfig, timeout time.Duration) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":             tag,
		"container":       name,
		"buildkit-config": bc,
		"driver":          "nerdctl",
	})
	logger.Debug("starting buildkitd")

	if err := nc.imageInspect(ctx, tag); err != nil {
		if err := nc.imagePull(ctx, tag); err != nil {
			return "", errors.Wrap(err, "pulling buildkitd image")
		}
	}

	existed, _ := nc.containerExists(ctx, name)
	if existed {
		status, err := nc.GetStatus(ctx, name)
		if err != nil {
			return "", errors.Wrap(err, "failed to get container status")
		}

		err = nc.handleContainerCreated(ctx, name, status, timeout)
		if err != nil {
			return "", errors.Wrap(err, "failed to handle container created condition")
		}

		// When status is StatusDead/StatusRemoving, we nened to create and start the container later(not to return directly).
		if status != containerType.StatusDead && status != containerType.StatusRemoving {
			return name, nil
		}
	}

	buildkitdCmd := "buildkitd"
	// TODO: support mirror CA keypair
	if len(bc.Registries) > 0 {
		cfg, err := bc.String()
		if err != nil {
			return "", errors.Wrap(err, "failed to generate buildkit config")
		}
		buildkitdCmd = fmt.Sprintf("mkdir /etc/buildkit && echo '%s' > /etc/buildkit/buildkitd.toml && buildkitd", cfg)
		logger.Debugf("setting buildkit config: %s", cfg)
	}

	out, err := nc.exec(ctx, "run", "-d",
		"--name", name,
		"--privileged",
		"--entrypoint", "sh",
		tag, "-c", buildkitdCmd)
	if err != nil {
		logger.WithError(err).Error("can not run buildkitd", out)
		return "", errors.Wrap(err, "running buildkitd")
	}

	err = nc.waitUntilRunning(ctx, name, timeout)

	return name, err
}

func (nc *nerdctlClient) Exec(ctx context.Context, cname string, cmd []string) error {
	return nil
}

func (nc *nerdctlClient) GetImageWithCacheHashLabel(ctx context.Context, image string, hash string) (dockerimage.Summary, error) {
	return dockerimage.Summary{}, nil
}
func (nc *nerdctlClient) RemoveImage(ctx context.Context, image string) error {
	return nil
}
func (nc *nerdctlClient) PushImage(ctx context.Context, image, platform string) error {
	return nil
}
func (nc *nerdctlClient) PruneImage(ctx context.Context) (dockerimage.PruneReport, error) {
	return dockerimage.PruneReport{}, nil
}
func (nc *nerdctlClient) Stats(ctx context.Context, cname string, statChan chan<- *driver.Stats, done <-chan bool) error {
	return nil
}

func (nc *nerdctlClient) GetStatus(ctx context.Context, cname string) (containerType.ContainerStatus, error) {
	container, err := nc.containerInspect(ctx, cname)
	if err != nil {
		return "", err
	}
	return containerType.ContainerStatus(container.State.Status), nil
}

// TODO(kweizh): use container engine to wrap docker and nerdctl
func (nc *nerdctlClient) waitUntilRunning(ctx context.Context,
	name string, timeout time.Duration) error {
	logger := logrus.WithField("container", name)
	logger.Debug("waiting to start")

	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-time.After(waitingInterval):
			_, err := nc.exec(ctx, "start", name)
			if err != nil {
				continue
			}

			c, err := nc.containerInspect(ctx, name)
			if err != nil {
				// Has not yet started. Keep waiting.
				continue
			}
			if c.State.Running {
				logger.Debug("the container is running")
				return nil
			}

		case <-ctxTimeout.Done():
			container, err := nc.containerInspect(ctx, name)
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

func (nc *nerdctlClient) waitUntilRemoved(ctx context.Context,
	name string, timeout time.Duration) error {
	logger := logrus.WithField("container", name)
	logger.Debug("waiting to be removed")

	// Wait for the container to be removed
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-time.After(waitingInterval):
			exist, err := nc.containerExists(ctxTimeout, name)
			if err != nil {
				return errors.Wrap(err, "failed to check if container has been removed")
			}
			if !exist {
				logger.Debug("the container has been removed")
				return nil
			}
		case <-ctxTimeout.Done():
			container, err := nc.containerInspect(ctx, name)
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

func (nc *nerdctlClient) handleContainerCreated(ctx context.Context,
	cname string, status containerType.ContainerStatus, timeout time.Duration) error {
	logger := logrus.WithFields(logrus.Fields{
		"container": cname,
		"status":    status,
	})

	switch status {
	case containerType.StatusPaused:
		logger.Info("container was paused, unpause it now...")
		if out, err := nc.exec(ctx, "unpause", cname); err != nil {
			logger.WithError(err).Error("can not run buildkitd", out)
			return errors.Wrap(err, "failed to unpause container")
		}
	case containerType.StatusExited:
		logger.Info("container exited, try to start it...")
		if out, err := nc.exec(ctx, "start", cname); err != nil {
			logger.WithError(err).Error("can not run buildkitd", out)
			return errors.Wrap(err, "failed to start exited container")
		}
	case containerType.StatusDead:
		logger.Info("container is dead, try to remove it...")
		if out, err := nc.exec(ctx, "remove", cname); err != nil {
			logger.WithError(err).Error("can not run buildkitd", out)
			return errors.Wrap(err, "failed to remove container")
		}
	case containerType.StatusCreated:
		logger.Info("container is being created.")
		if err := nc.waitUntilRunning(ctx, cname, timeout); err != nil {
			logger.WithError(err).Error("can not run buildkitd")
			return errors.Wrap(err, "failed to start container")
		}
	case containerType.StatusRemoving:
		// The remaining condition is StatusRemoving, we just need to wait.
		logger.Info("container is being removed.")
		if err := nc.waitUntilRemoved(ctx, cname, timeout); err != nil {
			logger.WithError(err).Error("can not run buildkitd")
			return errors.Wrap(err, "failed to remove container")
		}
	}
	// No process for StatusRunning

	return nil
}

func (nc *nerdctlClient) containerExists(ctx context.Context, tag string) (bool, error) {
	_, err := nc.containerInspect(ctx, tag)
	return err == nil, err
}

func (nc *nerdctlClient) containerInspect(ctx context.Context, tag string) (*container.InspectResponse, error) {
	out, err := nc.exec(ctx, "inspect", tag)
	if err != nil {
		// TODO(kweizh): check not found
		return nil, err
	}

	cs := []container.InspectResponse{}
	err = json.Unmarshal([]byte(out), &cs)
	if err != nil {
		logrus.WithError(err).Error(cs)
		return nil, err
	}

	return &(cs[0]), nil
}

func (nc *nerdctlClient) exec(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, nc.bin, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", errors.Wrap(err, "nerdctlClient error")
	}

	return out.String(), nil
}

// TODO(kweizh): return inspect result
func (nc *nerdctlClient) imageInspect(ctx context.Context, tag string) error {
	cmd := exec.CommandContext(ctx, nc.bin, "image", "inspect", tag)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// TODO(kweizh): check not found
		return err
	}

	return nil
}

// TODO(kweizh): return pull output
func (nc *nerdctlClient) imagePull(ctx context.Context, tag string) error {
	cmd := exec.CommandContext(ctx, nc.bin, "pull", tag)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
