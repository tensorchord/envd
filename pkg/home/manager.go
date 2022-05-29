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

package home

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

type Manager interface {
	CacheDir() string
	MarkCache(string, bool) error
	Cached(string) bool
	ConfigFile() string
}

type generalManager struct {
	cacheDir        string
	cacheStatusFile string
	configFile      string

	// TODO(gaocegege): Abstract CacheManager.
	cacheMap map[string]bool

	logger *logrus.Entry
}

var (
	defaultManager *generalManager
	once           sync.Once
)

func Initialize() error {
	once.Do(func() {
		defaultManager = &generalManager{
			cacheMap: make(map[string]bool),
		}
	})
	if err := defaultManager.init(); err != nil {
		return err
	}
	return nil
}

func GetManager() Manager {
	return defaultManager
}

func (m generalManager) CacheDir() string {
	return m.cacheDir
}

func (m generalManager) ConfigFile() string {
	return m.configFile
}

func (m generalManager) MarkCache(key string, cached bool) error {
	m.cacheMap[key] = cached
	return m.dumpCacheStatus()
}

func (m generalManager) Cached(key string) bool {
	return m.cacheMap[key]
}

func (m *generalManager) dumpCacheStatus() error {
	file, err := os.Create(m.cacheStatusFile)
	if err != nil {
		return errors.Wrap(err, "failed to create cache status file")
	}
	defer file.Close()

	e := gob.NewEncoder(file)
	if err := e.Encode(m.cacheMap); err != nil {
		return errors.Wrap(err, "failed to encode cache map")
	}
	return nil
}

func (m *generalManager) init() error {
	// Create $XDG_CONFIG_HOME/envd/config.envd
	config, err := xdg.ConfigFile("envd/config.envd")
	if err != nil {
		return errors.Wrap(err, "failed to get config file")
	}

	if err := fileutil.CreateIfNotExist(config); err != nil {
		return errors.Wrap(err, "failed to create config file")
	}
	m.configFile = config

	// Create $XDG_CACHE_HOME/envd
	_, err = xdg.CacheFile("envd/cache")
	if err != nil {
		return errors.Wrap(err, "failed to get cache")
	}
	m.cacheDir = filepath.Join(xdg.CacheHome, "envd")

	m.cacheStatusFile = filepath.Join(m.cacheDir, "cache.status")
	_, err = os.Stat(m.cacheStatusFile)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.WithField("filename", m.cacheStatusFile).Debug("Creating file")
			file, err := os.Create(m.cacheStatusFile)
			if err != nil {
				return errors.Wrap(err, "failed to create file")
			}
			err = file.Close()
			if err != nil {
				return errors.Wrap(err, "failed to close file")
			}
			if err := m.dumpCacheStatus(); err != nil {
				return errors.Wrap(err, "failed to dump cache status")
			}
		} else {
			return errors.Wrap(err, "failed to stat file")
		}
	}

	file, err := os.Open(m.cacheStatusFile)
	if err != nil {
		return errors.Wrap(err, "failed to open cache status file")
	}
	defer file.Close()
	e := gob.NewDecoder(file)
	if err := e.Decode(&m.cacheMap); err != nil {
		return errors.Wrap(err, "failed to decode cache map")
	}

	// Generate SSH keys when init
	if err := sshconfig.GenerateKeys(); err != nil {
		return errors.Wrap(err, "failed to generate ssh key")
	}

	m.logger = logrus.WithFields(logrus.Fields{
		"cache-dir":    m.cacheDir,
		"config":       m.configFile,
		"cache-status": m.cacheStatusFile,
	})

	m.logger.Debug("home manager initialized")
	return nil
}
