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
	"go.starlark.net/starlark"

	envddata "github.com/tensorchord/envd/pkg/data"
)

type DataSourceValue struct {
	source envddata.DataSource
}

func (d DataSourceValue) Init() error {
	return d.source.Init()
}

func (d DataSourceValue) GetHostDir() (string, error) {
	return d.source.GetHostDir()
}

func (d DataSourceValue) String() string {
	return d.source.Type()
}

func (d DataSourceValue) Type() string {
	return d.source.Type()
}

func (d DataSourceValue) Freeze() {}

func (d DataSourceValue) Truth() starlark.Bool { return true }

func (d DataSourceValue) Hash() (uint32, error) {
	return d.source.Hash()
}

func NewDataSourceValue(source envddata.DataSource) *DataSourceValue {
	return &DataSourceValue{source: source}
}
