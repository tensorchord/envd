// Copyright 2025 The envd Authors
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

package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/lang/ir"
)

// https://github.com/openai/codex
const (
	codexDefaultVersion = "0.98.0"
	codexReleaseUser    = "openai"
	codexReleaseRepo    = "codex"
)

func getLatestVersion(user, repo string) (string, error) {
	latestURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", user, repo)
	req, err := http.NewRequest(http.MethodGet, latestURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create request")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to get latest release")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("failed to get latest release: %s", resp.Status)
	}

	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", errors.Wrap(err, "failed to decode response")
	}
	if payload.TagName == "" {
		return "", errors.New("failed to get latest release: empty tag name")
	}

	version := strings.TrimPrefix(payload.TagName, "rust-v")
	version = strings.TrimPrefix(version, "v")
	return version, nil
}

func (g generalGraph) installAgentCodex(root llb.State, agent ir.CodeAgent) llb.State {
	base := llb.Image(curlImage)
	version := codexDefaultVersion
	if agent.Version != nil {
		version = *agent.Version
	} else {
		latestVersion, err := getLatestVersion(codexReleaseUser, codexReleaseRepo)
		if err != nil {
			logrus.WithError(err).WithField("default", codexDefaultVersion).Debug("failed to resolve latest codex version")
		} else {
			version = latestVersion
		}
	}
	logrus.WithField("codex_version", version).Debug("parse the agent version")
	builder := base.Run(
		llb.Shlexf(`sh -c "wget -qO- https://github.com/openai/codex/releases/download/rust-v%s/codex-$(uname -m)-unknown-linux-musl.tar.gz | tar -xz -C /tmp || exit 1"`, version),
		llb.WithCustomNamef("[internal] download codex %s", version),
	).Run(
		llb.Shlex(`sh -c "mv /tmp/codex-$(uname -m)-unknown-linux-musl /tmp/codex"`),
		llb.WithCustomNamef("[internal] prepare codex %s", version),
	).Root()
	root = root.File(
		llb.Copy(builder, "/tmp/codex", "/usr/bin/codex"),
		llb.WithCustomName("[internal] install codex"),
	)
	return root
}
