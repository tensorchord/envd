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

package json

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/version"
)

type versionJSON struct {
	Envd              string `json:"envd"`
	BuildDate         string `json:"build_date,omitempty"`
	GitCommit         string `json:"git_commit,omitempty"`
	GitTreeState      string `json:"git_tree_state,omitempty"`
	GitTag            string `json:"git_tag,omitempty"`
	GoVersion         string `json:"go_version,omitempty"`
	Compiler          string `json:"compiler,omitempty"`
	Platform          string `json:"platform,omitempty"`
	OSType            string `json:"os_type,omitempty"`
	OSVersion         string `json:"os_version,omitempty"`
	KernelVersion     string `json:"kernel_version,omitempty"`
	DockerHostVersion string `json:"docker_host_version,omitempty"`
	ContainerRuntimes string `json:"container_runtimes,omitempty"`
	DefaultRuntime    string `json:"default_runtime,omitempty"`
}

func PrintVersion(clicontext *cli.Context) error {
	short := clicontext.Bool("short")
	detail := clicontext.Bool("detail")
	ver := version.GetVersion()
	detailVer, err := formatter.GetDetailedVersion(clicontext)
	output := versionJSON{
		Envd: version.GetVersion().String(),
	}
	if short {
		return printJSON(output)
	}
	output.BuildDate = ver.BuildDate
	output.GitCommit = ver.GitCommit
	output.GitTreeState = ver.GitTreeState
	if ver.GitTag != "" {
		output.GitTag = ver.GitTag
	}
	output.GoVersion = ver.GoVersion
	output.Compiler = ver.Compiler
	output.Platform = ver.Platform
	if detail {
		if err != nil {
			fmt.Printf("Error in getting details from Docker Server: %s\n", err)
		} else {
			output.OSType = detailVer.OSType
			if detailVer.OSVersion != "" {
				output.OSVersion = detailVer.OSVersion
			}
			output.KernelVersion = detailVer.KernelVersion
			output.DockerHostVersion = detailVer.DockerVersion
			output.ContainerRuntimes = detailVer.ContainerRuntimes
			output.DefaultRuntime = detailVer.DefaultRuntime
		}
	}
	return printJSON(output)
}
