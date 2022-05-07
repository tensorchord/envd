// Copyright 2022 The MIDI Authors
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
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
)

const (
	cacheDirName = "cache"
)

type Manager interface {
	HomeDir() string
	CacheDir() string
	ConfigFile() string
}

type generalManager struct {
	homeDir    string
	cacheDir   string
	configFile string

	logger *logrus.Entry
}

var (
	defaultManager *generalManager
	once           sync.Once
)

func Initialize(homeDir, configFile string) error {
	once.Do(func() {
		defaultManager = &generalManager{}
	})
	if err := defaultManager.init(homeDir, configFile); err != nil {
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

func (m generalManager) HomeDir() string {
	return m.homeDir
}

func (m *generalManager) init(homeDir, configFile string) error {
	expandedDir, err := expandHome(homeDir)
	if err != nil {
		return errors.Wrap(err, "failed to expand home dir")
	}
	if err := os.MkdirAll(expandedDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create MIDI home directory")
	}
	m.homeDir = expandedDir

	cacheDir := filepath.Join(expandedDir, cacheDirName)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create MIDI cache directory")
	}
	m.cacheDir = cacheDir

	expandedFilePath, err := expandHome(configFile)
	if err != nil {
		return errors.Wrap(err, "failed to expand config file path")
	}

	_, err = os.Stat(expandedFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.WithField("config", expandedFilePath).Info("Creating config file")
			if _, err := os.Create(expandedFilePath); err != nil {
				return errors.Wrap(err, "failed to create config file")
			}
		} else {
			return errors.Wrap(err, "failed to stat config file")
		}
	}
	m.configFile = expandedFilePath

	m.logger = logrus.WithFields(logrus.Fields{
		"homeDir":  m.homeDir,
		"cacheDir": m.cacheDir,
		"config":   m.configFile,
	})

	m.logger.Debug("home manager initialized")
	return nil
}

func expandHome(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
