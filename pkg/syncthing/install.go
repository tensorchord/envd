package syncthing

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hashicorp/go-getter"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

// TODO: syncthing installation versioning
func getSyncthingVersion() string {
	return "1.22.2"
}

func getSyncthingInstallPath() string {
	return fmt.Sprintf("%s/bin", fileutil.DefaultCacheDir)
}

func getSyncthingBinPath() string {
	return fmt.Sprintf("%s/syncthing", getSyncthingInstallPath())
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

	return "", fmt.Errorf("%s-%s is not a supported platform for syncthing", os, arch)
}

func getSyncthingDownloadFolderName(os, arch, version string) string {
	return fmt.Sprintf("syncthing-%[1]s-%[2]s-v%[3]s", os, arch, version)
}

func IsInstalled() bool {
	if _, err := os.Stat(getSyncthingBinPath()); err != nil {
		return false
	}
	return true
}

func Install() error {
	logrus.Info("Installing syncthing")
	if IsInstalled() {
		logrus.Info("Syncthing is already installed, skipping installation")
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

	logrus.Info("Downloading syncthing from ", downloadUrl)
	client := &getter.Client{
		Src:  downloadUrl,
		Dst:  getSyncthingInstallPath(),
		Mode: getter.ClientModeDir,
	}

	if err := client.Get(); err != nil {
		return fmt.Errorf("failed to download syncthing from url %s: %s", client.Src, err)
	}

	var downloadFolder = fmt.Sprintf("%s/%s", getSyncthingInstallPath(), getSyncthingDownloadFolderName(operatingSystem, arch, version))

	err = os.Rename(fmt.Sprintf("%s/syncthing", downloadFolder), getSyncthingBinPath())
	if err != nil {
		return fmt.Errorf("failed to move syncthing binary: %s", err)
	}

	err = os.RemoveAll(downloadFolder)
	if err != nil {
		return fmt.Errorf("failed to remove syncthing download folder: %s", err)
	}

	if err := os.Chmod(getSyncthingBinPath(), 0755); err != nil {
		return fmt.Errorf("failed to set syncthing binary permissions: %s", err)
	}

	logrus.Info("Syncthing installed successfully!")

	return nil
}
