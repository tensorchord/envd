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

	"github.com/cockroachdb/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

const (
	cacheKey = "oh-my-zsh"
)

//go:embed install.sh
var installScript string

//go:embed zshrc
var zshrc string

type Manager interface {
	ZSHRC() string
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

func (m generalManager) ZSHRC() string {
	return zshrc
}

func (m generalManager) DownloadOrCache() (bool, error) {
	if home.GetManager().Cached(cacheKey) {
		logrus.WithFields(logrus.Fields{
			"cache-dir": m.OHMyZSHDir(),
		}).Debug("oh-my-zsh already exists in cache")
		return true, nil
	}
	url := "https://github.com/ohmyzsh/ohmyzsh.git"
	l := logrus.WithFields(logrus.Fields{
		"cache-dir": m.OHMyZSHDir(),
		"URL":       url,
	})

	// Cleanup the cache dir.
	if fileutil.RemoveAll(m.OHMyZSHDir()) != nil {
		return false, errors.New("failed to remove oh-my-zsh dir")
	}
	l.Debug("cache miss, downloading oh-my-zsh")

	// Init the git repository.
	repo, err := git.PlainInit(m.OHMyZSHDir(), false)
	if err != nil {
		return false, errors.Wrap(err, "failed to init oh-my-zsh repo")
	}
	cfg, err := repo.Config()
	if err != nil {
		return false, errors.Wrap(err, "failed to get repo config")
	}
	// Refer to https://github.com/tensorchord/envd/issues/183#issuecomment-1148113323
	cfg.Raw.AddOption("core", "", "eol", "lf")
	cfg.Raw.AddOption("core", "", "autocrlf", "false")
	cfg.Raw.AddOption("core", "", "filemode", "true")
	cfg.Raw.AddOption("core", "", "logallrefupdates", "true")
	cfg.Raw.AddOption("core", "", "repositoryformatversion", "0")
	cfg.Raw.AddOption("core", "", "bare", "false")
	cfg.Raw.AddOption("fsck", "", "zeroPaddedFilemode", "ignore")
	cfg.Raw.AddOption("fetch", "fsck", "zeroPaddedFilemode", "ignore")
	cfg.Raw.AddOption("receive", "fsck", "zeroPaddedFilemode", "ignore")
	cfg.Remotes["origin"] = &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
		Fetch: []config.RefSpec{
			config.RefSpec("+refs/heads/master:refs/remotes/origin/master"),
		},
	}
	cfg.Branches["master"] = &config.Branch{
		Name:   "master",
		Remote: "origin",
		Merge:  "refs/heads/master",
	}

	if err := cfg.Validate(); err != nil {
		return false, errors.Wrap(err, "failed to validate config")
	}
	if err := repo.SetConfig(cfg); err != nil {
		return false, errors.Wrap(err, "failed to set config")
	}

	if err := repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/master:refs/remotes/origin/master"),
		},
		Depth: 1,
	}); err != nil {
		return false, errors.Wrap(err, "failed to fetch oh-my-zsh")
	}
	wktree, err := repo.Worktree()
	if err != nil {
		return false, errors.Wrap(err, "failed to get worktree")
	}
	if err := wktree.Checkout(&git.CheckoutOptions{
		Branch: "refs/remotes/origin/master",
	}); err != nil {
		return false, errors.Wrap(err, "failed to checkout master")
	}

	if err := home.GetManager().MarkCache(cacheKey, true); err != nil {
		return false, errors.Wrap(err, "failed to update cache status")
	}
	l.Debug("oh-my-zsh is downloaded")
	return false, nil
}

func (m generalManager) OHMyZSHDir() string {
	return filepath.Join(home.GetManager().CacheDir(), "oh-my-zsh")
}
