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
package builder

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/containerd/containerd/platforms"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func ImageConfigStr(labels map[string]string) (string, error) {
	pl := platforms.Normalize(platforms.DefaultSpec())
	img := v1.Image{
		Config: v1.ImageConfig{
			Labels:     labels,
			WorkingDir: "/",
			Env:        []string{"PATH=" + DefaultPathEnv(pl.OS)},
		},
		Architecture: pl.Architecture,
		// Refer to https://github.com/tensorchord/envd/issues/269#issuecomment-1152944914
		OS: "linux",
		RootFS: v1.RootFS{
			Type: "layers",
		},
	}
	data, err := json.Marshal(img)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DefaultPathEnvUnix is unix style list of directories to search for
// executables. Each directory is separated from the next by a colon
// ':' character .
const DefaultPathEnvUnix = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/conda/bin"

// DefaultPathEnvWindows is windows style list of directories to search for
// executables. Each directory is separated from the next by a colon
// ';' character .
const DefaultPathEnvWindows = "c:\\Windows\\System32;c:\\Windows"

func DefaultPathEnv(os string) string {
	if os == "windows" {
		return DefaultPathEnvWindows
	}
	return DefaultPathEnvUnix
}

// parseOutput parses the output string and returns the output type and destination.
func parseOutput(output string) (string, string, error) {
	if output == "" {
		return "", "", nil
	}

	// Example: type=tar,dest=path
	matched, err := regexp.Match(`^type=[\s\S]+,dest=[\s\S]+$`, []byte(output))
	if err != nil {
		return "", "", errors.Errorf("failed to match output: %v", err)
	}
	if !matched {
		return "", "", errors.Errorf("unsupported format: %s", output)
	}

	fields := strings.Split(output, ",")
	outputMap := make(map[string]string, len(fields))
	for _, field := range fields {
		pair := strings.Split(field, "=")
		outputMap[pair[0]] = pair[1]
	}

	outputType := outputMap["type"]
	outputDest := outputMap["dest"]
	if outputType != "tar" {
		return "", "", errors.Errorf("unsupported output type: %s", outputType)
	}

	return outputType, outputDest, nil
}
