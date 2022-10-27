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

package envd

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
	servertypes "github.com/tensorchord/envd-server/api/types"
	"github.com/tensorchord/envd-server/client"
	"github.com/tensorchord/envd-server/errdefs"
	"github.com/tensorchord/envd-server/sshname"

	"github.com/tensorchord/envd/pkg/ssh"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/netutil"
)

type envdServerEngine struct {
	*client.Client
	IdentityToken string
}

func (e *envdServerEngine) ListImage(ctx context.Context) ([]types.EnvdImage, error) {
	return nil, errors.New("not implemented")
}

func (e envdServerEngine) Destroy(ctx context.Context, name string) (string, error) {
	err := e.EnvironmentRemove(ctx, e.IdentityToken, name)
	return name, err
}

func (e *envdServerEngine) ListImageDependency(ctx context.Context, image string) (*types.Dependency, error) {
	return nil, errors.New("not implemented")
}

func (e *envdServerEngine) GetImage(ctx context.Context, image string) (dockertypes.ImageSummary, error) {
	return dockertypes.ImageSummary{}, errors.New("not implemented")
}

func (e envdServerEngine) PruneImage(ctx context.Context) (dockertypes.ImagesPruneReport, error) {
	return dockertypes.ImagesPruneReport{}, errors.New("not implemented")
}

func (e *envdServerEngine) GetInfo(ctx context.Context) (*types.EnvdInfo, error) {
	return nil, errors.New("not implemented")
}

func (e *envdServerEngine) GPUEnabled(ctx context.Context) (bool, error) {
	return false, errors.New("not implemented")
}

func (e *envdServerEngine) PauseEnvironment(ctx context.Context, env string) (string, error) {
	return "", errors.New("pausing/resuming environments is not supported for the runner envd-server")
}

func (e *envdServerEngine) ResumeEnvironment(ctx context.Context, env string) (string, error) {
	return "", errors.New("pausing/resuming environments is not supported for the runner envd-server")
}

func (e *envdServerEngine) ListEnvironment(ctx context.Context) ([]types.EnvdEnvironment, error) {
	env, err := e.EnvironmentList(ctx, e.IdentityToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get environment")
	}
	res := []types.EnvdEnvironment{}
	for _, e := range env.Items {
		env, err := types.NewEnvironmentFromServer(e)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create env from the container")
		}
		res = append(res, *env)
	}

	return res, nil
}

func (e envdServerEngine) GenerateSSHConfig(name, iface, privateKeyPath string, startResult *StartResult) (sshconfig.EntryOptions, error) {
	username, err := sshname.Username(e.IdentityToken, startResult.Name)
	if err != nil {
		return sshconfig.EntryOptions{}, errors.Wrap(err, "failed to get the username")
	}

	eo := sshconfig.EntryOptions{
		Name:               name,
		IFace:              iface,
		Port:               startResult.SSHPort,
		PrivateKeyPath:     privateKeyPath,
		EnableHostKeyCheck: false,
		EnableAgentForward: false,
		User:               username,
	}
	return eo, nil
}

func (e envdServerEngine) Attach(name, iface, privateKeyPath string, startResult *StartResult) error {
	username, err := sshname.Username(e.IdentityToken, startResult.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get the username")
	}

	outputChannel := make(chan error)
	opt := ssh.DefaultOptions()
	opt.PrivateKeyPath = privateKeyPath
	opt.Port = startResult.SSHPort
	opt.AgentForwarding = false
	opt.User = username
	opt.Server = iface

	sshClient, err := ssh.NewClient(opt)
	if err != nil {
		outputChannel <- errors.Wrap(err, "failed to create the ssh client")
	}

	ports := startResult.Ports

	for _, p := range ports {
		if p.Port == 2222 {
			continue
		}

		// TODO(gaocegege): Use one remote port.
		localPort, err := netutil.GetFreePort()
		if err != nil {
			return errors.Wrap(err, "failed to get a free port")
		}
		localAddress := fmt.Sprintf("%s:%d", "localhost", localPort)
		remoteAddress := fmt.Sprintf("%s:%d", "localhost", p.Port)
		logrus.Infof("service \"%s\" is listening at %s\n", p.Name, localAddress)
		go func() {
			if err := sshClient.LocalForward(localAddress, remoteAddress); err != nil {
				outputChannel <- errors.Wrap(err, "failed to forward to local port")
			}
		}()
	}

	go func() {
		if err := sshClient.Attach(); err != nil {
			outputChannel <- errors.Wrap(err, "failed to attach to the container")
		}
		outputChannel <- nil
	}()

	if err := <-outputChannel; err != nil {
		return err
	}

	return nil
}

func (e *envdServerEngine) ListEnvDependency(
	ctx context.Context, name string) (*types.Dependency, error) {
	logger := logrus.WithFields(logrus.Fields{
		"env": name,
	})
	logger.Debug("getting dependencies")
	env, err := e.EnvironmentGet(ctx, e.IdentityToken, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get environment")
	}
	dep, err := types.NewDependencyFromLabels(env.Labels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dependency from the container")
	}
	return dep, nil
}

func (e *envdServerEngine) ListEnvPortBinding(
	ctx context.Context, name string) ([]types.PortBinding, error) {
	_, err := e.EnvironmentGet(ctx, e.IdentityToken, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get environment")
	}
	// TODO(gaocegege): Remove hard coded.
	res := []types.PortBinding{
		{
			Port:     "2222",
			Protocol: "TCP",
			HostIP:   "localhost",
			HostPort: "2222",
		},
	}
	return res, nil
}

func (e *envdServerEngine) CleanEnvdIfExists(ctx context.Context, name string, force bool) error {
	created, err := e.Exists(ctx, name)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}

	return e.EnvironmentRemove(ctx, e.IdentityToken, name)
}

// StartEnvd creates the container for the given tag and container name.
func (e *envdServerEngine) StartEnvd(ctx context.Context, so StartOptions) (*StartResult, error) {
	if so.EnvdServerSource == nil {
		return nil, errors.New("failed to get the envd server specific options")
	}

	req := servertypes.EnvironmentCreateRequest{
		Environment: servertypes.Environment{
			ObjectMeta: servertypes.ObjectMeta{
				Name: so.EnvironmentName,
			},
			Spec: servertypes.EnvironmentSpec{
				Image: so.Image,
			},
		},
	}

	resp, err := e.EnvironmentCreate(ctx, e.IdentityToken, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the environment")
	}

	if err := e.WaitUntilRunning(
		ctx, resp.Created.Name, so.Timeout); err != nil {
		return nil, errors.Wrap(err, "failed to wait until the container is running")
	}

	result := &StartResult{
		SSHPort: 2222,
		Address: "",
		Name:    resp.Created.Name,
		Ports:   resp.Created.Spec.Ports,
	}
	return result, nil
}

func (e *envdServerEngine) IsRunning(ctx context.Context, name string) (bool, error) {
	resp, err := e.EnvironmentGet(ctx, e.IdentityToken, name)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to list the environment")
	}
	// "Running" is hard-coded here.
	return resp.Status.Phase == "Running", nil
}

func (e *envdServerEngine) Exists(ctx context.Context, name string) (bool, error) {
	_, err := e.EnvironmentGet(ctx, e.IdentityToken, name)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to list the environment")
	}
	return true, nil
}

func (e *envdServerEngine) WaitUntilRunning(ctx context.Context, name string, timeout time.Duration) error {
	logger := logrus.WithField("container", name)
	logger.Debug("waiting to start")

	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-time.After(waitingInterval):
			isRunning, err := e.IsRunning(ctxTimeout, name)
			if err != nil {
				// Has not yet started. Keep waiting.
				return errors.Wrap(err, "failed to check if environment is running")
			}
			if isRunning {
				logger.Debug("the environment is running")
				return nil
			}

		case <-ctxTimeout.Done():
			return errors.Errorf("timeout %s: environment did not start", timeout)
		}
	}
}
