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
	"os"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
)

type Manager interface {
	CacheDir() string
	ConfigFile() string
}

type generalManager struct {
	cacheDir   string
	configFile string

	logger *logrus.Entry
}

var (
	defaultManager *generalManager
	once           sync.Once
)

func Initialize() error {
	once.Do(func() {
		defaultManager = &generalManager{}
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

func (m *generalManager) init() error {
	// Create $XDG_CONFIG_HOME/envd/config.envd
	config, err := xdg.ConfigFile("envd/config.envd")
	if err != nil {
		return errors.Wrap(err, "failed to get config file")
	}

	_, err = os.Stat(config)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.WithField("config", config).Info("Creating config file")
			if _, err := os.Create(config); err != nil {
				return errors.Wrap(err, "failed to create config file")
			}
		} else {
			return errors.Wrap(err, "failed to stat config file")
		}
	}
	m.configFile = config

	// Create $XDG_CACHE_HOME/envd
	_, err = xdg.CacheFile("envd/cache")
	if err != nil {
		return errors.Wrap(err, "failed to get cache")
	}
	m.cacheDir = filepath.Join(xdg.CacheHome, "envd")

	m.logger = logrus.WithFields(logrus.Fields{
		"cacheDir": m.cacheDir,
		"config":   m.configFile,
	})

	m.logger.Debug("home manager initialized")
	return nil
}
