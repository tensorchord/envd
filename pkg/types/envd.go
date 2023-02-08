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

package types

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/api/types"
	"github.com/moby/buildkit/util/system"
	servertypes "github.com/tensorchord/envd-server/api/types"

	"github.com/tensorchord/envd/pkg/util/netutil"
	"github.com/tensorchord/envd/pkg/version"
)

const (
	// environment PATH
	DefaultSystemPath = system.DefaultPathEnvUnix
	DefaultCondaPath  = "/opt/conda/envs/envd/bin:/opt/conda/bin:/home/envd/.local/bin"
	DefaultJuliaPath  = "/usr/local/julia/bin"
	// image
	PythonBaseImage   = "ubuntu:20.04"
	EnvdStarshipImage = "tensorchord/starship:v0.0.1"
	// supervisor
	HorustImage      = "tensorchord/horust:v0.2.1"
	HorustServiceDir = "/etc/horust/services"
	HorustLogDir     = "/var/log/horust"
	// env
	EnvdWorkDir = "ENVD_WORKDIR"
)

var EnvdSshdImage = fmt.Sprintf(
	"tensorchord/envd-sshd-from-scratch:%s",
	version.GetVersionForImageTag())

var BaseEnvironment = []struct {
	Name  string
	Value string
}{
	{"DEBIAN_FRONTEND", "noninteractive"},
	{"PATH", DefaultSystemPath},
	{"LANG", "en_US.UTF-8"},
	{"LC_ALL", "en_US.UTF-8"},
}
var BaseAptPackage = []string{
	"bash-static",
	"libtinfo5",
	"libncursesw5",
	// conda dependencies
	"bzip2",
	"ca-certificates",
	"libglib2.0-0",
	"libsm6",
	"libxext6",
	"libxrender1",
	"mercurial",
	"procps",
	"subversion",
	"wget",
	// envd dependencies
	"curl",
	"openssh-client",
	"git",
	"sudo",
	"vim",
	"make",
	"zsh",
	"locales",
}

type EnvdImage struct {
	servertypes.ImageMeta `json:",inline,omitempty"`

	EnvdManifest `json:",inline,omitempty"`
}

type EnvdEnvironment struct {
	servertypes.Environment `json:",inline,omitempty"`

	EnvdManifest `json:",inline,omitempty"`
}

type EnvdManifest struct {
	GPU          bool   `json:"gpu,omitempty"`
	CUDA         string `json:"cuda,omitempty"`
	CUDNN        string `json:"cudnn,omitempty"`
	BuildContext string `json:"build_context,omitempty"`
	Dependency   `json:",inline,omitempty"`
}

type EnvdInfo struct {
	types.Info
}

type EnvdContext struct {
	Current  string    `json:"current,omitempty"`
	Contexts []Context `json:"contexts,omitempty"`
}

type Context struct {
	Name           string      `json:"name,omitempty"`
	Builder        BuilderType `json:"builder,omitempty"`
	BuilderAddress string      `json:"builder_address,omitempty"`
	Runner         RunnerType  `json:"runner,omitempty"`
	RunnerAddress  *string     `json:"runner_address,omitempty"`
}

type BuilderType string

const (
	BuilderTypeMoby             BuilderType = "moby-worker"
	BuilderTypeDocker           BuilderType = "docker-container"
	BuilderTypeNerdctl          BuilderType = "nerdctl-container"
	BuilderTypeKubernetes       BuilderType = "kube-pod"
	BuilderTypeTCP              BuilderType = "tcp"
	BuilderTypeUNIXDomainSocket BuilderType = "unix"
)

type RunnerType string

const (
	RunnerTypeDocker     RunnerType = "docker"
	RunnerTypeEnvdServer RunnerType = "envd-server"
)

type Dependency struct {
	APTPackages  []string `json:"apt_packages,omitempty"`
	PyPIPackages []string `json:"pypi_packages,omitempty"`
}

type RepoInfo struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

type OwnerInfo struct {
	Uid int64 `json:"uid,omitempty"`
	Gid int64 `json:"gid,omitempty"`
}

type PortBinding struct {
	Name     string
	Port     string
	Protocol string
	HostIP   string
	HostPort string
}

type EnvdAuth struct {
	Current string       `json:"current,omitempty"`
	Auth    []AuthConfig `json:"auth,omitempty"`
}

type AuthConfig struct {
	Name     string `json:"name,omitempty"`
	JWTToken string `json:"jwt_token,omitempty"`
}

func NewImageFromSummary(image types.ImageSummary) (*EnvdImage, error) {
	img := EnvdImage{
		ImageMeta: servertypes.ImageMeta{
			Digest:  image.ID,
			Created: image.Created,
			Size:    image.Size,
			Labels:  image.Labels,
		},
	}
	if len(image.RepoTags) > 0 {
		img.Name = image.RepoTags[0]
	}
	m, err := newManifest(image.Labels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse manifest")
	}
	img.EnvdManifest = m
	return &img, nil
}

func NewImageFromMeta(meta servertypes.ImageMeta) (*EnvdImage, error) {
	img := EnvdImage{
		ImageMeta: meta,
	}
	manifest, err := newManifest(img.Labels)
	if err != nil {
		return nil, err
	}
	img.EnvdManifest = manifest
	return &img, nil
}

func NewEnvironmentFromContainer(ctr types.Container) (*EnvdEnvironment, error) {
	env := EnvdEnvironment{
		Environment: servertypes.Environment{
			Spec:   servertypes.EnvironmentSpec{Image: ctr.Image},
			Status: servertypes.EnvironmentStatus{Phase: ctr.Status},
		},
	}
	if name, ok := ctr.Labels[ContainerLabelName]; ok {
		env.Name = name
	}
	if jupyterAddr, ok := ctr.Labels[ContainerLabelJupyterAddr]; ok {
		env.Status.JupyterAddr = &jupyterAddr
	}
	if rstudioServerAddr, ok := ctr.Labels[ContainerLabelRStudioServerAddr]; ok {
		env.Status.RStudioServerAddr = &rstudioServerAddr
	}

	m, err := newManifest(ctr.Labels)
	if err != nil {
		return nil, err
	}
	env.EnvdManifest = m
	return &env, nil
}

func NewEnvironmentFromServer(ctr servertypes.Environment) (*EnvdEnvironment, error) {
	env := EnvdEnvironment{
		Environment: ctr,
	}

	m, err := newManifest(ctr.Labels)
	if err != nil {
		return nil, err
	}
	env.EnvdManifest = m
	return &env, nil
}

func newManifest(labels map[string]string) (EnvdManifest, error) {
	manifest := EnvdManifest{}
	if gpuEnabled, ok := labels[ImageLabelGPU]; ok {
		manifest.GPU = gpuEnabled == "true"
	}
	if cuda, ok := labels[ImageLabelCUDA]; ok {
		manifest.CUDA = cuda
	}
	if cudnn, ok := labels[ImageLabelCUDNN]; ok {
		manifest.CUDNN = cudnn
	}
	if context, ok := labels[ImageLabelContext]; ok {
		manifest.BuildContext = context
	}
	dep, err := NewDependencyFromLabels(labels)
	if err != nil {
		return manifest, err
	}
	manifest.Dependency = *dep
	return manifest, nil
}

func NewDependencyFromContainerJSON(ctr types.ContainerJSON) (*Dependency, error) {
	return NewDependencyFromLabels(ctr.Config.Labels)
}

func NewDependencyFromImageSummary(img types.ImageSummary) (*Dependency, error) {
	return NewDependencyFromLabels(img.Labels)
}

func NewPortBindingFromContainerJSON(ctr types.ContainerJSON) ([]PortBinding, error) {
	var labels []servertypes.EnvironmentPort
	err := json.Unmarshal([]byte(ctr.Config.Labels[ImageLabelPorts]), &labels)
	if err != nil {
		return nil, err
	}

	var portMap = make(map[string]string)
	for _, label := range labels {
		portMap[fmt.Sprint(label.Port)] = label.Name
	}

	config := ctr.HostConfig.PortBindings
	var ports []PortBinding
	for port, bindings := range config {
		if len(bindings) <= 0 {
			continue
		}
		binding := bindings[len(bindings)-1]
		ports = append(ports, PortBinding{
			Name:     portMap[port.Port()],
			Port:     port.Port(),
			Protocol: port.Proto(),
			HostIP:   binding.HostIP,
			HostPort: binding.HostPort,
		})
	}

	return ports, nil
}

func NewDependencyFromLabels(label map[string]string) (*Dependency, error) {
	dep := Dependency{}
	if pkgs, ok := label[ImageLabelAPT]; ok {
		lst, err := parseAPTPackages(pkgs)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse apt packages")
		}
		dep.APTPackages = lst
	}
	if pypiCommands, ok := label[ImageLabelPyPI]; ok {
		pkgs, err := parsePyPICommands(pypiCommands)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse pypi commands")
		}
		dep.PyPIPackages = pkgs
	}
	return &dep, nil
}

func parseAPTPackages(lst string) ([]string, error) {
	var pkgs []string
	err := json.Unmarshal([]byte(lst), &pkgs)
	return pkgs, err
}

func parsePyPICommands(lst string) ([]string, error) {
	var pkgs []string
	err := json.Unmarshal([]byte(lst), &pkgs)
	return pkgs, err
}

func (c Context) GetSSHHostname(sshdHost string) (string, error) {
	if c.RunnerAddress == nil {
		return sshdHost, nil
	}

	// TODO(gaocegege): Check ENVD_SERVER_HOST.
	hostname, err := netutil.GetHost(*c.RunnerAddress)
	if err != nil {
		return "", err
	}
	return hostname, nil
}
