set -euo pipefail && \
UNAME_M="$(uname -m)" && \
if [ "${UNAME_M}" = "x86_64" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-x86_64.sh"; \
	SHA256SUM="238abad23f8d4d8ba89dd05df0b0079e278909a36e06955f12bbef4aa94e6131"; \
elif [ "${UNAME_M}" = "aarch64" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-aarch64.sh"; \
	SHA256SUM="4e0723b9d76aa491cf22511dac36f4fdec373e41d2a243ff875e19b8df39bf94"; \
fi && \
wget "${MINICONDA_URL}" -O /tmp/miniconda.sh && \
echo "${SHA256SUM}  /tmp/miniconda.sh" > /tmp/shasum && \
if [ "${CONDA_VERSION}" != "latest" ]; then sha256sum -c -s /tmp/shasum; fi
