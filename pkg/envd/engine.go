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
	"context"
	"time"

	dockertypes "github.com/docker/docker/api/types"

	"github.com/tensorchord/envd/pkg/lang/ir"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/types"
)

// Engine is the core engine to manage the envd environments.
type Engine interface {
	SSHClient
	ImageClient
	EnvironmentClient
	VersionClient
}

type EnvironmentClient interface {
	PauseEnvironment(ctx context.Context, env string) (string, error)
	ResumeEnvironment(ctx context.Context, env string) (string, error)
	GetEnvironment(ctx context.Context, env string) (*types.EnvdEnvironment, error)
	ListEnvironment(ctx context.Context) ([]types.EnvdEnvironment, error)
	ListEnvRuntimeGraph(ctx context.Context, env string) (*ir.RuntimeGraph, error)
	ListEnvDependency(ctx context.Context, env string) (*types.Dependency, error)
	ListEnvPortBinding(ctx context.Context, env string) ([]types.PortBinding, error)

	CleanEnvdIfExists(ctx context.Context, name string, force bool) error
	// StartEnvd creates the container for the given tag and container name.
	StartEnvd(ctx context.Context, so StartOptions) (*StartResult, error)

	IsRunning(ctx context.Context, name string) (bool, error)
	Exists(ctx context.Context, name string) (bool, error)
	WaitUntilRunning(ctx context.Context, name string, timeout time.Duration) error
	Destroy(ctx context.Context, name string) (string, error)
}

type SSHClient interface {
	GenerateSSHConfig(name, iface, privateKeyPath string,
		startResult *StartResult) (sshconfig.EntryOptions, error)
	Attach(name, iface, privateKeyPath string,
		startResult *StartResult, g ir.Graph) error
}

type ImageClient interface {
	ListImage(ctx context.Context) ([]types.EnvdImage, error)
	ListImageDependency(ctx context.Context, image string) (*types.Dependency, error)
	GetImage(ctx context.Context, image string) (types.EnvdImage, error)
	PruneImage(ctx context.Context) (dockertypes.ImagesPruneReport, error)
}

type VersionClient interface {
	GetInfo(ctx context.Context) (*types.EnvdInfo, error)
	GPUEnabled(ctx context.Context) (bool, error)
}
