#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR=`dirname $0`

GIT_TAG_VERSION=$(git describe --tags --abbrev=0)
DOCKER_HUB_ORG="${DOCKER_HUB_ORG:-tensorchord}"
ENVD_OS="${ENVD_OS:-ubuntu20.04}"
JULIA_VERSION="${JULIA_VERSION:-1.8rc1}"
RLANG_VERSION="${RLANG_VERSION:-4.2}"

cd ${ROOT_DIR}
# ubuntu 22.04 build require moby/buildkit version greater than 0.8.1
if ! docker buildx inspect cuda; then
    docker buildx create --use --platform linux/x86_64,linux/arm64,linux/ppc64le --driver-opt image=moby/buildkit:v0.10.3
fi

# https://github.com/docker/buildx/issues/495#issuecomment-754688157
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes

# TODO(gaocegege): Support linux/arm64
docker buildx build \
    --build-arg ENVD_VERSION=${GIT_TAG_VERSION} \
    --build-arg ENVD_SSHD_IMAGE=tensorchord/envd-sshd-from-scratch \
    -t ${DOCKER_HUB_ORG}/r-base:${RLANG_VERSION}-envd-${GIT_TAG_VERSION} \
    --pull --push --platform linux/x86_64 \
    -f r${RLANG_VERSION}.Dockerfile .
docker buildx build \
    --build-arg ENVD_VERSION=${GIT_TAG_VERSION} \
    --build-arg ENVD_SSHD_IMAGE=tensorchord/envd-sshd-from-scratch \
    -t ${DOCKER_HUB_ORG}/julia:${JULIA_VERSION}-${ENVD_OS}-envd-${GIT_TAG_VERSION} \
    --pull --push --platform linux/x86_64,linux/arm64 \
    -f julia${JULIA_VERSION}-${ENVD_OS}.Dockerfile .
cd - > /dev/null
