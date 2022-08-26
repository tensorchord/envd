# Add Runtime Redistributed Metadata To Image Metadata
Authors:
- [nullday](https://github.com/aseaday)

## Summary

The `envd` would add (runing graph metadata)[https://github.com/tensorchord/envd/blob/630ada172bdf876c3b749329fdbe284c108051f2/pkg/lang/ir/types.go#L70] would be encoded into a ASCII string and added to be image(OCI Spec) as a label. 

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

