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

package e2e

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/ssh"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
)

func (e *Example) BuildImage(force bool) func() {
	return func() {
		logrus.Info("building quick-start image")
		buildContext := "testdata/" + e.Name
		args := []string{
			"envd.test", "--debug", "build", "--path", buildContext, "--tag", e.Tag,
		}
		if force {
			args = append(args, "--force")
		}
		err := e.app.Run(args)
		if err != nil {
			panic(err)
		}
	}
}

func (e *Example) RemoveImage() func() {
	return func() {
		ctx := context.TODO()
		dockerClient, err := docker.NewClient(ctx)
		if err != nil {
			panic(err)
		}
		err = dockerClient.RemoveImage(ctx, e.Tag)
		if err != nil {
			panic(err)
		}
	}
}

func GetDockerClient(ctx context.Context) docker.Client {
	dockerClient, err := docker.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return dockerClient
}

type Example struct {
	Name string
	Tag  string
	app  app.EnvdApp
}

func NewExample(name string, testcaseAbbr string) *Example {
	tag := name + ":" + testcaseAbbr
	return &Example{
		Name: name,
		Tag:  tag,
		app:  app.New(),
	}
}

func (e *Example) Exec(cmd string) string {
	sshClient := e.getSSHClient()
	ret, err := sshClient.ExecWithOutput(cmd)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(ret), "\n")
}

func (e *Example) RunContainer() func() {
	return func() {
		buildContext := "testdata/" + e.Name
		args := []string{
			"envd.test", "--debug", "up", "--path", buildContext, "--tag", e.Tag, "--detach", "--force",
		}
		err := e.app.Run(args)
		if err != nil {
			panic(err)
		}
	}
}

func (e *Example) DestroyContainer() func() {
	return func() {
		buildContext := "testdata/" + e.Name
		args := []string{
			"envd.test", "--debug", "destroy", "--path", buildContext,
		}
		err := e.app.Run(args)
		if err != nil {
			panic(err)
		}
	}
}

func (e *Example) getSSHClient() ssh.Client {
	localhost := "127.0.0.1"
	port, err := sshconfig.GetPort(e.Name)
	if err != nil {
		panic(err)
	}
	priv_path := sshconfig.GetPrivateKey()
	sshClient, err := ssh.NewClient(
		localhost, "envd", port, true, priv_path, "")
	if err != nil {
		panic(err)
	}
	return sshClient
}
