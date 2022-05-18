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

package compileui

import (
	"time"

	"github.com/tensorchord/envd/pkg/editor/vscode"
)

type Action string

const (
	ActionStart Action = "start"
	ActionEnd   Action = "end"
)

type Result struct {
	plugins []*PluginInfo
	ZSHInfo *ZSHInfo
}

type PluginInfo struct {
	vscode.Plugin
	startTime *time.Time
	endTime   *time.Time
	cached    bool
}

type ZSHInfo struct {
	OHMYZSH   string
	startTime *time.Time
	endTime   *time.Time
	cached    bool
}
