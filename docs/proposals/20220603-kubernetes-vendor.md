# Multi-runtime Support

Authors:
- [Ce Gao](https://github.com/gaocegege)

## Summary

envd builds and runs the environments with Docker. This proposal is to support Kubernetes as a runtime option. After this feature, envd can build the environment with Docker and run on Kubernetes.

## Motivation

Docker runtime only works for individual developers. There are many data scientists developing models on Kubernetes. In the current design, they may need to run `envd build`, then push the image to a OCI image registry manually. The image can be used after a deployment or a pod is created in Kubernetes.

We can provide native support for ease of use.

## Goals

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
- [Support context in envd CLI to switch runtimes between Docker and Kubernetes](https://github.com/tensorchord/envd/issues/92)
- Pause and unpause environments on Kubernetes (Kubernetes does not support the primitives)
- Build the envd environments with the buildkitd pod on Kubernetes
- Manage envd images on Kubernetes

## Proposal

### User Stories

### Implementation Details/Notes/Constraints

#### Runtime Interface

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
    runtime := runtime.New(getCLIFlag("runtime-vendor"))
    runtime.StartEnvironment(...)
}
```

#### `StartEnvironment` Implementation



## Design Details

### Test Plan

### Graduation Criteria
