# Add Runtime Configuration to Image for Envd Image Redistribution
Authors:
- [nullday](https://github.com/aseaday)

## Summary

The `envd` would add (runing graph metadata)[https://github.com/tensorchord/envd/blob/630ada172bdf876c3b749329fdbe284c108051f2/pkg/lang/ir/types.go#L70] would be encoded into a ASCII string and added to be image(OCI Spec) as a config. 

we named the above mentioned encoded config `Envd Runtime Graph Label`

## Motivation

This proposal is part of effort to define what the artifact `envd` delivery and decoupling the phases of build and runing. It will be friendly for running a envd environment even at the absence of `build.env`.

the concept of `running context` is a also needed as addition of `build contxt` for example:

- An engineer build a easy-to-use env for his/her interns for the quick use of company tools.
- Kubernetes remote runtime support in the future.

## Goals
- Provide internal API:
    - read runtime metadata
    - write runtime metadata
- `up` and `run` command doesn't need to interpret the `build.env`. `run` command can run a built env image with attachment of current working directory as `running context`.

## Implementations

There are two parts of runtime configuration that can be used in envd.
- OCI sepecifications specific:
    - ExposedPort
    - Entrypoint
    - Env
    - Cmd
- Custome Labels

For some parts of runtime configuration, we could use the OCI part such as environment variables. We still need to deal with extra parts such as port bindings which not covered by the OCI spec.

```golang
type RuntimeGraph struct {
	RuntimeCommands map[string]string
	RuntimeDaemon   [][]string
	RuntimeEnviron  map[string]string
	RuntimeExpose   []ExposeItem
}
```

- *RuntimeCommands* will not be added to the runtime label because it is one-time only
- *RuntimeEnviron* is completely in accordance with the OCI defined configuration `Env`. So we don't need to add it to the runtime label.
- *RuntimeDaemon* and *RuntimeExpose* will be added to the runtime label and we also should update the config the OCI configuration's Entrypoint and ExposedPort accordingly.

we use the following labels:

- ai.tensorchord.envd.runtimeGraph.version
- ai.tensorchord.envd.runtimeGraph.Daemon
- ai.tensorchord.envd.runtimeGraph.Expose


#### Why msgpack

MessagePack is supported by over 50 programming languages and environments. And it will be shorter than JSON.