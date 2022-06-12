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
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types/filters"

	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/types"
)

func labels(name string, jupyterConfig *ir.JupyterConfig, sshPort int) map[string]string {
	res := make(map[string]string)
	res[types.ContainerLabelName] = name
	res[types.ContainerLabelSSHPort] = strconv.Itoa(sshPort)
	if jupyterConfig != nil {
		res[types.ContainerLabelJupyterAddr] = fmt.Sprintf("http://localhost:%d", jupyterConfig.Port)
	}
	return res
}

func dockerfilters(gpu bool) filters.Args {
	f := filters.NewArgs()
	f.Add("label", fmt.Sprintf("%s=%s", types.ImageLabelVendor, types.ImageVendorEnvd))
	if gpu {
		f.Add("label", fmt.Sprintf("%s=true", types.ImageLabelGPU))
	}
	return f
}

func dockerfiltersWithName(name string) filters.Args {
	f := filters.NewArgs()
	f.Add("reference", name)
	return f
}
