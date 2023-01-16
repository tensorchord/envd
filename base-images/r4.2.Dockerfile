ARG ENVD_VERSION
ARG ENVD_SSHD_IMAGE
FROM rocker/r-ver:4.2 as base

FROM base as base-amd64

FROM base as base-arm64

FROM ${ENVD_SSHD_IMAGE}:${ENVD_VERSION} AS envd

FROM base-${TARGETARCH}

ARG TARGETARCH

LABEL maintainer "envd-maintainers <envd-maintainers@tensorchord.ai>"

ENV DEBIAN_FRONTEND noninteractive
ENV PATH="/usr/bin:${PATH}"
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

RUN apt-get update && \
    apt-get install -y --no-install-recommends --no-install-suggests --fix-missing \
    apt-utils bash-static libtinfo5 libncursesw5 \
    # rstudio dependencies
    file libapparmor1 libclang-dev libcurl4-openssl-dev libedit2 libobjc4 wget libssl-dev \
    libpq5 psmisc procps python3-setuptools pwgen lsb-release \
    # envd dependencies
    python3 curl openssh-client git tini sudo zsh vim \
    && rm -rf /var/lib/apt/lists/* \
    # prompt
    && curl --proto '=https' --tlsv1.2 -sSf https://starship.rs/install.sh | sh -s -- -y

RUN set -x && \
    UNAME_M="$(uname -m)" && \
    if [ "${UNAME_M}" = "x86_64" ]; then \
      RSTUDIO_URL="https://download2.rstudio.org/server/jammy/amd64/rstudio-server-2022.12.0-353-amd64.deb"; \
    elif [ "${UNAME_M}" = "aarch64" ]; then \
      RSTUDIO_URL="https://rstudio.org/download/latest/latest/server/focal/rstudio-server-latest-arm64.deb"; \
    fi && \
    DOWNLOAD_FILE=rstudio-server.deb && \
    wget "${RSTUDIO_URL}" -O ${DOWNLOAD_FILE} && \
    dpkg -i "$DOWNLOAD_FILE" && \
    rm ${DOWNLOAD_FILE} && rm -f /var/lib/rstudio-server/secure-cookie-key

COPY --from=envd /usr/bin/envd-sshd /var/envd/bin/envd-sshd
