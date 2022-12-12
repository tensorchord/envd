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

package v0

import "github.com/tensorchord/envd/pkg/util/fileutil"

const (
	osDefault              = "ubuntu20.04"
	languageDefault        = "python"
	languageVersionDefault = "3"
	CUDNNVersionDefault    = "8"

	aptSourceFilePath = "/etc/apt/sources.list"
	pypiIndexFilePath = "/etc/pip.conf"

	pypiConfigTemplate = `
[global]
index-url=%s
%s

[install]
src = /tmp
`
)

var (
	// used inside the container
	defaultConfigDir   = fileutil.EnvdHomeDir(".config")
	starshipConfigPath = fileutil.EnvdHomeDir(".config", "starship.toml")
)
