package home

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
)

func (m generalManager) SSHBinaryFile() string {
	return "midi-ssh"
}

func (m generalManager) DownloadOrCacheSSHBinary(url string) error {
	//TODO(gaocegege): Check the checksums.txt first.
	// https://github.com/tensorchord/MIDI/releases/download/v0.0.1-alpha.1/checksums.txt

	filename := filepath.Join(m.CacheDir(), m.SSHBinaryFile())
	logger := logrus.WithFields(logrus.Fields{
		"url":      url,
		"filename": filename,
	})
	out, err := os.Create(filename)

	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to download ssh binary")
	}
	logger.Debugf("downloading %s", filename)

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to download ssh binary")
	}

	if err := os.Chmod(filename, 0755); err != nil {
		return errors.Wrap(err, "failed to chmod ssh binary")
	}

	return nil
}
