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
	"github.com/cockroachdb/errors"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

type configManager interface {
	ConfigFile() string
}

func (m *generalManager) initConfig() error {
	// Create $HOME/.config/envd/config.envd
	config, err := fileutil.ConfigFile("config.envd")
	if err != nil {
		return errors.Wrap(err, "failed to get config file")
	}

	if err := fileutil.CreateIfNotExist(config); err != nil {
		return errors.Wrap(err, "failed to create config file")
	}
	m.configFile = config
	return nil
}

func (m *generalManager) ConfigFile() string {
	return m.configFile
}
