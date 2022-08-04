#!/usr/bin/env bash

ROOT_DIR=`dirname $0`

ENVD_VERSION="${ENVD_VERSION:-$GIT_TAG_VERSION}"

cd ${ROOT_DIR}

envd build --export-cache type=registry,ref=docker.io/gaocegege/python-cache:3.9-envd-v${ENVD_VERSION} --force

cd - > /dev/null
