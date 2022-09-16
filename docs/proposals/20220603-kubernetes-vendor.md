# Multi-runtime Support

Authors:
- [Ce Gao](https://github.com/gaocegege)

## Summary

envd builds and runs the environments with Docker. This proposal is to support Kubernetes as a runtime option. After this feature, envd can build the environment with Docker and run on Kubernetes.

## Motivation

Docker runtime only works for individual developers. There are many data scientists developing models on Kubernetes. In the current design, they may need to run `envd build`, then push the image to a OCI image registry manually. The image can be used after a deployment or a pod is created in Kubernetes.

We can provide native support for ease of use.

## Goals

- Build the envd environments with the buildkitd pod on Kubernetes
- Run environments on Kubernetes
  - Compose Kubernetes resources (e.g. deployments, services, pods) when run `envd up`.
  - Push the image to a registry.
  - Port forward for sshd and jupyter notebook service.
  - Sync the pod file system with the local host file system.
- Manage environments on Kubernetes
  - List environments running on Kubernetes, add a new field `vendor` when run `envd get envs`
  - List dependencies of the environment on Kubernetes

## Non-Goals

- Monitoring on Kubernetes
- Pause and unpause environments on Kubernetes (Kubernetes does not support the primitives)
- Manage envd images on Kubernetes

## Proposal

### User stories

#### CI/CD with envd

As a ML infra engineer, I want to build the images with CI/CD tools. Docker may not be set up in the CI/CD runners, so that envd should support to build with buildkitd running on the Kubernetes.

#### Remote development

As a AI/ML engineer, I want to develop on the remote cluster (managed by Kubernetes).

### CLI behavior

We need to support primitives like `run`:

```
$ envd run --env my-env --image my-image
```

User may use `envd` to build the image, and use it on Kubernetes. Thus they need to push the image to a registry. `envd up` may not satisfy the requirements.

The end-to-end process will be:

```
$ envd context create --name test --builder-name test --use --builder kube-pod --runner kubernetes --runner-config kube-pod-config.yaml
$ envd build
$ envd push
$ envd run --env test --image test

or 

$ envd context create --name test --builder-name test --use --builder kube-pod --runner kubernetes --runner-config kube-pod-config.yaml
$ envd up
```

### Implementation details/Notes/Constraints

#### envd API server



#### envd contexts

envd uses buildkitd container in the local docker host to build the images by default. `context` command will be introduced to support remote build.

```
NAME:
   envd context - Manage envd contexts

USAGE:
   envd context command [command options] [arguments...]

COMMANDS:
   create   Create envd context
   ls       List envd contexts
   rm       Remove envd context
   use      Use the specified envd context
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)
```

The structure of the envd context is:

```go
type Context struct {
	Name          string      `json:"name,omitempty"`
	Builder       BuilderType `json:"builder,omitempty"`
	BuilderSocket string      `json:"builder_socket,omitempty"`

	Runner     RunnerType `json:"runner,omitempty"`
	DockerHost *string    `json:"docker_host,omitempty"`
	KubeConfig *string    `json:"kube_config,omitempty"`
}

type BuilderType string

const (
	BuilderTypeDocker     BuilderType = "docker-container"
	BuilderTypeKubernetes BuilderType = "kube-pod"
)

type RunnerType string

const (
	RunnerTypeDocker     RunnerType = "docker"
	RunnerTypeKubernetes RunnerType = "kubernetes"
)
```

The default envd context is:

```go
Context{
    Name:          "default",
    Builder:       BuilderTypeDocker,
    BuilderSocket: "unix:///var/run/docker.sock",
    Runner:        RunnerTypeDocker,
    DockerHost:    nil,
    KubeConfig:    nil,
}
```

Users may use `envd context create` to create new context and use it:

```
$ envd context create --name my-context --builder-type docker-container --runner-type docker --docker-host unix:///var/run/docker.sock
```

#### Entrypoint

Currently, entrypoint is set during the container creation:

```go
	config := &container.Config{
		Entrypoint: []string{
			"tini",
			"--",
			"bash",
			"-c",
		},
	}
	config.Entrypoint = append(config.Entrypoint,
		entrypointSH(g, config.WorkingDir, sshPort))
```

To support the Kubernetes, envd should generate entrypoint and cmd for the image during the build process, instead of in the runtime.

#### Bi-directional filesystem sync

A per-pod sidecar is introduced to sync files between host and Kubernetes pod. [syncthing](https://docs.syncthing.net/intro/getting-started.html) will run in the per-pod sidecar.

The sidecar design does not involve changing the image.

#### Port forwarding

There are several ports will be used in envd:

- sshd server port. envd randomly selects a host port for sshd.
- jupyter notebook port.
- RStudio server port (in the future).

[portforwarder](https://github.com/kubernetes/kubernetes/blob/v1.6.0-alpha.0/pkg/kubectl/cmd/portforward.go#L107) in client-go will be used, to forward pod ports to the host.

**Pending Issues:**

- `ssh envd` will be broken on Kubernetes if we do not have a daemon process to forward the ports.

#### Runtime interface

A new interface `Runtime` will be introduced. And the interface `docker.Client` will be updated.

```go
type Runtime interface {
	StartEnvironment()

    ListEnvironments()
    GetEnvironment()
    PauseEnvironment()
    ResumeEnvironment()
    DestroyEnvironment()

    GPUEnabled()
}
```

`docker.Client` is used to:

- Destroy the environment via `envd destroy`
- Check if GPU is supported and run the docker container via `envd up`
- Load the image into the local docker host via `builder`
- Check if the buildkitd container is running in `buildkit.Client`
- List images, environments in `envd.Engine`

`envd destroy`, `envd up`, and `envd.Engine` should not use `docker.Client` any more. It should be migrated to `Runtime`. Because these func calls are related to runtime.

`builder` and `buildkit.Client` still needs to use `docker.Client` since we keep using docker to build the images.

The pseudocode of the new logic will be like:

```go
// Destroy the environment
func destroy(clicontext *cli.Context) error {
    runtime := runtime.New(getCLIFlag("runtime-vendor"))
    runtime.DestroyEnvironment(...)
}

func up(clicontext *cli.Context) error {
    // Build the image first
    builder.Build()
    // Push the image to a registry
    pushImage()
    runtime := runtime.New(getCLIFlag("runtime-vendor"))
    runtime.StartEnvironment(...)
}
```

## Design Details

### Test Plan

[kind](https://github.com/kubernetes-sigs/kind)-based integration test should be added. 
