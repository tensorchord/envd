#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR=`dirname $0`

GIT_TAG_VERSION=$(git describe --tags --abbrev=0 | sed -r 's/[v]+//g') # remove v from version
ENVD_VERSION="${ENVD_VERSION:-$GIT_TAG_VERSION}"
DOCKER_IMAGE_TAG="${DOCKER_IMAGE_TAG:-v$ENVD_VERSION}"
DOCKER_HUB_ORG="${DOCKER_HUB_ORG:-tensorchord}"
PYTHON_VERSION="${PYTHON_VERSION:-3.9}"
ENVD_OS="${ENVD_OS:-ubuntu20.04}"
JULIA_VERSION="${JULIA_VERSION:-1.8rc1}"
RLANG_VERSION="${RLANG_VERSION:-4.2}"

cd ${ROOT_DIR}
# ubuntu 22.04 build require moby/buildkit version greater than 0.8.1
if ! docker buildx inspect cuda; then
    docker buildx create --use --platform linux/x86_64,linux/arm64,linux/ppc64le --driver-opt image=moby/buildkit:v0.10.3 --name cuda --node cuda
fi

# https://github.com/docker/buildx/issues/495#issuecomment-754688157
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes

docker buildx build \
    --build-arg ENVD_VERSION=${ENVD_VERSION} \
    --build-arg ENVD_SSH_IMAGE=ghcr.io/tensorchord/envd-sshd-from-scratch \
    --pull --push --platform linux/x86_64,linux/arm64 \
    -t ${DOCKER_HUB_ORG}/python:${PYTHON_VERSION}-${ENVD_OS}-envd-${DOCKER_IMAGE_TAG} \
    -f python${PYTHON_VERSION}-${ENVD_OS}.Dockerfile .
docker buildx build --build-arg IMAGE_NAME=docker.io/nvidia/cuda \
    --build-arg ENVD_VERSION=${ENVD_VERSION} \
    --build-arg ENVD_SSH_IMAGE=ghcr.io/tensorchord/envd-sshd-from-scratch \
    --pull --push --platform linux/x86_64,linux/arm64 \
    -t ${DOCKER_HUB_ORG}/python:${PYTHON_VERSION}-${ENVD_OS}-cuda11.2-cudnn8-envd-${DOCKER_IMAGE_TAG} \
    -f python${PYTHON_VERSION}-${ENVD_OS}-cuda11.2.Dockerfile .

# TODO(gaocegege): Support linux/arm64
docker buildx build \
    --build-arg ENVD_VERSION=${ENVD_VERSION} \
    --build-arg ENVD_SSH_IMAGE=ghcr.io/tensorchord/envd-sshd-from-scratch \
    -t ${DOCKER_HUB_ORG}/r-base:${RLANG_VERSION}-envd-${DOCKER_IMAGE_TAG} \
    --pull --push --platform linux/x86_64 \
    -f r${RLANG_VERSION}.Dockerfile .
docker buildx build \
    --build-arg ENVD_VERSION=${ENVD_VERSION} \
    --build-arg ENVD_SSH_IMAGE=ghcr.io/tensorchord/envd-sshd-from-scratch \
    -t ${DOCKER_HUB_ORG}/julia:${JULIA_VERSION}-${ENVD_OS}-envd-${DOCKER_IMAGE_TAG} \
    --pull --push --platform linux/x86_64,linux/arm64 \
    -f julia${JULIA_VERSION}-${ENVD_OS}.Dockerfile .
cd - > /dev/null
