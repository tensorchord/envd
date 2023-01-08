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
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/driver/docker"
	"github.com/tensorchord/envd/pkg/envd"
	ir "github.com/tensorchord/envd/pkg/lang/ir/v0"
	"github.com/tensorchord/envd/pkg/types"
)

func BuildContextDirWithName(name string) string {
	return filepath.Join("testdata", name)
}

func ResetEnvdApp() {
	ir.DefaultGraph = ir.NewGraph()
}

func (e *Example) BuildImage(force bool) func() {
	return func() {
		logrus.Infof("building %s image in %s", e.Name, e.BuildContextPath)
		args := []string{
			"envd.test", "--debug", "build",
			"--path", e.BuildContextPath, "--tag", e.Tag,
		}
		if force {
			args = append(args, "--force")
		}
		ResetEnvdApp()
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

func GetEngine(ctx context.Context) envd.Engine {
	opt := envd.Options{
		Context: &types.Context{
			Runner: types.RunnerTypeDocker,
		},
	}
	engine, err := envd.New(ctx, opt)
	if err != nil {
		panic(err)
	}
	return engine
}

type Example struct {
	Tag              string
	BuildContextPath string
	// Name is the filepath.Base(BuildContextPath).
	Name string

	app *app.EnvdApp
}

func NewExample(path string, testcaseAbbr string) *Example {
	name := filepath.Base(path)
	tag := name + ":" + testcaseAbbr
	app := app.New()
	return &Example{
		Tag:              tag,
		Name:             name,
		BuildContextPath: path,
		app:              &app,
	}
}

func (e *Example) Exec(cmd string) (string, error) {
	args := []string{
		"envd.test", "exec", "--name", e.Name, "--raw", cmd,
	}

	buffer := new(bytes.Buffer)
	e.app.Writer = buffer

	ResetEnvdApp()
	err := e.app.Run(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to start `run` command")
	}
	return strings.Trim(buffer.String(), "\n"), nil
}

func (e *Example) ExecRuntimeCommand(cmd string) (string, error) {
	buildContext := e.BuildContextPath
	args := []string{
		"envd.test", "--debug", "exec", "-p", buildContext, "--command", cmd,
	}

	buffer := new(bytes.Buffer)
	e.app.Writer = buffer

	ResetEnvdApp()
	err := e.app.Run(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to start `run` command")
	}
	return strings.Trim(buffer.String(), "\n"), nil
}

func (e *Example) RunContainer() func() {
	return func() {
		buildContext := e.BuildContextPath
		args := []string{
			"envd.test", "--debug", "up", "--path", buildContext, "--tag", e.Tag, "--detach", "--force",
		}
		ResetEnvdApp()
		err := e.app.Run(args)
		if err != nil {
			panic(err)
		}
	}
}

func (e *Example) DestroyContainer() func() {
	return func() {
		buildContext := e.BuildContextPath
		args := []string{
			"envd.test", "--debug", "destroy", "--path", buildContext,
		}
		err := e.app.Run(args)
		if err != nil {
			panic(err)
		}
	}
}
