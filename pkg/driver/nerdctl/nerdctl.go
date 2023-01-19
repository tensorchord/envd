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
	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/driver"
)

type nerdctlClient struct {
	bin string
}

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

func (nc *nerdctlClient) StartBuildkitd(ctx context.Context,
	tag, name, mirror string, timeout time.Duration) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":       tag,
		"container": name,
		"mirror":    mirror,
		"driver":    "nerdctl",
	})
	logger.Debug("starting buildkitd")

	if err := nc.imageInspect(ctx, tag); err != nil {
		if err := nc.imagePull(ctx, tag); err != nil {
			return "", errors.Wrap(err, "pulling buildkitd image")
		}
	}

	existed, _ := nc.containerExists(ctx, name)
	if !existed {
		buildkitdCmd := "buildkitd"
		if mirror != "" {
			cfg := fmt.Sprintf(`
[registry."docker.io"]
	mirrors = ["%s"]`, mirror)
			buildkitdCmd = fmt.Sprintf("mkdir /etc/buildkit && echo '%s' > /etc/buildkit/buildkitd.toml && buildkitd", cfg)
			logger.Debugf("setting buildkit config: %s", cfg)
		}

		out, err := nc.exec(ctx, nil, "run", "-d",
			"--name", name,
			"--privileged",
			"--entrypoint", "sh",
			tag, "-c", buildkitdCmd)
		if err != nil {
			logrus.Error("can not run buildkitd", out, err)
			return "", errors.Wrap(err, "running buildkitd")
		}
	}

	err := nc.waitUntilRunning(ctx, name, timeout)

	return name, err
}

func (nc *nerdctlClient) Exec(ctx context.Context, cname string, cmd []string) error {
	return nil
}

func (nc *nerdctlClient) GetImageWithCacheHashLabel(ctx context.Context, image string, hash string) (types.ImageSummary, error) {
	return types.ImageSummary{}, nil
}
func (nc *nerdctlClient) RemoveImage(ctx context.Context, image string) error {
	return nil
}
func (nc *nerdctlClient) PruneImage(ctx context.Context) (types.ImagesPruneReport, error) {
	return types.ImagesPruneReport{}, nil
}
func (nc *nerdctlClient) Stats(ctx context.Context, cname string, statChan chan<- *driver.Stats, done <-chan bool) error {
	return nil
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
		case <-time.After(time.Second):
			_, err := nc.exec(ctx, nil, "start", name)
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

func (nc *nerdctlClient) containerExists(ctx context.Context, tag string) (bool, error) {
	_, err := nc.containerInspect(ctx, tag)
	return err == nil, err
}

func (nc *nerdctlClient) containerInspect(ctx context.Context, tag string) (*types.ContainerJSON, error) {
	out, err := nc.exec(ctx, nil, "inspect", tag)
	if err != nil {
		//TODO(kweizh): check not found
		return nil, err
	}

	cs := []types.ContainerJSON{}
	err = json.Unmarshal([]byte(out), &cs)
	if err != nil {
		logrus.Error(cs, err)
		return nil, err
	}

	return &(cs[0]), nil
}

func (nc *nerdctlClient) exec(ctx context.Context, stdin io.Reader, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, nc.bin, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	if stdin != nil {
		cmd.Stdin = stdin
	}

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
		//TODO(kweizh): check not found
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
