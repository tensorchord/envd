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

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

type cacheManager interface {
	CacheDir() string
	MarkCache(string, bool) error
	Cached(string) bool
	CleanCache() error
}

func (m *generalManager) initCache() error {
	// Create $HOME/.cache/envd/
	m.cacheDir = fileutil.DefaultCacheDir

	cacheStatusFile, err := fileutil.CacheFile("cache.status")
	if err != nil {
		return errors.Wrap(err, "failed to get cache.status file path")
	}
	m.cacheStatusFile = cacheStatusFile
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
	return nil
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

func (m generalManager) CacheDir() string {
	return m.cacheDir
}

func (m generalManager) CleanCache() error {
	if m.cacheDir == "" {
		return nil
	}
	logrus.Debug("cleaning up host cache directory")
	return os.RemoveAll(m.cacheDir)
}
