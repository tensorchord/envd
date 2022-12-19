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

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/version"
)

var CommandVersion = &cli.Command{
	Name:     "version",
	Category: CategoryOther,
	Aliases:  []string{"v"},
	Usage:    "Print envd version information",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "short",
			Usage:   "Only print the version number",
			Value:   false,
			Aliases: []string{"s"},
		},
		&cli.BoolFlag{
			Name:    "detail",
			Usage:   "Print details about the envd environment",
			Value:   false,
			Aliases: []string{"d"},
		},
		&cli.StringFlag{
			Name:     "format",
			Usage:    "Format of output, could be \"json\" or \"table\"",
			Aliases:  []string{"f"},
			Value:    "table",
			Required: false,
			Action: func(clicontext *cli.Context, v string) error {
				switch v {
				case
					"table",
					"json":
					return nil
				}
				return errors.Errorf("Argument format only allows \"json\" and \"table\", found %v", v)
			},
		},
	},
	Action: printVersion,
}

func printVersion(clicontext *cli.Context) error {
	format := clicontext.String("format")
	switch format {
	case "table":
		return printTableVersion(clicontext)
	case "json":
		return printJsonVersion(clicontext)
	}
	return nil
}

func printTableVersion(clicontext *cli.Context) error {
	short := clicontext.Bool("short")
	detail := clicontext.Bool("detail")
	ver := version.GetVersion()
	detailVer, err := getDetailedVersion(clicontext)
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

type versionJson struct {
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

func printJsonVersion(clicontext *cli.Context) error {
	short := clicontext.Bool("short")
	detail := clicontext.Bool("detail")
	ver := version.GetVersion()
	detailVer, err := getDetailedVersion(clicontext)
	output := versionJson{
		Envd: version.GetVersion().String(),
	}
	if short {
		return PrintJson(output)
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
	return PrintJson(output)
}

func getDetailedVersion(clicontext *cli.Context) (detailedVersion, error) {
	context, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return detailedVersion{}, errors.Wrap(err, "failed to get the current context")
	}
	opt := envd.Options{
		Context: context,
	}
	engine, err := envd.New(clicontext.Context, opt)
	if err != nil {
		return detailedVersion{}, errors.Wrap(
			err, "failed to create engine for docker server",
		)
	}

	info, err := engine.GetInfo(clicontext.Context)
	if err != nil {
		return detailedVersion{}, errors.Wrap(
			err, "failed to get detailed version info from docker server",
		)
	}

	return detailedVersion{
		OSVersion:         info.OSVersion,
		OSType:            info.OSType,
		KernelVersion:     info.KernelVersion,
		DockerVersion:     info.ServerVersion,
		Architecture:      info.Architecture,
		DefaultRuntime:    info.DefaultRuntime,
		ContainerRuntimes: GetRuntimes(info),
	}, nil
}

type detailedVersion struct {
	OSVersion         string
	OSType            string
	KernelVersion     string
	Architecture      string
	DockerVersion     string
	ContainerRuntimes string
	DefaultRuntime    string
}

func GetRuntimes(info *types.EnvdInfo) string {
	runtimesMap := info.Runtimes
	keys := make([]string, 0, len(runtimesMap))
	for k := range runtimesMap {
		keys = append(keys, k)
	}
	return "[" + strings.Join(keys, ",") + "]"
}

func PrintJson(v any) error {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	// avoid escaped from <none> into "\u003cnone\u003e"
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}
	fmt.Print(buffer.String())
	return nil
}
