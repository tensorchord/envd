set -euo pipefail && \
UNAME_M="$(uname -m)" && \
if [ "${UNAME_M}" = "x86_64" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-x86_64.sh"; \
	SHA256SUM="4ee9c3aa53329cd7a63b49877c0babb49b19b7e5af29807b793a76bdb1d362b4"; \
elif [ "${UNAME_M}" = "s390x" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-s390x.sh"; \
	SHA256SUM="e5e5e89cdcef9332fe632cd25d318cf71f681eef029a24495c713b18e66a8018"; \
elif [ "${UNAME_M}" = "aarch64" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-aarch64.sh"; \
	SHA256SUM="00c7127a8a8d3f4b9c2ab3391c661239d5b9a88eafe895fd0f3f2a8d9c0f4556"; \
elif [ "${UNAME_M}" = "ppc64le" ]; then \
	MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-${CONDA_VERSION}-Linux-ppc64le.sh"; \
	SHA256SUM="8ee1f8d17ef7c8cb08a85f7d858b1cb55866c06fcf7545b98c3b82e4d0277e66"; \
fi && \
wget "${MINICONDA_URL}" -O /tmp/miniconda.sh && \
echo "${SHA256SUM} /tmp/miniconda.sh" > /tmp/shasum && \
if [ "${CONDA_VERSION}" != "latest" ]; then sha256sum --check --status /tmp/shasum; fi && \
mkdir -p /opt && \
sh /tmp/miniconda.sh -b -u -p /opt/conda && \
rm /tmp/miniconda.sh /tmp/shasum && \
echo ". /opt/conda/etc/profile.d/conda.sh" >> ~/.bashrc && \
echo "conda activate base" >> ~/.bashrc && \
echo -e "channels:\n  - defaults" > /opt/conda/.condarc && \
find /opt/conda/ -follow -type f -name '*.a' -delete && \
find /opt/conda/ -follow -type f -name '*.js.map' -delete && \
/opt/conda/bin/conda clean -afy
