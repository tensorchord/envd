#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR=`dirname $0`

GIT_TAG_VERSION=$(git describe --tags --abbrev=0 | sed -r 's/[v]+//g') # remove v from version
ENVD_VERSION="${ENVD_VERSION:-$GIT_TAG_VERSION}"
DOCKER_HUB_ORG="${DOCKER_HUB_ORG:-tensorchord}"
BUILD_FUNC="${BUILD_FUNC:-build}"
TAG_SUFFIX="${TAG_SUFFIX:-}"

cd ${ROOT_DIR}

envd --debug build -f build.envd:${BUILD_FUNC} --export-cache type=registry,ref=docker.io/${DOCKER_HUB_ORG}/python-cache:envd-v${ENVD_VERSION}${TAG_SUFFIX} --force
# envd build -f build.envd:build_gpu_11_2 --export-cache type=registry,ref=docker.io/${DOCKER_HUB_ORG}/python-cache:envd-v${ENVD_VERSION}-cuda-11.2.0-cudnn-8 --force
# envd build -f build.envd:build_gpu_11_3 --export-cache type=registry,ref=docker.io/${DOCKER_HUB_ORG}/python-cache:envd-v${ENVD_VERSION}-cuda-11.3.0-cudnn-8 --force
# envd build -f build.envd:build_gpu_11_6 --export-cache type=registry,ref=docker.io/${DOCKER_HUB_ORG}/python-cache:envd-v${ENVD_VERSION}-cuda-11.6.0-cudnn-8 --force

cd - > /dev/null
