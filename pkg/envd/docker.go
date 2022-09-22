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
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
	"github.com/tensorchord/envd/pkg/util/netutil"
)

const (
	localhost = "127.0.0.1"
)

var (
	waitingInternal = 1 * time.Second
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

func dockerFiltersWithCacheLabel(name string, hash string) filters.Args {
	f := filters.NewArgs()
	f.Add("reference", name)
	f.Add("label", fmt.Sprintf("%s=%s", types.ImageLabelCacheHash, hash))
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
		envdImg, err := types.NewImage(img)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create envd image from the docker image")
		}
		envdImgs = append(envdImgs, *envdImg)
	}
	return envdImgs, nil
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
		env, err := types.NewEnvironment(ctr)
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
	logger = logrus.WithField("container", env)
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
	logger = logrus.WithField("container", env)
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
	dep, err := types.NewDependencyFromImage(img)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dependency from image")
	}
	return dep, nil
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
	ports := types.NewPortBindingFromContainerJSON(ctr)
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

	// force delete the container no matter it is running or not.
	if force {
		return e.ContainerRemove(ctx, name, dockertypes.ContainerRemoveOptions{
			Force: true,
		})
	}

	running, _ := e.IsRunning(ctx, name)
	if err != nil {
		return err
	}
	if running {
		logrus.Errorf("container %s is running, cannot clean envd, please save your data and stop the running container if you need to envd up again.", name)
		return errors.Newf("\"%s\" is stil running, please run `envd destroy --name %s` stop it first", name, name)
	}
	return e.ContainerRemove(ctx, name, dockertypes.ContainerRemoveOptions{})
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

// StartEnvd creates the container for the given tag and container name.
func (e dockerEngine) StartEnvd(ctx context.Context, tag, name, buildContext string,
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

	resp, err := e.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create the container")
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
			return "", "", errors.New("port is already allocated in the host")
		}
		return "", "", errors.Wrap(err, "failed to run the container")
	}

	container, err := e.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to inspect the container")
	}

	if err := e.WaitUntilRunning(
		ctx, container.Name, timeout); err != nil {
		return "", "", errors.Wrap(err, "failed to wait until the container is running")
	}

	return container.Name, container.NetworkSettings.IPAddress, nil
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
		case <-time.After(waitingInternal):
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
	return nv.Path != "", nil
}

func (e dockerEngine) GetImage(ctx context.Context, image string) (dockertypes.ImageSummary, error) {
	images, err := e.ImageList(ctx, dockertypes.ImageListOptions{
		Filters: dockerFiltersWithName(image),
	})
	if err != nil {
		return dockertypes.ImageSummary{}, err
	}
	if len(images) == 0 {
		return dockertypes.ImageSummary{},
			errors.Errorf("image %s not found", image)
	}
	return images[0], nil
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

func labels(name string, g ir.Graph,
	sshPortInHost, jupyterPortInHost, rstudioServerPortInHost int) map[string]string {
	res := make(map[string]string)
	res[types.ContainerLabelName] = name
	res[types.ContainerLabelSSHPort] = strconv.Itoa(sshPortInHost)
	if g.JupyterConfig != nil {
		res[types.ContainerLabelJupyterAddr] =
			fmt.Sprintf("http://%s:%d", localhost, jupyterPortInHost)
	}
	if g.RStudioServerConfig != nil {
		res[types.ContainerLabelRStudioServerAddr] =
			fmt.Sprintf("http://%s:%d", localhost, rstudioServerPortInHost)
	}

	return res
}
