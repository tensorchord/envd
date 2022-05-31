ARG ENVD_VERSION
ARG ENVD_SSH_IMAGE
FROM ubuntu:20.04 as base

FROM base as base-amd64

FROM base as base-arm64

FROM ${ENVD_SSH_IMAGE}:${ENVD_VERSION} AS envd

FROM base-${TARGETARCH}

ARG TARGETARCH

LABEL maintainer "envd-maintainers <envd-maintainers@tensorchord.ai>"

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends --no-install-suggests --fix-missing bash-static \
    apt-utils libtinfo5 libncursesw5 \
    # envd dependencies
    python3 curl openssh-client git tini sudo python3-pip zsh vim \
    && rm -rf /var/lib/apt/lists/*

COPY --from=envd /usr/bin/envd-ssh /var/envd/bin/envd-ssh
