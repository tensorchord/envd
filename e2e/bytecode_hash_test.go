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
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func appendSomeToFile(path string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	blank := "\n\n"
	f.Write([]byte(blank))
	if err := f.Close(); err != nil {
		panic(err)
	}
}

var _ = Describe("bytecode hash cache target", func() {
	exampleName := "quick-start"
	It("add some blank to build.envd", func() {
		ctx := context.TODO()
		BuildImage(exampleName, false)()
		dockerClient := GetDockerClient(ctx)
		imageSum, err := dockerClient.GetImage(ctx, exampleName+":e2etest")
		Expect(err).NotTo(HaveOccurred())
		oldCreated := imageSum.Created
		appendSomeToFile("testdata/" + exampleName + "/build.envd")
		BuildImage(exampleName, false)()
		imageSum, err = dockerClient.GetImage(ctx, exampleName+":e2etest")
		Expect(err).NotTo(HaveOccurred())
		newCreated := imageSum.Created
		Expect(oldCreated).To(Equal(newCreated))
		RemoveExampleImage(exampleName)
	})
})
