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

const (
	ContainerLabelName              = "ai.tensorchord.envd.name"
	ContainerLabelJupyterAddr       = "ai.tensorchord.envd.jupyter.address"
	ContainerLabelRStudioServerAddr = "ai.tensorchord.envd.rstudio.server.address"
	ContainerLabelSSHPort           = "ai.tensorchord.envd.ssh.port"

	ImageLabelContainerName = "ai.tensorchord.envd.container.name"
	ImageLabelVendor        = "ai.tensorchord.envd.vendor"
	ImageLabelGPU           = "ai.tensorchord.envd.gpu"
	ImageLabelRepo          = "ai.tensorchord.envd.repo"
	ImageLabelPorts         = "ai.tensorchord.envd.ports"
	ImageLabelAPT           = "ai.tensorchord.envd.apt.packages"
	ImageLabelPyPI          = "ai.tensorchord.envd.pypi.commands"
	ImageLabelR             = "ai.tensorchord.envd.r.packages"
	ImageLabelCUDA          = "ai.tensorchord.envd.gpu.cuda"
	ImageLabelCUDNN         = "ai.tensorchord.envd.gpu.cudnn"
	ImageLabelContext       = "ai.tensorchord.envd.build.context"
	ImageLabelCacheHash     = "ai.tensorchord.envd.build.digest"
	ImageLabelSyntaxVer     = "ai.tensorchord.envd.syntax.version"
	RuntimeGraphCode        = "ai.tensorchord.envd.graph.runtime"
	GeneralGraphCode        = "ai.tensorchord.envd.graph.general"

	ImageVendorEnvd = "envd"
)
