// Copyright 2023 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tensorchord/envd/e2e/v1"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var _ = Describe("e2e moby builder test", Ordered, func() {
	exampleName := "build-test"
	defaultContext := "default"
	mobyContext := "envd-test-moby"
	ctx := types.Context{
		Name:    mobyContext,
		Builder: types.BuilderTypeMoby,
		Runner:  types.RunnerTypeDocker,
	}

	BeforeAll(func() {
		Expect(home.Initialize()).To(Succeed())
		err := home.GetManager().ContextCreate(ctx, true)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should find a new context", func() {
		ctx, err := home.GetManager().ContextGetCurrent()
		Expect(err).NotTo(HaveOccurred())
		Expect(ctx.Name).To(Equal(mobyContext))
		Expect(ctx.Builder).To(Equal(types.BuilderTypeMoby))
	})

	It("build images with moby", func() {
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), "e2e")
		e.BuildImage(true)()
	})

	AfterAll(func() {
		err := home.GetManager().ContextUse(defaultContext)
		Expect(err).NotTo(HaveOccurred())
		err = home.GetManager().ContextRemove(mobyContext)
		Expect(err).NotTo(HaveOccurred())
		context, err := home.GetManager().ContextGetCurrent()
		Expect(err).NotTo(HaveOccurred())
		Expect(context.Name).To(Equal(defaultContext))
	})
})
