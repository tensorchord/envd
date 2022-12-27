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

package json

import (
	"fmt"

	"github.com/tensorchord/envd/pkg/app/formatter"
	"github.com/tensorchord/envd/pkg/types"
)

type envJSONLs struct {
	Name      string `json:"name"`
	Endpoint  string `json:"endpoint,omitempty"`
	SSHTarget string `json:"ssh_target"`
	Image     string `json:"image"`
	GPU       bool   `json:"gpu"`
	CUDA      string `json:"cuda,omitempty"`
	CUDNN     string `json:"cudnn,omitempty"`
	Status    string `json:"status"`
}

func PrintEnvironments(envs []types.EnvdEnvironment) error {
	output := []envJSONLs{}
	for _, env := range envs {
		item := envJSONLs{
			Name:      env.Name,
			Endpoint:  formatter.FormatEndpoint(env),
			SSHTarget: fmt.Sprintf("%s.envd", env.Name),
			Image:     env.Spec.Image,
			GPU:       env.GPU,
			CUDA:      env.CUDA,
			CUDNN:     env.CUDNN,
			Status:    env.Status.Phase,
		}
		output = append(output, item)
	}
	return printJSON(output)
}

type envJSONDescribe struct {
	Ports        []envJSONPort       `json:"ports,omitempty"`
	Dependencies []envJSONDependency `json:"dependencies,omitempty"`
}
type envJSONPort struct {
	Name          string `json:"name"`
	ContainerPort string `json:"container_port"`
	Protocol      string `json:"protocol"`
	HostIP        string `json:"host_ip"`
	HostPort      string `json:"host_port"`
}

type envJSONDependency struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func PrintEnvironmentDescriptions(dep *types.Dependency, ports []types.PortBinding) error {
	output := envJSONDescribe{}

	for _, port := range ports {
		port := envJSONPort{
			Name:          port.Name,
			ContainerPort: port.Port,
			Protocol:      port.Protocol,
			HostIP:        port.HostIP,
			HostPort:      port.HostPort,
		}
		output.Ports = append(output.Ports, port)
	}
	for _, p := range dep.PyPIPackages {
		dependency := envJSONDependency{
			Name: p,
			Type: "Python",
		}
		output.Dependencies = append(output.Dependencies, dependency)
	}
	for _, p := range dep.APTPackages {
		dependency := envJSONDependency{
			Name: p,
			Type: "APT",
		}
		output.Dependencies = append(output.Dependencies, dependency)
	}
	return printJSON(output)
}
