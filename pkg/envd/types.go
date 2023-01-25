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

package envd

import (
	"time"

	"github.com/tensorchord/envd-server/api/types"

	"github.com/tensorchord/envd/pkg/lang/ir"
)

const (
	Localhost = "127.0.0.1"
)

var (
	waitingInterval = 1 * time.Second
)

type StartOptions struct {
	Image           string
	EnvironmentName string
	BuildContext    string
	NumGPU          int
	NumCPU          string
	NumMem          string
	Timeout         time.Duration
	ShmSize         int
	Forced          bool
	SshdHost        string

	EngineSource
}

type EngineSource struct {
	DockerSource     *DockerSource
	EnvdServerSource *EnvdServerSource
}

type DockerSource struct {
	Graph        ir.Graph
	MountOptions []string
}

type EnvdServerSource struct{}

type StartResult struct {
	// TODO(gaocegege): Make result a chan, to send running status to the receiver.
	SSHPort int
	Address string
	Name    string

	Ports []types.EnvironmentPort
}
