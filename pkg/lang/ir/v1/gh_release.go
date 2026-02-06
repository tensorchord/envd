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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var (
	githubAPIBaseURL      = "https://api.github.com"
	latestVersionCacheTTL = time.Hour
	maxReleaseNum         = 10
)

type cacheEntry struct {
	Version   string    `json:"version"`
	ExpiresAt time.Time `json:"expires_at"`
}

func getLatestReleaseVersion(user, repo string) (string, error) {
	now := time.Now()
	if version, ok, err := readLatestVersionCache(user, repo, now); err != nil {
		logrus.WithError(err).Debug("failed to read latest release cache")
	} else if ok {
		return version, nil
	}

	latestURL := fmt.Sprintf("%s/repos/%s/%s/releases", githubAPIBaseURL, user, repo)
	req, err := http.NewRequest(http.MethodGet, latestURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create request")
	}
	q := req.URL.Query()
	q.Set("per_page", strconv.Itoa(maxReleaseNum))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "envd")
	if token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
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

	var releases []struct {
		TagName    string `json:"tag_name"`
		PreRelease bool   `json:"prerelease"`
		Draft      bool   `json:"draft"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", errors.Wrap(err, "failed to decode response")
	}
	if len(releases) == 0 {
		return "", errors.New("failed to get latest release: empty response")
	}

	version := ""
	for _, release := range releases {
		if release.Draft || release.PreRelease || release.TagName == "" {
			continue
		}
		version = release.TagName
		break
	}
	if version == "" {
		return "", errors.Newf("failed to get latest release: no stable release found in the %d releases", maxReleaseNum)
	}
	if err := writeLatestVersionCache(user, repo, cacheEntry{
		Version:   version,
		ExpiresAt: now.Add(latestVersionCacheTTL),
	}); err != nil {
		logrus.WithError(err).Debug("failed to write latest release cache")
	}
	return version, nil
}

func readLatestVersionCache(user, repo string, now time.Time) (string, bool, error) {
	cachePath, err := cacheFilePath(user, repo)
	if err != nil {
		return "", false, err
	}
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, errors.Wrap(err, "failed to read cache file")
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", false, errors.Wrap(err, "failed to decode cache file")
	}
	if entry.Version == "" || now.After(entry.ExpiresAt) {
		return "", false, nil
	}
	return entry.Version, true, nil
}

func writeLatestVersionCache(user, repo string, entry cacheEntry) error {
	cachePath, err := cacheFilePath(user, repo)
	if err != nil {
		return err
	}
	dir := filepath.Dir(cachePath)
	tmp, err := os.CreateTemp(dir, "github-release-*.tmp")
	if err != nil {
		return errors.Wrap(err, "failed to create temp cache file")
	}
	defer func() {
		_ = os.Remove(tmp.Name())
	}()
	if err := json.NewEncoder(tmp).Encode(entry); err != nil {
		_ = tmp.Close()
		return errors.Wrap(err, "failed to encode cache file")
	}
	if err := tmp.Close(); err != nil {
		return errors.Wrap(err, "failed to close temp cache file")
	}
	if err := os.Rename(tmp.Name(), cachePath); err != nil {
		return errors.Wrap(err, "failed to move cache file")
	}
	return nil
}

func cacheFilePath(user, repo string) (string, error) {
	name := fmt.Sprintf("github-release-%s-%s.json", sanitizeCacheComponent(user), sanitizeCacheComponent(repo))
	return fileutil.CacheFile(name)
}

func sanitizeCacheComponent(component string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_")
	return replacer.Replace(component)
}
