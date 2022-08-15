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

package data

import (
	"fmt"

	"github.com/tensorchord/envd/pkg/home"

	"go.starlark.net/starlark"
)

type EnvdManagedDataSource struct {
	name        string
	hostDataDir string
}

func (e *EnvdManagedDataSource) Init() error {
	manager := home.GetManager()
	hostDataDir, err := manager.InitDataDir(e.name)
	if err != nil {
		return err
	}
	e.hostDataDir = hostDataDir
	return nil
}

func (e *EnvdManagedDataSource) GetHostDir() (string, error) {
	return e.hostDataDir, nil
}

func (e *EnvdManagedDataSource) Type() string {
	return "envd managed data source"
}

func (e *EnvdManagedDataSource) Hash() (uint32, error) {
	return starlark.String(fmt.Sprintf("envd://%s", e.name)).Hash()
}

func NewEnvdManagedDataSource(name string) *EnvdManagedDataSource {
	return &EnvdManagedDataSource{name: name}
}
