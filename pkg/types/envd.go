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

	"github.com/docker/docker/api/types"
)

type EnvdImage struct {
	types.ImageSummary

	GPU        bool   `json:"gpu,omitempty"`
	CUDA       string `json:"cuda,omitempty"`
	CUDNN      string `json:"cudnn,omitempty"`
	Dependency `json:",inline,omitempty"`
}

type EnvdEnvironment struct {
	types.Container

	// The name of the environment.
	Name        string `json:"name,omitempty"`
	GPU         bool   `json:"gpu,omitempty"`
	CUDA        string `json:"cuda,omitempty"`
	CUDNN       string `json:"cudnn,omitempty"`
	JupyterAddr string `json:"jupyter_addr,omitempty"`
	Dependency  `json:",inline,omitempty"`
}

type Dependency struct {
	APTPackages  []string `json:"apt_packages,omitempty"`
	PyPIPackages []string `json:"pypi_packages,omitempty"`
}

const (
	ContainerLabelName        = "ai.tensorchord.envd.name"
	ContainerLabelJupyterAddr = "ai.tensorchord.envd.jupyter.address"

	ImageLabelVendor = "ai.tensorchord.envd.vendor"
	ImageLabelGPU    = "ai.tensorchord.envd.gpu"
	ImageLabelAPT    = "ai.tensorchord.envd.apt.packages"
	ImageLabelPyPI   = "ai.tensorchord.envd.pypi.packages"
	ImageLabelCUDA   = "ai.tensorchord.envd.gpu.cuda"
	ImageLabelCUDNN  = "ai.tensorchord.envd.gpu.cudnn"

	ImageVendorEnvd = "envd"
)

func NewImage(image types.ImageSummary) (*EnvdImage, error) {
	img := EnvdImage{
		ImageSummary: image,
	}
	if gpuEnabled, ok := image.Labels[ImageLabelGPU]; ok {
		img.GPU = gpuEnabled == "true"
		img.CUDA = image.Labels[ImageLabelCUDA]
		img.CUDNN = image.Labels[ImageLabelCUDNN]
	}
	dep, err := newDependencyFromLabels(image.Labels)
	if err != nil {
		return nil, err
	}
	img.Dependency = *dep
	return &img, nil
}

func NewEnvironment(ctr types.Container) (*EnvdEnvironment, error) {
	env := EnvdEnvironment{
		Container: ctr,
	}
	if gpuEnabled, ok := ctr.Labels[ImageLabelGPU]; ok {
		env.GPU = gpuEnabled == "true"
		env.CUDA = ctr.Labels[ImageLabelCUDA]
		env.CUDNN = ctr.Labels[ImageLabelCUDNN]
	}
	if name, ok := ctr.Labels[ContainerLabelName]; ok {
		env.Name = name
	}
	if jupyterAddr, ok := ctr.Labels[ContainerLabelJupyterAddr]; ok {
		env.JupyterAddr = jupyterAddr
	}

	dep, err := newDependencyFromLabels(ctr.Labels)
	if err != nil {
		return nil, err
	}
	env.Dependency = *dep
	return &env, nil
}

func NewDependencyFromContainerJSON(ctr types.ContainerJSON) (*Dependency, error) {
	return newDependencyFromLabels(ctr.Config.Labels)
}

func NewDependencyFromImage(img types.ImageSummary) (*Dependency, error) {
	return newDependencyFromLabels(img.Labels)
}

func GetImageName(image EnvdImage) string {
	if len(image.ImageSummary.RepoTags) != 0 {
		return image.ImageSummary.RepoTags[0]
	}
	return "<none>"
}

func newDependencyFromLabels(label map[string]string) (*Dependency, error) {
	dep := Dependency{}
	if pkgs, ok := label[ImageLabelAPT]; ok {
		lst, err := parseAPTPackages(pkgs)
		if err != nil {
			return nil, err
		}
		dep.APTPackages = lst
	}
	if pkgs, ok := label[ImageLabelPyPI]; ok {
		lst, err := parsePyPIPackages(pkgs)
		if err != nil {
			return nil, err
		}
		dep.PyPIPackages = lst
	}
	return &dep, nil
}

func parseAPTPackages(lst string) ([]string, error) {
	var pkgs []string
	if err := json.Unmarshal([]byte(lst), &pkgs); err != nil {
		return nil, err
	}
	return pkgs, nil
}

func parsePyPIPackages(lst string) ([]string, error) {
	var pkgs []string
	if err := json.Unmarshal([]byte(lst), &pkgs); err != nil {
		return nil, err
	}
	return pkgs, nil
}
