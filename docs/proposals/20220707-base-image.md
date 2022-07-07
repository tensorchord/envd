# Base image customization

Authors:
- [Ce Gao](https://github.com/gaocegege)

## Summary

The base image used in envd is hard coded in [pkg/lang/ir package](https://github.com/tensorchord/envd/tree/main/pkg/lang/ir). They are built from [base-images](https://github.com/tensorchord/envd/tree/main/base-images).

This proposal is to support customization of the base image.

## Motivation

AI/ML infra teams need to change base images especially when they have their internal storage backends (so that they need the SDK to access the storage layer).

Besides this, it also requires customization if envd integrates with [Flyte](https://flyte.org/) (See [feat(integration): Integrate with Flyte #541](https://github.com/tensorchord/envd/issues/541))

## Glossary of Terms

N/A

## Goals

- Customize base images during build process
- Support amd64/arm64
- Support [multi-target build](https://github.com/tensorchord/envd/issues/403)

## Non-Goals

- Support non-Ubuntu base images
- Support [feat(lang): Extend envd build process to help system administrators customize #557](https://github.com/tensorchord/envd/issues/557)

## Proposal

### User Story

#### Reproduce papers 

As a researcher, I need to reproduce the papers. These papers may provide the docker images. Thus I do not want to reproduce it in envd.

#### Enterprise users

Enterprise users may want to use envd in their own infra. Thus they need to customize the base image.

#### Integrate envd with other systems

Some other projects like [Flyte](https://flyte.org/) may use envd to build the docker images. Thus they need to customize the base image.

#### Build images for large scale training or deployment

Besides development, envd may be used in large scale training or deployment. Thus it needs to customize the base image, so that the image size can be reduced. For example, there is no need to install jupyter notebook for serving purpose.

### Language

```python
def build():
    base(image="docker.io/tensorchord/custom-image")
```

New argument `image` is introduced in `base` to support customization. There are some requirements for the image.

- Should be Ubuntu-based image
- Should have `pip` installed

Besides this, we need to introduce a new func `config.ssh`. If users specify the base image, and there is no func calls like `config.jupyter` or `config.ssh`, we consider it as an image not for development. Thus there are some different situations:

#### Custom image for serving or other purposes

```python
def build():
    base(image="docker.io/tensorchord/custom-image")
    # We consider it as an image not for development.
    install.python_packages(packages=["tensorflow"])
```

#### Custom image for development

```python
def build():
    base(image="docker.io/tensorchord/custom-image")
    # We consider it as an image for development. 
    # Thus we should install sshd and jupyter. for it.
    config.jupyter()
    config.ssh()
```

#### Default image for development

```python
def build():
    # We consider it as an image for development if there is
    # no image specified in `base` func. But we do not need 
    # to install sshd and jupyter since the default base image
    # contains them.
    base(os="ubuntu", version="20.04")
```

#### Default image for serving or other purposes

```python
def build():
    # We consider it as an image not for development if there is
    # no image specified in `base` func.
    base(os="ubuntu", version="20.04")
    config.ssh(mode="disable")
    install.python_packages(packages=["tensorflow"])
```

## Development Process

Already supported by envd:

- Default image for development

### Stage 1

Support:

- Custom image for serving or other purposes

### Stage 2

Support:

- Default image for serving or other purposes
- Custom image for development (Huge engineering work)
