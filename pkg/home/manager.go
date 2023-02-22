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
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/driver/docker"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
)

type Manager interface {
	configManager
	contextManager
	cacheManager
	dataManager
	authManager
}

type generalManager struct {
	cacheDir        string
	cacheStatusFile string
	configFile      string
	contextFile     string
	authFile        string

	// TODO(gaocegege): Abstract CacheManager.
	cacheMap map[string]bool
	context  types.EnvdContext
	auth     types.EnvdAuth

	logger *logrus.Entry
}

var (
	defaultManager *generalManager
	once           sync.Once
)

func Initialize() error {
	builder := types.BuilderTypeDocker
	dockerVersion, err := docker.GetDockerVersion()
	if err == nil && dockerVersion > 22 {
		builder = types.BuilderTypeMoby
	}
	once.Do(func() {
		defaultManager = &generalManager{
			cacheMap: make(map[string]bool),
			context: types.EnvdContext{
				Current: "default",
				Contexts: []types.Context{
					{
						Name:           "default",
						Builder:        builder,
						BuilderAddress: "envd_buildkitd",
						Runner:         types.RunnerTypeDocker,
						RunnerAddress:  nil,
					},
				},
			},
		}
	})
	return defaultManager.init()
}

func GetManager() Manager {
	return defaultManager
}

func (m *generalManager) init() error {
	if err := m.initConfig(); err != nil {
		return errors.Wrap(err, "failed to initialize config")
	}

	if err := m.initContext(); err != nil {
		return errors.Wrap(err, "failed to initialize context")
	}

	if err := m.initCache(); err != nil {
		return errors.Wrap(err, "failed to initialize cache")
	}

	if err := m.initAuth(); err != nil {
		return errors.Wrap(err, "failed to initialize auth")
	}

	if err := sshconfig.GenerateKeys(); err != nil {
		return errors.Wrap(err, "failed to generate ssh key")
	}

	m.logger = logrus.WithFields(logrus.Fields{
		"cache-dir":    m.cacheDir,
		"config-file":  m.configFile,
		"cache-status": m.cacheStatusFile,
		"context-file": m.contextFile,
		"cache-map":    m.cacheMap,
		"context":      m.context,
	})

	m.logger.Debug("home manager initialized")
	return nil
}
