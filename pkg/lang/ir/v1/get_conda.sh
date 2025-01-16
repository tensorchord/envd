set -euo pipefail && \
UNAME_M="$(uname -m)" && \
if [ "${UNAME_M}" = "x86_64" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-x86_64.sh"; \
	SHA256SUM="807774bae6cd87132094458217ebf713df436f64779faf9bb4c3d4b6615c1e3a"; \
elif [ "${UNAME_M}" = "s390x" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-s390x.sh"; \
	SHA256SUM="bb499b18dbcbb2d89b22f91fe26fe661f5ed1f1944fdc743560d69cd52a2468f"; \
elif [ "${UNAME_M}" = "aarch64" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-aarch64.sh"; \
	SHA256SUM="a8846ade7a5ddd9b6a6546590054d70d1c2cbe4fbe8c79fb70227e8fd93ef9f8"; \
fi && \
wget "${MINICONDA_URL}" -O /tmp/miniconda.sh && \
echo "${SHA256SUM}  /tmp/miniconda.sh" > /tmp/shasum && \
if [ "${CONDA_VERSION}" != "latest" ]; then sha256sum -c -s /tmp/shasum; fi
