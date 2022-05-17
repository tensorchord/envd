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

package main

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tensorchord/envd/pkg/docker"
)

var _ = Describe("up command", func() {
	buildContext := "testdata"
	args := []string{
		"envd.test", "--debug", "up", "--path", buildContext, "--detach",
	}
	BeforeEach(func() {
		cli, err := docker.NewClient(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		_ = cli.Destroy(context.TODO(), buildContext)
	})
	When("given the right arguments", func() {
		It("should up and destroy successfully", func() {
			_, err := run(args)
			Expect(err).NotTo(HaveOccurred())
			destroyArgs := []string{
				"envd.test", "--debug", "destroy", "--path", buildContext,
			}
			_, err = run(destroyArgs)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
