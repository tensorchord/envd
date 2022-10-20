#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR=`dirname $0`

GIT_TAG_VERSION=$(git describe --tags --abbrev=0 | sed -r 's/[v]+//g') # remove v from version
ENVD_VERSION="${ENVD_VERSION:-$GIT_TAG_VERSION}"
DOCKER_HUB_ORG="${DOCKER_HUB_ORG:-tensorchord}"
TAG_SUFFIX="${TAG_SUFFIX:-}"

cd ${ROOT_DIR}

docker buildx build \
    --build-arg ENVD_VERSION=${ENVD_VERSION} \
    -t ${DOCKER_HUB_ORG}/envd:${ENVD_VERSION} \
    --pull --push --platform linux/x86_64,linux/arm64 \
    -f envd-daemonless.Dockerfile ../../

cd - > /dev/null
