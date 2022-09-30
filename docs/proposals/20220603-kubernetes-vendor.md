# Multi-runtime Support

Authors:
- [Ce Gao](https://github.com/gaocegege)

## Summary

envd builds and runs the environments with Docker. This proposal is to support Kubernetes as a runtime option. After this feature, envd can build the environment with Docker and run on Kubernetes.

## Motivation

Docker runtime only works for individual developers. There are many data scientists developing models on Kubernetes. In the current design, they may need to run `envd build` and then manually push the image to an OCI image registry. The image can be used after a deployment or a pod is created in Kubernetes.

We can provide native support for ease of use.

## Goals

- Run environments on Kubernetes
  - Build the envd environments with the buildkitd pod
  - Compose Kubernetes resources (e.g. services, pods) when run `envd up`
  - Push the image to a registry
  - Port forward for sshd and jupyter notebook service
- Manage environments on Kubernetes
  - List environments running on Kubernetes
- Garbage collection
  - Delete idle pods

## Non-Goals

- Sync the pod file system with the local host file system.
  - It should be discussed in [tensorchord/envd#928](https://github.com/tensorchord/envd/discussions/928)
- Monitoring on Kubernetes
- Pause and unpause environments on Kubernetes (Kubernetes does not support the primitives)
- Manage envd images on Kubernetes

## Terms

- buildkit: https://github.com/moby/buildkit is the library behind `docker build`.
- containerssh: https://containerssh.io/ is a ssh proxy server.

## Proposal

### User stories

#### CI/CD with envd

As an ML infra engineer, I want to build the images with CI/CD tools. Docker may not be set up in the CI/CD runners so envd should support building with buildkitd running on the Kubernetes.

#### Remote development

As an AI/ML engineer, I want to develop on the remote cluster (managed by Kubernetes).

### CLI behavior

We need to support primitives like `run`:

```
$ envd run --env my-env --image my-image
```

User may use `envd` to build the image, and use it on Kubernetes. Thus they need to push the image to a registry. `envd up` may not satisfy the requirements.

The end-to-end process will be:

```
$ envd context create --name test --builder-name test --use --builder kube-pod --runner server --runner-addr http://localhost:2222
$ envd login
$ envd build
$ envd push
$ envd run --env test --image test

or 

$ envd context create --name test --builder-name test --use --builder kube-pod --runner server --runner-addr http://localhost:2222
$ envd login
$ envd up
```

The user first declares which runner will be used. `server` should be used if the user expects running on Kubernetes. The `server` is the short name of envd API server.

Then the user needs to login to the server with the public key.

After that, the user can manage the environments on Kubernetes.

### Design Details

#### envd API server

![](https://user-images.githubusercontent.com/5100735/190376297-f2ce5938-ac29-4eaf-8253-61461bcef8bd.png)

envd API server provides RESTful HTTP API to the envd CLI, and acts as the ssh proxy server (TCP), with the help of [libcontainerssh](https://github.com/ContainerSSH/libcontainerssh)

envd CLI communicates with the API server to:

- Register public key to the server for auth.
- Create backend environment to connect. If it is on Kubernetes, the backend environment will be run inside a pod.
- Connect to the ssh proxy server.

#### Reconnection

If the user reconnects to the pod, envd API server needs to know the correct pod to reconnect.

Thus the backend pod should have a unique identifier. When the user run envd attach (Or maybe envd ssh), the only information given to envd API server, is the username. Thus the CLI should encode the project name, the user, and other info, and use it as the username.

Or, the server grants the user a unique random username when the envd CLI logins. The envd CLI uses the unique username to communicate with the API server.

The main challenge here is to identify the correct pod with the existing SSH protocol.

#### Garbage collection

The pod can be deleted if it is idle within a given time threshold. This feature is related to the containerssh audit. A new daemon process `envd-server-gc` can be introduced. It watches the audit log file directory, then delete the idle pods.

#### Port forwarding

There are several ports will be used in envd:

- sshd server port. envd randomly selects a host port for sshd.
- jupyter notebook port.
- RStudio server port.

The section is to be added. SSH port is forwarded by the ssh proxy server. Other ports may need ssh tunneling.

#### Data and code integration

`runtime.docker.mount` and other docker-specific funcs will be ignored in kubernetes context. `runtime.kubernetes.mount()` should create/use the corresponding PV/PVC in the backend pod.

```python
runtime.kubernetes.mount("")
```

### Implementation details/Notes/Constraints

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
	BuilderAddr   string      `json:"builder_addr,omitempty"`

	Runner     RunnerType `json:"runner,omitempty"`
	RunnerAddr *string    `json:"runner_addr,omitempty"`
}

type BuilderType string

const (
	BuilderTypeDocker     BuilderType = "docker-container"
	BuilderTypeKubernetes BuilderType = "kube-pod"
)

type RunnerType string

const (
	RunnerTypeDocker     RunnerType = "docker"
	RunnerTypeKubernetes RunnerType = "server"
)
```

The default envd context is:

```go
Context{
    Name:          "default",
    Builder:       BuilderTypeDocker,
    BuilderSocket: "unix:///var/run/docker.sock",
    Runner:        RunnerTypeDocker,
    RunnerAddr:    nil,
}
```

Users may use `envd context create` to create new context and use it:

```
$ envd context create --name my-context --builder-type docker-container --runner-type docker --docker-host unix:///var/run/docker.sock
```

##### `envd-server` runner

envd CLI registers the local public key to the envd server, and request the address of the ssh proxy server `containerssh` if the user use the `envd-server` runner.

User registration is not supported in the current design. envd CLI sends a POST request to the envd server.

```
POST /v1/auth
request: {
	"public_key": CONTENT
}
response: {
	"identity_token": "sadasfafd",
	"port": "2222",
	"status_message": "Login successfully"
}
```

The server returns an identity token, to identify the user. The identity token will be persisted in `~/.config/envd/auth`. The port of the containerssh will be returned too. It should be kept in the context.

#### Engine interface

`pkg/envd/engine.Engine` should be generalized to be the interface, docker and envd engine should be the two implementations for the interface.

```go
type Engine interface {
	StartEnvironment()

  ListEnvironments()
  GetEnvironment()
  PauseEnvironment()
  ResumeEnvironment()
  DestroyEnvironment()

  GPUEnabled()
}
```

And, the exising `docker.Client` will only be used to:

- Load the image into the local docker host via `builder`
- Check if the buildkitd container is running in `buildkit.Client`

`builder` and `buildkit.Client` still needs to use `docker.Client` since we keep using docker to build the images.

The pseudocode of the new logic will be like:

```go
// Destroy the environment
func destroy(clicontext *cli.Context) error {
    runtime := engine.New(getCLIFlag("runtime-vendor"))
    runtime.DestroyEnvironment(...)
}

func up(clicontext *cli.Context) error {
    // Build the image first
    builder.Build()
    // Push the image to a registry
    pushImage()
    runtime := engine.New(getCLIFlag("runtime-vendor"))
    runtime.StartEnvironment(...)
}
```

#### `envd create`

`envd create` will be introduced as the new command. It builds the image, and launches the container/pod. If the runner is `envd-server`, the image needs to be pushed first.

Here, the endpoint from `envd-server` needs to be designed carefully.

```
POST /v1/environments
To be determined.
```

#### `envd ssh`

`envd ssh` will be introduced as the new command. It creates the ssh client config, then connect to the ssh proxy server `containerssh`. The username should be the contains the identitytoken and project name.

The `envd-server` should:

- Get the identity token and the project name, then find the correct pod to attach.
- Update `containerssh` app config with the `/v1/config` endpoint.

#### `envd-sshd`

The `envd-sshd` is a tiny sshd implementation which is embeded into the image, and will be used in the resulting environment. envd CLI generates the entrypoint for envrironment image. The command for `envd-sshd` will be like:

```
envd-sshd --port <port> --key <public key path>
```

The information needed by the sshd is passed to the binary via CLI arguments. It does not work with the Kubernetes design. The binary `envd-sshd` should read the environment variable `ENVD_SSH_KEY_PATH` to load the keys, instead of the CLI argument. The backend pods will be created with the same public key mounted in runtime. The public key is generated by the envd API server.

The envd API server checks if the public key and signature provided by the user is valid. If both are valid, the server will use the generated private key and public key to connect to the backend pod.

### Test Plan

[kind](https://github.com/kubernetes-sigs/kind)-based integration test should be added. 
