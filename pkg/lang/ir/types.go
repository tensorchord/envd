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

package ir

import (
	"github.com/opencontainers/go-digest"
)

// The results during runtime should be maintained here
type RuntimeGraph struct {
	RuntimeCommands   map[string]string `json:"commands,omitempty"`
	RuntimeDaemon     [][]string        `json:"daemon,omitempty"`
	RuntimeInitScript [][]string        `json:"init_script,omitempty"`
	RuntimeEnviron    map[string]string `json:"environ,omitempty"`
	RuntimeEnvPaths   []string          `json:"env_paths,omitempty"`
	RuntimeExpose     []ExposeItem      `json:"expose,omitempty"`
}

type CopyInfo struct {
	Source      string
	Destination string
}

type MountInfo struct {
	Source      string
	Destination string
}

type HTTPInfo struct {
	URL      string
	Checksum digest.Digest
	Filename string
}

type RStudioServerConfig struct {
}

type Language struct {
	Name    string
	Version *string
}

type CondaConfig struct {
	CondaPackages      []string
	AdditionalChannels []string
	CondaChannel       *string
	CondaEnvFileName   string
	UseMicroMamba      bool
}

type GitConfig struct {
	Name   string
	Email  string
	Editor string
}

type ExposeItem struct {
	EnvdPort      int
	HostPort      int
	ServiceName   string
	ListeningAddr string
}

type JupyterConfig struct {
	Token string
	Port  int64
}

type RunBuildCommand struct {
	Commands  []string
	MountHost bool
}

type APTConfig struct {
	Name       string
	Enabled    string
	Types      string
	URIs       string
	Suites     string
	Components string
	Signed     string
	Arch       string
}
