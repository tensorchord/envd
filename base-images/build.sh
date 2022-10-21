#!/usr/bin/env bash
# Copyright 2022 The envd Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -euo pipefail

ROOT_DIR=`dirname $0`

GIT_TAG_VERSION=$(git describe --tags --abbrev=0)
DOCKER_HUB_ORG="${DOCKER_HUB_ORG:-tensorchord}"
ENVD_OS="${ENVD_OS:-ubuntu20.04}"
JULIA_VERSION="${JULIA_VERSION:-1.8rc1}"
RLANG_VERSION="${RLANG_VERSION:-4.2}"

cd ${ROOT_DIR}
# ubuntu 22.04 build require moby/buildkit version greater than 0.8.1
docker buildx create --use --platform linux/x86_64,linux/arm64,linux/ppc64le --driver-opt image=moby/buildkit:v0.10.3

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
