set -o pipefail && \
UNAME_M="$(uname -m)" && \
if [ "${UNAME_M}" = "x86_64" ]; then \
    JULIA_URL="https://julialang-s3.julialang.org/bin/linux/x64/1.10/julia-1.10.10-linux-x86_64.tar.gz"; \
    SHA256SUM="6a78a03a71c7ab792e8673dc5cedb918e037f081ceb58b50971dfb7c64c5bf81"; \
elif [ "{UNAME_M}" = "aarch64" ]; then \
    JULIA_URL="https://julialang-s3.julialang.org/bin/linux/aarch64/1.10/julia-1.10.10-linux-aarch64.tar.gz" \
    SHA256SUM="a4b157ed68da10471ea86acc05a0ab61c1a6931ee592a9b236be227d72da50ff"; \
fi && \

wget "${JULIA_URL}" -O /tmp/julia.tar.gz && \
echo "${SHA256SUM}  /tmp/julia.tar.gz" > /tmp/sha256sum && \
sha256sum -c -s /tmp/sha256sum
EXIT_CODE=$?
if [ $EXIT_CODE -ne 0 ]; then
    echo "CHECKSUM FAILED" && \
    rm /tmp/julia.tar.gz && \
    wget "${JULIA_URL}" -O /tmp/julia.tar.gz && \
    sha256sum -c -s /tmp/sha256sum
else
    echo "CHECKSUM PASSED"
fi
