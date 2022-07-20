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

func BuildImage(exampleName string) func() {
	return func() {
		logrus.Info("building quick-start image")
		err := BuildExampleImage(exampleName, app.New())
		if err != nil {
			panic(err)
		}
	}
}

func RemoveImage(exampleName string) func() {
	return func() {
		err := RemoveExampleImage(exampleName)
		if err != nil {
			panic(err)
		}
	}
}

func RunContainer(exampleName string) func() {
	return func() {
		err := RunExampleContainer(exampleName, app.New())
		if err != nil {
			panic(err)
		}
	}
}

func DestoryContainer(exampleName string) func() {
	return func() {
		err := DestroyExampleContainer(exampleName, app.New())
		if err != nil {
			panic(err)
		}
	}
}

type Example struct {
	Name string
}

func example(name string) *Example {
	return &Example{
		Name: name,
	}
}

func (e *Example) Exec(cmd string) string {
	sshClient := getSSHClient(e.Name)
	ret, err := sshClient.ExecWithOutput(cmd)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(ret), "\n")
}

func BuildExampleImage(exampleName string, app app.EnvdApp) error {
	buildContext := "testdata/" + exampleName
	tag := exampleName + ":e2etest"
	args := []string{
		"envd.test", "--debug", "build", "--path", buildContext, "--tag", tag, "--force",
	}
	err := app.Run(args)
	return err
}

func RemoveExampleImage(exampleName string) error {
	ctx := context.TODO()
	dockerClient, err := docker.NewClient(ctx)
	if err != nil {
		return err
	}
	err = dockerClient.RemoveImage(ctx, exampleName+":e2etest")
	if err != nil {
		return err
	}
	return nil
}

func RunExampleContainer(exampleName string, app app.EnvdApp) error {
	buildContext := "testdata/" + exampleName
	tag := exampleName + ":e2etest"
	args := []string{
		"envd.test", "--debug", "up", "--path", buildContext, "--tag", tag, "--detach", "--force",
	}
	err := app.Run(args)
	return err
}

func DestroyExampleContainer(exampleName string, app app.EnvdApp) error {
	buildContext := "testdata/" + exampleName
	args := []string{
		"envd.test", "--debug", "destroy", "--path", buildContext,
	}
	err := app.Run(args)
	return err
}

func getSSHClient(exampleName string) ssh.Client {
	localhost := "127.0.0.1"
	port, err := sshconfig.GetPort(exampleName)
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
