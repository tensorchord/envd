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
)

var (
	// Package is filled at linking time
	Package = "github.com/tensorchord/envd"

	// Revision is filled with the VCS (e.g. git) revision being used to build
	// the program at linking time.
	Revision = ""

	version         = "0.0.0+unknown"
	buildDate       = "1970-01-01T00:00:00Z" // output from `date -u +'%Y-%m-%dT%H:%M:%SZ'`
	gitCommit       = ""                     // output from `git rev-parse HEAD`
	gitTag          = ""                     // output from `git describe --exact-match --tags HEAD` (if clean tree state)
	gitTreeState    = ""                     // determined from `git status --porcelain`. either 'clean' or 'dirty'
	developmentFlag = "false"
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

func (v Version) String() string {
	return v.Version
}

// GetVersionForImageTag gets the version for an image tag.
func GetVersionForImageTag() string {
	if gitTag != "" {
		return gitTag
	}
	// Empty version tag only appears in dev built.
	// Set to `latest` if so.
	return "latest"
}

// SetGitTagForE2ETest sets the gitTag for test purpose.
func SetGitTagForE2ETest(tag string) {
	gitTag = tag
}

// GetEnvdVersion gets Envd version information
func GetEnvdVersion() string {
	var versionStr string

	if gitCommit != "" && gitTag != "" &&
		gitTreeState == "clean" && developmentFlag == "false" {
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
