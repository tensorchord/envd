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
	"bytes"
	"context"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/docker"
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
	app  *app.EnvdApp
}

func NewExample(name string, testcaseAbbr string) *Example {
	tag := name + ":" + testcaseAbbr
	app := app.New()
	return &Example{
		Name: name,
		Tag:  tag,
		app:  &app,
	}
}

func (e *Example) Exec(cmd string) (string, error) {
	args := []string{
		"envd.test", "run", "--name", e.Name, "--raw", cmd,
	}

	buffer := new(bytes.Buffer)
	e.app.Writer = buffer

	err := e.app.Run(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to start `run` command")
	}
	return strings.Trim(buffer.String(), "\n"), nil
}

func (e *Example) ExecRuntimeCommand(cmd string) (string, error) {
	args := []string{
		"envd.test", "run", "--name", e.Name, "--command", cmd,
	}

	buffer := new(bytes.Buffer)
	e.app.Writer = buffer

	err := e.app.Run(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to start `run` command")
	}
	return strings.Trim(buffer.String(), "\n"), nil
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
