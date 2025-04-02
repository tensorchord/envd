set -euo pipefail && \
UNAME_M="$(uname -m)" && \
if [ "${UNAME_M}" = "x86_64" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-x86_64.sh"; \
	SHA256SUM="d8c1645776c0758214e4191c605abe5878002051316bd423f2b14b22d6cb4251"; \
elif [ "${UNAME_M}" = "s390x" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-s390x.sh"; \
	SHA256SUM="0b4d5a3f16dcb2d230ba5dfdfdb848c854006aab6dd1bd3dbf29fcddf04b07a4"; \
elif [ "${UNAME_M}" = "aarch64" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-aarch64.sh"; \
	SHA256SUM="8a1d4407fce7ec552ac6ed655ce93d83549e02b819cacefbb7f640f9051e638b"; \
fi && \
wget "${MINICONDA_URL}" -O /tmp/miniconda.sh && \
echo "${SHA256SUM}  /tmp/miniconda.sh" > /tmp/shasum && \
if [ "${CONDA_VERSION}" != "latest" ]; then sha256sum -c -s /tmp/shasum; fi
