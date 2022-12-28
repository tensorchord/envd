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

package table

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/version"
)

const indent = "  "

func PrintVersion(clicontext *cli.Context) error {
	IdentPrintf := func(format string, a ...any) {
		fmt.Printf(indent+format, a...)
	}

	short := clicontext.Bool("short")
	detail := clicontext.Bool("detail")
	ver := version.GetVersion()
	detailVer, err := formatter.GetDetailedVersion(clicontext)
	fmt.Printf("envd: %s\n", ver)
	if short {
		return nil
	}
	IdentPrintf("BuildDate: %s\n", ver.BuildDate)
	IdentPrintf("GitCommit: %s\n", ver.GitCommit)
	IdentPrintf("GitTreeState: %s\n", ver.GitTreeState)
	if ver.GitTag != "" {
		IdentPrintf("GitTag: %s\n", ver.GitTag)
	}
	IdentPrintf("GoVersion: %s\n", ver.GoVersion)
	IdentPrintf("Compiler: %s\n", ver.Compiler)
	IdentPrintf("Platform: %s\n", ver.Platform)
	if detail {
		if err != nil {
			fmt.Printf("Error in getting details from Docker Server: %s\n", err)
		} else {
			IdentPrintf("OSType: %s\n", detailVer.OSType)
			if detailVer.OSVersion != "" {
				IdentPrintf("OSVersion: %s\n", detailVer.OSVersion)
			}
			IdentPrintf("KernelVersion: %s\n", detailVer.KernelVersion)
			IdentPrintf("DockerHostVersion: %s\n", detailVer.DockerVersion)
			IdentPrintf("ContainerRuntimes: %s\n", detailVer.ContainerRuntimes)
			IdentPrintf("DefaultRuntime: %s\n", detailVer.DefaultRuntime)
		}
	}
	return nil
}
