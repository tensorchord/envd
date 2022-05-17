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

package shell

import (
	_ "embed"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

//go:embed install.sh
var installScript string

type Manager interface {
	InstallScript() string
	DownloadOrCache() (bool, error)
	OHMyZSHDir() string
}

type generalManager struct {
}

func NewManager() Manager {
	return &generalManager{}
}

func (m generalManager) InstallScript() string {
	return installScript
}

func (m generalManager) DownloadOrCache() (bool, error) {
	if home.GetManager().Cached("oh-my-zsh") {
		logrus.WithFields(logrus.Fields{
			"cache-dir": m.OHMyZSHDir(),
		}).Debug("oh-my-zsh already exists in cache")
		return true, nil
	}
	url := "https://github.com/ohmyzsh/ohmyzsh"
	l := logrus.WithFields(logrus.Fields{
		"cache-dir": m.OHMyZSHDir(),
		"URL":       url,
	})

	// Cleanup the cache dir.
	if fileutil.RemoveAll(m.OHMyZSHDir()) != nil {
		return false, errors.New("failed to remove oh-my-zsh dir")
	}
	l.Debug("cache miss, downloading oh-my-zsh")
	_, err := git.PlainClone(m.OHMyZSHDir(), false, &git.CloneOptions{
		URL:   url,
		Depth: 1,
	})
	if err != nil {
		return false, err
	}

	if err := home.GetManager().MarkCache("oh-my-zsh", true); err != nil {
		return false, errors.Wrap(err, "failed to update cache status")
	}
	l.Debug("oh-my-zsh is downloaded")
	return false, nil
}

func (m generalManager) OHMyZSHDir() string {
	return filepath.Join(home.GetManager().CacheDir(), "oh-my-zsh")
}
