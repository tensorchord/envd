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

package vscode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"
)

const (
	PLATFORM_WIN32_X64    = "win32-x64"
	PLATFORM_WIN32_IA32   = "win32-ia32"
	PLATFORM_WIN32_ARM64  = "win32-arm64"
	PLATFORM_LINUX_X64    = "linux-x64"
	PLATFORM_LINUX_ARM64  = "linux-arm64"
	PLATFORM_LINUX_ARMHF  = "linux-armhf"
	PLATFORM_DARWIN_X64   = "darwin-x64"
	PLATFORM_DARWIN_ARM64 = "darwin-arm64"
	PLATFORM_ALPINE_X64   = "alpine-x64"
)

func ConvertLLBPlatform(platform *v1.Platform) (string, error) {
	// Convert opencontainers style platform to VSCode extension style platform.
	switch platform.OS {
	case "windows":
		switch platform.Architecture {
		case "amd64":
			return PLATFORM_WIN32_X64, nil
		case "386":
			return PLATFORM_WIN32_IA32, nil
		case "arm64":
			return PLATFORM_WIN32_ARM64, nil
		}
	case "linux":
		switch platform.Architecture {
		case "amd64":
			return PLATFORM_LINUX_X64, nil
		case "arm64":
			return PLATFORM_LINUX_ARM64, nil
		case "arm":
			return PLATFORM_LINUX_ARMHF, nil
		}
	case "darwin":
		switch platform.Architecture {
		case "amd64":
			return PLATFORM_DARWIN_X64, nil
		case "arm64":
			return PLATFORM_DARWIN_ARM64, nil
		}
	case "alpine":
		switch platform.Architecture {
		case "amd64":
			return PLATFORM_ALPINE_X64, nil
		}
	}

	return "", errors.Errorf("unsupported platform: %s/%s", platform.OS, platform.Architecture)
}

func GetLatestVersionURL(p Plugin) (string, error) {
	// Auto-detect the version.
	// Refer to https://github.com/tensorchord/envd/issues/161#issuecomment-1129475975
	latestURL := fmt.Sprintf(vendorOpenVSXTemplate, p.Publisher, p.Extension, p.Platform)
	resp, err := http.Get(latestURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to get latest version")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("failed to get latest version: %s", resp.Status)
	}
	jsonResp := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return "", errors.Wrap(err, "failed to decode response")
	}
	if jsonResp["downloads"] == nil {
		return "", errors.New("failed to get latest version: no downloads")
	}
	downloads := jsonResp["downloads"].(map[string]interface{})
	if downloads["universal"] != nil {
		return downloads["universal"].(string), nil
	}
	if downloads[p.Platform] == nil {
		return "", errors.Errorf("failed to get latest version: no target platform %s", p.Platform)
	}
	return downloads[p.Platform].(string), nil
}

func ParsePlugin(p string) (*Plugin, error) {
	indexPublisher := strings.Index(p, ".")
	if indexPublisher == -1 {
		return nil, errors.New("invalid publisher")
	}
	publisher := p[:indexPublisher]

	indexExtension := strings.LastIndex(p[indexPublisher:], "-")
	if indexExtension == -1 {
		extension := p[indexPublisher+1:]
		logrus.WithFields(logrus.Fields{
			"publisher": publisher,
			"extension": extension,
		}).Debug("vscode plugin is parsed without version")
		return &Plugin{
			Publisher: publisher,
			Extension: extension,
		}, nil
	}

	indexExtension = indexPublisher + indexExtension
	extension := p[indexPublisher+1 : indexExtension]
	version := p[indexExtension+1:]
	if _, err := strconv.Atoi(version[0:1]); err != nil {
		extension := p[indexPublisher+1:]
		logrus.WithFields(logrus.Fields{
			"publisher": publisher,
			"extension": extension,
		}).Debug("vscode plugin is parsed without version")
		return &Plugin{
			Publisher: publisher,
			Extension: extension,
		}, nil
	}
	logrus.WithFields(logrus.Fields{
		"publisher": publisher,
		"extension": extension,
		"version":   version,
	}).Debug("vscode plugin is parsed")
	return &Plugin{
		Publisher: publisher,
		Extension: extension,
		Version:   &version,
	}, nil
}
