/*
   Copyright The envd Authors.
   Copyright The BuildKit Authors.
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package version

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/types"
)

var (
	// Package is filled at linking time
	Package = "github.com/tensorchord/envd"

	// Revision is filled with the VCS (e.g. git) revision being used to build
	// the program at linking time.
	Revision = ""

	version      = "0.0.0+unknown"
	buildDate    = "1970-01-01T00:00:00Z" // output from `date -u +'%Y-%m-%dT%H:%M:%SZ'`
	gitCommit    = ""                     // output from `git rev-parse HEAD`
	gitTag       = ""                     // output from `git describe --exact-match --tags HEAD` (if clean tree state)
	gitTreeState = ""                     // determined from `git status --porcelain`. either 'clean' or 'dirty'
)

// Version contains envd version information
type Version struct {
	Version      string
	BuildDate    string
	GitCommit    string
	GitTag       string
	GitTreeState string
	GoVersion    string
	Compiler     string
	Platform     string
}

type DetailedVersion struct {
	OSVersion         string
	OSType            string
	KernelVersion     string
	Architecture      string
	DockerVersion     string
	ContainerRuntimes string
	DefaultRuntime    string
}

func (v Version) String() string {
	return v.Version
}

// GetEnvdVersion gets Envd version information
func GetEnvdVersion() string {
	var versionStr string

	if gitCommit != "" && gitTag != "" && gitTreeState == "clean" {
		// if we have a clean tree state and the current commit is tagged,
		// this is an official release.
		versionStr = gitTag
	} else {
		// otherwise formulate a version string based on as much metadata
		// information we have available.
		if strings.HasPrefix(version, "v") {
			versionStr = version
		} else {
			versionStr = "v" + version
		}
		if len(gitCommit) >= 7 {
			versionStr += "+" + gitCommit[0:7]
			if gitTreeState != "clean" {
				versionStr += ".dirty"
			}
		} else {
			versionStr += "+unknown"
		}
	}
	return versionStr
}

func GetRuntimes(info *types.EnvdInfo) string {
	runtimesMap := info.Runtimes
	keys := make([]string, 0, len(runtimesMap))
	for k := range runtimesMap {
		keys = append(keys, k)
	}
	return "[" + strings.Join(keys, ",") + "]"
}

// GetVersion returns the version information
func GetVersion() Version {
	return Version{
		Version:      GetEnvdVersion(),
		BuildDate:    buildDate,
		GitCommit:    gitCommit,
		GitTag:       gitTag,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func GetDetailedVersion(clicontext *cli.Context) (DetailedVersion, error) {
	engine, err := envd.New(clicontext.Context)
	if err != nil {
		return DetailedVersion{}, errors.Wrap(
			err, "failed to create engine for docker server",
		)
	}

	info, err := engine.GetInfo(clicontext.Context)
	if err != nil {
		return DetailedVersion{}, errors.Wrap(
			err, "failed to get detailed version info from docker server",
		)
	}

	return DetailedVersion{
		OSVersion:         info.OSVersion,
		OSType:            info.OSType,
		KernelVersion:     info.KernelVersion,
		DockerVersion:     info.ServerVersion,
		Architecture:      info.Architecture,
		DefaultRuntime:    info.DefaultRuntime,
		ContainerRuntimes: GetRuntimes(info),
	}, nil
}

var (
	reRelease *regexp.Regexp
	reDev     *regexp.Regexp
	reOnce    sync.Once
)

func UserAgent() string {
	version := GetVersion().String()

	reOnce.Do(func() {
		reRelease = regexp.MustCompile(`^(v[0-9]+\.[0-9]+)\.[0-9]+$`)
		reDev = regexp.MustCompile(`^(v[0-9]+\.[0-9]+)\.[0-9]+`)
	})

	if matches := reRelease.FindAllStringSubmatch(version, 1); len(matches) > 0 {
		version = matches[0][1]
	} else if matches := reDev.FindAllStringSubmatch(version, 1); len(matches) > 0 {
		version = matches[0][1] + "-dev"
	}

	return "envd/" + version
}
