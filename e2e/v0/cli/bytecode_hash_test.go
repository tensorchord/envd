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

package cli

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v0"
)

func appendSomeToFile(path string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	blank := "\n\n"
	_, err = f.Write([]byte(blank))
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

var _ = Describe("bytecode hash cache target", func() {
	exampleName := "quick-start"
	It("add some blank to build.envd", func() {
		testcase := "add-blank"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		ctx := context.TODO()
		e.BuildImage(false)()
		engine := e2e.GetEngine(ctx)
		imageSum, err := engine.GetImage(ctx, e.Tag)
		Expect(err).NotTo(HaveOccurred())
		oldCreated := imageSum.Created
		appendSomeToFile("testdata/" + exampleName + "/build.envd")
		e.BuildImage(false)()
		imageSum, err = engine.GetImage(ctx, e.Tag)
		Expect(err).NotTo(HaveOccurred())
		newCreated := imageSum.Created
		Expect(oldCreated).To(Equal(newCreated))
		e.RemoveImage()()
	})
})
