package cli

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tensorchord/envd/e2e"
	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/driver/docker"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var _ = Describe("run command", func() {
	buildContext := "testdata/run-test"
	buildImageName := "testdata/run-test:dev"
	env := "run-test"
	baseArgs := []string{
		"envd.test", "--debug",
	}

	BeforeAll(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
		envdApp := app.New()
		err := envdApp.Run(append(baseArgs, "bootstrap"))
		Expect(err).NotTo(HaveOccurred())
		_, err = docker.NewClient(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		c := types.Context{Runner: types.RunnerTypeDocker}
		opt := envd.Options{Context: &c}
		envdEngine, err := envd.New(context.TODO(), opt)
		Expect(err).NotTo(HaveOccurred())
		_, err = envdEngine.Destroy(context.TODO(), env)
		Expect(err).NotTo(HaveOccurred())

	})
	When("given the right arguments", func() {
		It("should build, run and destroy successfully", func() {
			// build env image
			buildArgs := append(baseArgs, "build", "--tag", buildImageName, "--path", buildContext)
			e2e.ResetEnvdApp()
			envdApp := app.New()
			err := envdApp.Run(buildArgs)
			Expect(err).NotTo(HaveOccurred())
			// run env
			runArgs := append(baseArgs, "run", "--image", buildImageName, "--name", env, "--detach")
			err = envdApp.Run(runArgs)
			Expect(err).NotTo(HaveOccurred())
			// destroy env
			destroyArgs := append(baseArgs, "destroy", "--name", env)
			err = envdApp.Run(destroyArgs)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
