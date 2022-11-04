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

GIT_TAG_VERSION=$(git describe --tags --abbrev=0 | sed -r 's/[v]+//g') # remove v from version
ENVD_VERSION="${ENVD_VERSION:-$GIT_TAG_VERSION}"
DOCKER_HUB_ORG="${DOCKER_HUB_ORG:-tensorchord}"
BUILD_FUNC="${BUILD_FUNC:-build}"
TAG_SUFFIX="${TAG_SUFFIX:-}"

cd ${ROOT_DIR}

envd --debug build -f build.envd:${BUILD_FUNC} --export-cache type=registry,ref=docker.io/${DOCKER_HUB_ORG}/python-cache:envd-v${ENVD_VERSION}${TAG_SUFFIX} --force

cd - > /dev/null
