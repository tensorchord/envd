set -euo pipefail && \
ARCH="$(uname -m)" && \
if [[ "${ARCH}" == "aarch64" ]]; then \
    ARCH="aarch64"; \
elif [[ "${ARCH}" == "ppc64le" ]]; then \
    ARCH="ppc64le"; \
else \
    ARCH="64"; \
fi && \
mkdir -p ${MAMBA_BIN_DIR} && \
curl -Ls https://micro.mamba.pm/api/micromamba/linux-${ARCH}/${MAMBA_VERSION} | tar -xvj -C ${MAMBA_BIN_DIR} --strip-components=1 bin/micromamba && \
chown $(id -u):$(id -g) ${MAMBA_BIN_DIR}/micromamba
ln -s ${MAMBA_BIN_DIR}/micromamba ${MAMBA_BIN_DIR}/conda && \
echo -e "channels:\n  - defaults" > ${MAMBA_ROOT_PREFIX}/.mambarc
echo -e "#!/bin/sh\n\. ${MAMBA_ROOT_PREFIX}/etc/profile.d/micromamba.sh || return \$?\nmicromamba activate \"\$@\"" > ${MAMBA_BIN_DIR}/activate
