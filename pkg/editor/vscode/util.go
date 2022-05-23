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
	"github.com/sirupsen/logrus"
)

func GetLatestVersionURL(p Plugin) (string, error) {
	// Auto-detect the version.
	// Refer to https://github.com/tensorchord/envd/issues/161#issuecomment-1129475975
	latestURL := fmt.Sprintf(vendorOpenVSXTemplate, p.Publisher, p.Extension)
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
	if jsonResp["files"] == nil {
		return "", errors.New("failed to get latest version: no files")
	}
	files := jsonResp["files"].(map[string]interface{})
	if files["download"] == nil {
		return "", errors.New("failed to get latest version: no download url")
	}
	return files["download"].(string), nil
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
