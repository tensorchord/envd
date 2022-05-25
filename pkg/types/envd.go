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
