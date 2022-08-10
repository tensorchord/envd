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

	"go.starlark.net/starlark"
)

type EnvdManagedDataSource struct {
	name string
}

func (d EnvdManagedDataSource) Type() string {
	return "Envd Managed Data Source"
}

func (d EnvdManagedDataSource) String() string {
	return "Envd Managed Data Source"
}

func (d EnvdManagedDataSource) Freeze()              {}
func (d EnvdManagedDataSource) Truth() starlark.Bool { return false }
func (d EnvdManagedDataSource) Hash() (uint32, error) {
	return hashString(fmt.Sprintf("envd://%s", d.name)), nil
}

func (d *EnvdManagedDataSource) Init() {
	panic("not implemented") // TODO: Implement
}

func (d *EnvdManagedDataSource) GetHostDir() string {
	panic("not implemented") // TODO: Implement
}

func NewEnvdManagedDataSource(name string) *EnvdManagedDataSource {
	return &EnvdManagedDataSource{name: name}
}
