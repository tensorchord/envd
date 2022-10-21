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
ENVD_VERSION="${ENVD_VERSION:-$GIT_TAG_VERSION}"
DOCKER_HUB_ORG="${DOCKER_HUB_ORG:-tensorchord}"
TAG_SUFFIX="${TAG_SUFFIX:-}"

cd ${ROOT_DIR}

docker buildx build \
    --build-arg ENVD_VERSION=${ENVD_VERSION} \
    -t ${DOCKER_HUB_ORG}/envd:${ENVD_VERSION}-rootless \
    --pull --push --platform linux/x86_64,linux/arm64 \
    -f envd-daemonless.Dockerfile ../../

cd - > /dev/null
