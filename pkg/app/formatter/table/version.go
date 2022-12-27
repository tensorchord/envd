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

func PrintVersion(clicontext *cli.Context) error {
	short := clicontext.Bool("short")
	detail := clicontext.Bool("detail")
	ver := version.GetVersion()
	detailVer, err := formatter.GetDetailedVersion(clicontext)
	fmt.Printf("envd: %s\n", ver)
	if short {
		return nil
	}
	fmt.Printf("  BuildDate: %s\n", ver.BuildDate)
	fmt.Printf("  GitCommit: %s\n", ver.GitCommit)
	fmt.Printf("  GitTreeState: %s\n", ver.GitTreeState)
	if ver.GitTag != "" {
		fmt.Printf("  GitTag: %s\n", ver.GitTag)
	}
	fmt.Printf("  GoVersion: %s\n", ver.GoVersion)
	fmt.Printf("  Compiler: %s\n", ver.Compiler)
	fmt.Printf("  Platform: %s\n", ver.Platform)
	if detail {
		if err != nil {
			fmt.Printf("Error in getting details from Docker Server: %s\n", err)
		} else {
			fmt.Printf("  OSType: %s\n", detailVer.OSType)
			if detailVer.OSVersion != "" {
				fmt.Printf("  OSVersion: %s\n", detailVer.OSVersion)
			}
			fmt.Printf("  KernelVersion: %s\n", detailVer.KernelVersion)
			fmt.Printf("  DockerHostVersion: %s\n", detailVer.DockerVersion)
			fmt.Printf("  ContainerRuntimes: %s\n", detailVer.ContainerRuntimes)
			fmt.Printf("  DefaultRuntime: %s\n", detailVer.DefaultRuntime)
		}
	}
	return nil
}
