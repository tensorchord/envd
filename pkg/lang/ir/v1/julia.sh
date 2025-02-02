set -o pipefail && \
UNAME_M="$(uname -m)" && \
if [ "${UNAME_M}" = "x86_64" ]; then \
    JULIA_URL="https://julialang-s3.julialang.org/bin/linux/x64/1.10/julia-1.10.8-linux-x86_64.tar.gz"; \
    SHA256SUM="0410175aeec3df63173c15187f2083f179d40596d36fd3a57819cc5f522ae735"; \
elif [ "{UNAME_M}" = "aarch64" ]; then \
    JULIA_URL="https://julialang-s3.julialang.org/bin/linux/aarch64/1.10/julia-1.10.8-linux-aarch64.tar.gz" \
    SHA256SUM="8d63dd12595a08edc736be8d6c4fea1840f137b81c62079d970dbd1be448b8cd"; \
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
