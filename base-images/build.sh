#!/usr/bin/env bash

ROOT_DIR=`dirname $0`

cd ${ROOT_DIR}
# ubuntu 22.04 build require moby/buildkit version greater than 0.8.1
if ! docker buildx inspect cuda; then
    docker buildx create --use --platform linux/x86_64,linux/arm64,linux/ppc64le --driver-opt image=moby/buildkit:v0.10.3 --name cuda --node cuda
fi
docker buildx build --build-arg IMAGE_NAME=docker.io/nvidia/cuda \
    --build-arg ENVD_VERSION=0.0.1-alpha.5 \
    --build-arg ENVD_SSH_IMAGE=ghcr.io/tensorchord/envd-ssh-from-scratch \
    --build-arg HTTP_PROXY=${HTTP_PROXY} \
     --build-arg HTTPS_PROXY=${HTTPS_PROXY} \
    --pull --push --platform linux/x86_64,linux/arm64 \
    -t gaocegege/python:3.8-ubuntu2004-cuda11.6-cudnn8 \
    -f python3.8-ubuntu2004-cuda11.6.Dockerfile .
cd - > /dev/null
