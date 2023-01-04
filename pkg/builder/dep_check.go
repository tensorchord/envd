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
	"context"
	"os"

	"github.com/tensorchord/envd/pkg/driver/docker"
)

func (b generalBuilder) checkIfNeedBuild(ctx context.Context) bool {
	if b.graph.GetHTTP() != nil {
		return true
	}
	depsFiles := []string{
		b.ConfigFilePath,
	}
	depsFiles = b.GetDepsFilesHandler(depsFiles)
	isUpdated, err := b.checkDepsFileUpdate(ctx, b.Tag, b.ManifestFilePath, depsFiles)
	if err != nil {
		b.logger.Debugf("failed to check manifest update: %s", err)
	}
	if !isUpdated {
		b.logger.Infof("manifest is not updated, skip building")
		return false
	}
	return true
}

// nolint:unparam
func (b generalBuilder) checkDepsFileUpdate(ctx context.Context, tag string, manifest string, deps []string) (bool, error) {
	dockerClient, err := docker.NewClient(ctx)
	if err != nil {
		return true, err
	}

	image, err := dockerClient.GetImageWithCacheHashLabel(ctx, tag, b.manifestCodeHash)
	if err != nil {
		return true, err
	}
	imageCreatedTime := image.Created

	latestTimestamp := int64(0)
	for _, dep := range deps {
		file, err := os.Stat(dep)
		if err != nil {
			return true, err
		}
		modifiedtime := file.ModTime().Unix()
		// Only need to use the latest modified time
		if modifiedtime > latestTimestamp {
			latestTimestamp = modifiedtime
		}
	}
	if latestTimestamp > imageCreatedTime {
		return true, nil
	}
	return false, nil
}
