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
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"

	envdconfig "github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/lang/version"
	"github.com/tensorchord/envd/pkg/ssh"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
	"github.com/tensorchord/envd/pkg/util/netutil"
)

type dockerEngine struct {
	*client.Client
}

func dockerFilters(gpu bool) filters.Args {
	f := filters.NewArgs()
	f.Add("label", fmt.Sprintf("%s=%s", types.ImageLabelVendor, types.ImageVendorEnvd))
	if gpu {
		f.Add("label", fmt.Sprintf("%s=true", types.ImageLabelGPU))
	}
	return f
}

func dockerFiltersWithName(name string) filters.Args {
	f := filters.NewArgs()
	f.Add("reference", name)
	return f
}

func (e dockerEngine) ListImage(ctx context.Context) ([]types.EnvdImage, error) {
	images, err := e.ImageList(ctx, dockertypes.ImageListOptions{
		Filters: dockerFilters(false),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the images")
	}

	envdImgs := make([]types.EnvdImage, 0)
	for _, img := range images {
		envdImg, err := types.NewImageFromSummary(img)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create envd image `%s` from the docker image", img.ID)
		}
		envdImgs = append(envdImgs, *envdImg)
	}
	return envdImgs, nil
}

func (e dockerEngine) GetEnvironment(ctx context.Context, env string) (*types.EnvdEnvironment, error) {
	ctrs, err := e.ContainerList(ctx, dockertypes.ContainerListOptions{
		Filters: dockerFiltersWithName(env),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get container: %s", env)
	}

	if len(ctrs) <= 0 {
		return nil, errors.Newf("can not find the container: %s", env)
	}

	environment, err := types.NewEnvironmentFromContainer(ctrs[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to create env from the container")
	}
	return environment, nil
}

func (e dockerEngine) ListEnvironment(
	ctx context.Context) ([]types.EnvdEnvironment, error) {
	ctrs, err := e.ContainerList(ctx, dockertypes.ContainerListOptions{
		Filters: dockerFilters(false),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list containers")
	}

	envs := make([]types.EnvdEnvironment, 0)
	for _, ctr := range ctrs {
		env, err := types.NewEnvironmentFromContainer(ctr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create env from the container")
		}
		envs = append(envs, *env)
	}
	return envs, nil
}

func (e dockerEngine) PauseEnvironment(ctx context.Context, env string) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"env": env,
	})
	logger.Debug("pausing environment")
	err := e.ContainerPause(ctx, env)
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
	return env, nil
}

func (e dockerEngine) ResumeEnvironment(ctx context.Context, env string) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"env": env,
	})
	logger.Debug("resuming environment")

	err := e.ContainerUnpause(ctx, env)
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
	return env, nil
}

// ListImageDependency gets the dependencies of the given environment.
func (e dockerEngine) ListImageDependency(ctx context.Context, image string) (*types.Dependency, error) {
	logger := logrus.WithFields(logrus.Fields{
		"image": image,
	})
	logger.Debug("getting dependencies")
	images, err := e.ImageList(ctx, dockertypes.ImageListOptions{
		Filters: dockerFiltersWithName(image),
	})
	if err != nil {
		return nil, err
	}
	if len(images) == 0 {
		return nil, errors.Errorf("image %s not found", image)
	}

	img := images[0]
	dep, err := types.NewDependencyFromImageSummary(img)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dependency from image")
	}
	return dep, nil
}

func (e dockerEngine) getVerFromImageLabel(ctx context.Context, env string) (version.Getter, error) {
	ctr, err := e.GetImage(ctx, env)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to inspect image: %s", env)
	}
	ver, ok := ctr.Labels[types.ImageLabelSyntaxVer]
	if !ok {
		return version.NewByVersion("v0"), nil
	}
	return version.NewByVersion(ver), nil
}

func (e dockerEngine) listEnvGeneralGraph(ctx context.Context, env string, g ir.Graph) (ir.Graph, error) {
	ctr, err := e.GetImage(ctx, env)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to inspect image: %s", env)
	}
	code, ok := ctr.Labels[types.GeneralGraphCode]
	if !ok {
		return nil, errors.Newf("failed to get runtime graph label from image: %s", env)
	}
	logrus.WithField("env", env).Debugf("general graph: %s", code)
	newg, err := g.GeneralGraphFromLabel([]byte(code))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create runtime graph from the image: %s", env)
	}
	return newg, err
}

func (e dockerEngine) ListEnvRuntimeGraph(ctx context.Context, env string) (*ir.RuntimeGraph, error) {
	ctr, err := e.ContainerInspect(ctx, env)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to inspect container: %s", env)
	}
	code, ok := ctr.Config.Labels[types.RuntimeGraphCode]
	if !ok {
		return nil, errors.Newf("failed to get runtime graph label from container: %s", env)
	}
	logrus.WithField("env", env).Debugf("runtime graph: %s", code)
	rg := ir.RuntimeGraph{}
	err = rg.Load([]byte(code))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create runtime graph from the container: %s", env)
	}
	return &rg, err
}

// ListEnvDependency gets the dependencies of the given environment.
func (e dockerEngine) ListEnvDependency(
	ctx context.Context, env string) (*types.Dependency, error) {
	logger := logrus.WithFields(logrus.Fields{
		"env": env,
	})
	logger.Debug("getting dependencies")
	ctr, err := e.ContainerInspect(ctx, env)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get container")
	}
	dep, err := types.NewDependencyFromContainerJSON(ctr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dependency from the container")
	}
	return dep, nil
}

func (e dockerEngine) ListEnvPortBinding(ctx context.Context, env string) ([]types.PortBinding, error) {
	logrus.WithField("env", env).Debug("getting env port bindings")
	ctr, err := e.ContainerInspect(ctx, env)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get container")
	}

	ports, err := types.NewPortBindingFromContainerJSON(ctr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get port binding from container json")
	}

	return ports, nil
}

func (e dockerEngine) GetInfo(ctx context.Context) (*types.EnvdInfo, error) {
	info, err := e.Info(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get docker client info")
	}
	return &types.EnvdInfo{
		Info: info,
	}, nil
}

func (e dockerEngine) CleanEnvdIfExists(ctx context.Context, name string, force bool) error {
	created, err := e.Exists(ctx, name)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	return e.ContainerRemove(ctx, name, dockertypes.ContainerRemoveOptions{Force: force})
}

func (e dockerEngine) Exists(ctx context.Context, cname string) (bool, error) {
	_, err := e.ContainerInspect(ctx, cname)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (e dockerEngine) IsRunning(ctx context.Context, cname string) (bool, error) {
	container, err := e.ContainerInspect(ctx, cname)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return container.State.Running, nil
}

func (e dockerEngine) GenerateSSHConfig(name, iface, privateKeyPath string,
	startResult *StartResult) (sshconfig.EntryOptions, error) {
	eo := sshconfig.EntryOptions{
		Name:               name,
		IFace:              iface,
		Port:               startResult.SSHPort,
		PrivateKeyPath:     privateKeyPath,
		EnableHostKeyCheck: false,
		EnableAgentForward: true,
	}
	return eo, nil
}

func (e dockerEngine) Attach(name, iface, privateKeyPath string,
	startResult *StartResult, g ir.Graph) error {
	opt := ssh.DefaultOptions()
	opt.Server = iface
	opt.PrivateKeyPath = privateKeyPath
	opt.Port = startResult.SSHPort
	sshClient, err := ssh.NewClient(opt)
	if err != nil {
		return errors.Wrap(err, "failed to create the ssh client")
	}
	opt.Server = iface

	if err := sshClient.Attach(); err != nil {
		return errors.Wrap(err, "failed to attach to the container")
	}
	return nil
}

// StartEnvd creates the container for the given tag and container name.
func (e dockerEngine) StartEnvd(ctx context.Context, so StartOptions) (*StartResult, error) {
	logger := logrus.WithFields(logrus.Fields{
		"tag":           so.Image,
		"environment":   so.EnvironmentName,
		"gpu":           so.NumGPU,
		"build-context": so.BuildContext,
	})
	if so.NumGPU != 0 {
		nvruntimeExists, err := e.GPUEnabled(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to check if nvidia-runtime is installed")
		}
		if !nvruntimeExists {
			return nil, errors.New("GPU is required but nvidia container runtime is not installed, please refer to https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#docker")
		}
	}

	sshPortInHost, err := netutil.GetFreePort()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a free port")
	}

	err = e.CleanEnvdIfExists(ctx, so.EnvironmentName, so.Forced)
	if err != nil {
		return nil, errors.Wrap(err, "failed to clean the envd environment")
	}
	config := &container.Config{
		Image:        so.Image,
		ExposedPorts: nat.PortSet{},
	}
	base := fileutil.EnvdHomeDir(filepath.Base(so.BuildContext))
	config.WorkingDir = base

	var g ir.Graph
	if so.Image == "" {
		g = so.DockerSource.Graph
	} else {
		// Use specified image as the running image
		getter, err := e.getVerFromImageLabel(ctx, so.Image)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get the version from the image label")
		}
		defaultGraph := getter.GetDefaultGraph()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get the graph from the image")
		}
		g, err = e.listEnvGeneralGraph(ctx, so.Image, defaultGraph)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get the graph from the image")
		}
		so.DockerSource.Graph = g
	}

	if so.DockerSource == nil || so.DockerSource.Graph == nil {
		return nil, errors.New("failed to get the docker-specific options")
	}

	mountOption := make([]mount.Mount, 0,
		len(so.DockerSource.MountOptions)+len(g.GetMount())+1)
	for _, option := range so.DockerSource.MountOptions {
		mStr := strings.Split(option, ":")
		if len(mStr) != 2 {
			return nil, errors.Newf("Invalid mount options %s", option)
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
	for _, m := range g.GetMount() {
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
		Source: so.BuildContext,
		Target: base,
	})

	logger.WithFields(logrus.Fields{
		"mount-path":  so.BuildContext,
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

	// shared memory size
	if so.ShmSize > 0 {
		hostConfig.ShmSize = int64(so.ShmSize) * 1024 * 1024
	}

	// Configure ssh port.
	natPort := nat.Port(fmt.Sprintf("%d/tcp", envdconfig.SSHPortInContainer))
	hostConfig.PortBindings[natPort] = []nat.PortBinding{
		{
			HostIP:   so.SshdHost,
			HostPort: strconv.Itoa(sshPortInHost),
		},
	}

	var jupyterPortInHost int
	// TODO(gaocegege): Avoid specific logic to set the port.
	// Add a func to builder to generate all the ports from the build process.
	if g.GetJupyterConfig() != nil {
		jc := g.GetJupyterConfig()
		if jc.Port != 0 {
			jupyterPortInHost = int(jc.Port)
		} else {
			var err error
			jupyterPortInHost, err = netutil.GetFreePort()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get a free port")
			}
		}
		natPort := nat.Port(fmt.Sprintf("%d/tcp", envdconfig.JupyterPortInContainer))
		hostConfig.PortBindings[natPort] = []nat.PortBinding{
			{
				HostIP:   Localhost,
				HostPort: strconv.Itoa(jupyterPortInHost),
			},
		}
		config.ExposedPorts[natPort] = struct{}{}
	}
	var rStudioPortInHost int
	if g.GetRStudioServerConfig() != nil {
		var err error
		rStudioPortInHost, err = netutil.GetFreePort()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get a free port")
		}
		natPort := nat.Port(fmt.Sprintf("%d/tcp", envdconfig.RStudioServerPortInContainer))
		hostConfig.PortBindings[natPort] = []nat.PortBinding{
			{
				HostIP:   Localhost,
				HostPort: strconv.Itoa(rStudioPortInHost),
			},
		}
		config.ExposedPorts[natPort] = struct{}{}
	}

	if len(g.GetExposedPorts()) > 0 {

		for _, item := range g.GetExposedPorts() {
			var err error
			if item.HostPort == 0 {
				item.HostPort, err = netutil.GetFreePort()
				if err != nil {
					return nil, errors.Wrap(err, "failed to get a free port")
				}
			}
			natPort := nat.Port(fmt.Sprintf("%d/tcp", item.EnvdPort))
			hostConfig.PortBindings[natPort] = []nat.PortBinding{
				{
					HostIP:   item.ListeningAddr,
					HostPort: strconv.Itoa(item.HostPort),
				},
			}
			config.ExposedPorts[natPort] = struct{}{}
		}
	}

	if so.NumGPU != 0 {
		logger.Debug("GPU is enabled.")
		hostConfig.DeviceRequests = deviceRequests(so.NumGPU)
	}

	config.Labels = e.labels(g, so.EnvironmentName,
		sshPortInHost, jupyterPortInHost, rStudioPortInHost)

	logger = logger.WithFields(logrus.Fields{
		"entrypoint":  config.Entrypoint,
		"working-dir": config.WorkingDir,
	})
	logger.Debugf("starting %s container", so.EnvironmentName)

	resp, err := e.ContainerCreate(ctx, config, hostConfig, nil, nil, so.EnvironmentName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the container")
	}

	for _, w := range resp.Warnings {
		logger.Warnf("run with warnings: %s", w)
	}

	if err := e.ContainerStart(
		ctx, resp.ID, dockertypes.ContainerStartOptions{}); err != nil {
		errCause := errors.UnwrapAll(err)
		// Hack to check if the port is already allocated.
		if strings.Contains(errCause.Error(), "port is already allocated") {
			logrus.Debugf("failed to allocate the port: %s", err)
			return nil, errors.New("port is already allocated in the host")
		}
		return nil, errors.Wrap(err, "failed to run the container")
	}

	container, err := e.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to inspect the container")
	}

	if err := e.WaitUntilRunning(
		ctx, container.Name, so.Timeout); err != nil {
		return nil, errors.Wrap(err, "failed to wait until the container is running")
	}

	result := &StartResult{
		SSHPort: sshPortInHost,
		Address: container.NetworkSettings.IPAddress,
		Name:    container.Name,
	}
	return result, nil
}

func (e dockerEngine) Destroy(ctx context.Context, name string) (string, error) {
	logger := logrus.WithField("container", name)

	ctr, err := e.ContainerInspect(ctx, name)
	if err != nil {
		errCause := errors.UnwrapAll(err).Error()
		if strings.Contains(errCause, "No such container") {
			// If the container is not found, it is already destroyed or the name is wrong.
			logger.Infof("cannot find container %s, maybe it's already destroyed or the name is wrong", name)
			return "", nil
		}
		return "", errors.Wrap(err, "failed to inspect the container")
	}

	// Refer to https://docs.docker.com/engine/reference/commandline/container_kill/
	if err := e.ContainerKill(ctx, name, "KILL"); err != nil {
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

	if err := e.ContainerRemove(ctx, name, dockertypes.ContainerRemoveOptions{}); err != nil {
		return "", errors.Wrap(err, "failed to remove the container")
	}

	if _, err := e.ImageRemove(ctx, ctr.Image, dockertypes.ImageRemoveOptions{}); err != nil {
		return "", errors.Errorf("remove image %s failed: %w", ctr.Image, err)
	}
	logrus.Infof("image(%s) is destroyed", ctr.Image)

	return name, nil
}

func (e dockerEngine) WaitUntilRunning(ctx context.Context,
	name string, timeout time.Duration) error {
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
				return errors.Wrap(err, "failed to check if container is running")
			}
			if isRunning {
				logger.Debug("the container is running")
				return nil
			}

		case <-ctxTimeout.Done():
			container, err := e.ContainerInspect(ctx, name)
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

func (e dockerEngine) GPUEnabled(ctx context.Context) (bool, error) {
	info, err := e.GetInfo(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get docker info")
	}
	logrus.WithField("info", info).Debug("docker info")
	nv := info.Runtimes["nvidia"]
	if nv.Path != "" {
		return true, nil
	} else if strings.HasSuffix(info.KernelVersion, "WSL2") {
		logrus.Warn("We couldn't detect if your runtime support GPU on WSL2, we will continue to run your environment.")
		return true, nil
	} else {
		return false, nil
	}
}

func (e dockerEngine) GetImage(ctx context.Context, image string) (types.EnvdImage, error) {
	images, err := e.ImageList(ctx, dockertypes.ImageListOptions{
		Filters: dockerFiltersWithName(image),
	})
	if err != nil {
		return types.EnvdImage{}, err
	}
	if len(images) == 0 {
		return types.EnvdImage{},
			errors.Errorf("image %s not found", image)
	}
	img, err := types.NewImageFromSummary(images[0])
	if err != nil {
		return types.EnvdImage{}, err
	}
	return *img, nil
}

func (e dockerEngine) PruneImage(ctx context.Context) (dockertypes.ImagesPruneReport, error) {
	pruneReport, err := e.ImagesPrune(ctx, filters.Args{})
	if err != nil {
		return dockertypes.ImagesPruneReport{}, errors.Wrap(err, "failed to prune images")
	}
	return pruneReport, nil
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

func (e dockerEngine) labels(g ir.Graph, name string,
	sshPortInHost, jupyterPortInHost, rstudioServerPortInHost int) map[string]string {
	res := make(map[string]string)
	res[types.ContainerLabelName] = name
	res[types.ContainerLabelSSHPort] = strconv.Itoa(sshPortInHost)
	if g.GetJupyterConfig() != nil {
		res[types.ContainerLabelJupyterAddr] =
			fmt.Sprintf("http://%s:%d", Localhost, jupyterPortInHost)
	}
	if g.GetRStudioServerConfig() != nil {
		res[types.ContainerLabelRStudioServerAddr] =
			fmt.Sprintf("http://%s:%d", Localhost, rstudioServerPortInHost)
	}

	return res
}
