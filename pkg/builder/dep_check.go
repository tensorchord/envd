package builder

import (
	"context"
	"os"

	"github.com/tensorchord/envd/pkg/docker"
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
