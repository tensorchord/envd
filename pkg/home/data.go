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

	"github.com/cockroachdb/errors"
)

type dataManager interface {
	InitDataDir(name string) (string, error)
}

func (m *generalManager) InitDataDir(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home dir")
	}
	newDataDir := filepath.Join(home, ".envd", "data", name)
	err = os.Mkdir(newDataDir, 0644)
	if err != nil {
		return "", errors.Wrap(err, "failed to create data dir")
	}
	return newDataDir, nil
}
