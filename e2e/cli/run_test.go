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

	When("given the right arguments", func() {
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
		envdApp = app.New()
		// build env image
		buildArgs := append(baseArgs, "build", "--tag", buildImageName, "--path", buildContext)
		err = envdApp.Run(buildArgs)
		Expect(err).NotTo(HaveOccurred())
		// run env
		runArgs := append(baseArgs, "run", "--image", buildImageName, "--name", env, "--detach")
		err = envdApp.Run(runArgs)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterAll(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
		envdApp := app.New()
		e2e.ResetEnvdApp()
		err := envdApp.Run([]string{"envd.test", "--debug", "bootstrap"})
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
})
