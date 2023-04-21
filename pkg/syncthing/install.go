// Copyright 2023 The envd Authors
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

package syncthing

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/go-getter"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func getSyncthingVersion() string {
	// TODO: Better versioning
	return "1.22.2"
}

func getSyncthingInstallPath() string {
	return filepath.Join(fileutil.DefaultCacheDir, "bin")
}

func GetSyncthingBinPath() string {
	return filepath.Join(getSyncthingInstallPath(), "syncthing")
}

func getSyncthingDownloadURL(os, arch, version string) (string, error) {
	// TODO: double check os/arch support
	fileExtension := "tar.gz"
	if os == "windows" || os == "macos" {
		fileExtension = "zip"
	}

	downloadUrl := "https://github.com/syncthing/syncthing/releases/download/v%[3]s/syncthing-%[1]s-%[2]s-v%[3]s.%[4]s"
	switch os {
	case "linux":
		switch arch {
		case "amd64", "arm64":
			return fmt.Sprintf(downloadUrl, os, arch, version, fileExtension), nil
		}
	case "macos":
		switch arch {
		case "amd64", "arm64":
			return fmt.Sprintf(downloadUrl, os, arch, version, fileExtension), nil
		}
	case "windows":
		switch arch {
		case "amd64":
			return fmt.Sprintf(downloadUrl, os, arch, version, fileExtension), nil
		}
	}

	return "", errors.New(fmt.Sprintf("%s-%s is not a supported platform for syncthing", os, arch))
}

func getSyncthingDownloadFolderName(os, arch, version string) string {
	return fmt.Sprintf("syncthing-%[1]s-%[2]s-v%[3]s", os, arch, version)
}

func IsInstalled() bool {
	if _, err := os.Stat(GetSyncthingBinPath()); err != nil {
		return false
	}
	return true
}

func (s *Syncthing) CleanupSyncthing() error {
	logrus.Debug("Cleaning up syncthing")

	err := os.RemoveAll(s.HomeDirectory)
	if err != nil {
		return errors.Wrap(err, "failed to remove syncthing config file: ")
	}

	return nil
}

func InstallSyncthing() error {
	logrus.Debug("Installing syncthing")
	if IsInstalled() {
		logrus.Debug("Syncthing is already installed, skipping installation")
		return nil
	}

	var operatingSystem = runtime.GOOS
	var arch = runtime.GOARCH
	var version = getSyncthingVersion()
	if operatingSystem == "darwin" {
		operatingSystem = "macos"
	}

	downloadUrl, err := getSyncthingDownloadURL(operatingSystem, arch, version)
	if err != nil {
		return err
	}

	logrus.Debug("Downloading syncthing from ", downloadUrl)
	client := &getter.Client{
		Src:  downloadUrl,
		Dst:  getSyncthingInstallPath(),
		Mode: getter.ClientModeDir,
	}

	if err := client.Get(); err != nil {
		return err
	}

	var downloadFolder = fmt.Sprintf("%s/%s", getSyncthingInstallPath(), getSyncthingDownloadFolderName(operatingSystem, arch, version))

	err = os.Rename(fmt.Sprintf("%s/syncthing", downloadFolder), GetSyncthingBinPath())
	if err != nil {
		return err
	}

	err = os.RemoveAll(downloadFolder)
	if err != nil {
		return err
	}

	if err := os.Chmod(GetSyncthingBinPath(), 0755); err != nil {
		return err
	}

	logrus.Info("Syncthing installed successfully!")

	return nil
}
