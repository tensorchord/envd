#!/usr/bin/env bash

set -e

ROOT_DIR=`dirname $0`

GIT_TAG_VERSION=$(git describe --tags --abbrev=0 | sed -r 's/[v]+//g') # remove v from version
ENVD_VERSION="${ENVD_VERSION:-$GIT_TAG_VERSION}"

cd ${ROOT_DIR}

envd build --export-cache type=registry,ref=docker.io/tensorchord/python-cache:3.9-envd-v${ENVD_VERSION} --force

cd - > /dev/null
