package types

import (
	"github.com/docker/docker/api/types"
)

type EnvdEnvironment struct {
	types.Container

	// The name of the environment.
	Name        string `json:"name"`
	GPU         bool   `json:"gpu"`
	JupyterAddr string `json:"jupyter_addr"`
}

const (
	ContainerLabelGPU         = "ai.tensorchord.envd.gpu"
	ContainerLabelVendor      = "ai.tensorchord.envd.vendor"
	ContainerLabelName        = "ai.tensorchord.envd.name"
	ContainerLabelJupyterAddr = "ai.tensorchord.envd.jupyter.address"

	ImageLabelAPT  = "ai.tensorchord.envd.apt.packages"
	ImageLabelPyPI = "ai.tensorchord.envd.pypi.packages"
)

func FromContainer(ctr types.Container) EnvdEnvironment {
	env := EnvdEnvironment{
		Container: ctr,
	}
	if gpuEnabled, ok := ctr.Labels[ContainerLabelGPU]; ok {
		env.GPU = gpuEnabled == "true"
	}
	if name, ok := ctr.Labels[ContainerLabelName]; ok {
		env.Name = name
	}
	if jupyterAddr, ok := ctr.Labels[ContainerLabelJupyterAddr]; ok {
		env.JupyterAddr = jupyterAddr
	}

	return env
}
